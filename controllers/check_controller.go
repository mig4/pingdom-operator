/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/go-logr/logr"
	"github.com/russellcardullo/go-pingdom/pingdom"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	observabilityv1alpha1 "gitlab.com/mig4/pingdom-operator/api/v1alpha1"
	"gitlab.com/mig4/pingdom-operator/controllers/finalizer"
	checkreconciler "gitlab.com/mig4/pingdom-operator/controllers/resources/check"
)

// CheckReconciler reconciles a Check object
type CheckReconciler struct {
	client.Client
	Log      logr.Logger
	PdApiKey string
}

// +kubebuilder:rbac:groups=observability.pingdom.mig4.gitlab.io,resources=checks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=observability.pingdom.mig4.gitlab.io,resources=checks/status,verbs=get;update;patch

func (r *CheckReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("resource", "check", "namespacedName", req.NamespacedName)

	log.V(1).Info("entered Reconcile, retrieving Check")
	var check observabilityv1alpha1.Check
	if err := r.Get(ctx, req.NamespacedName, &check); err != nil {
		log.Error(err, "Check not found")
		// Ignore NotFound errors, we'll get a new notification when the
		// object exists.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// TODO: move to a defaulting webhook
	r.initCheckDefaults(log, req, &check)

	// Initialise Pingdom client
	pdClient, err := r.initPingdomClient(
		ctx, check.GetNamespace(), check.Spec.CredentialsSecret.Name,
	)
	if err != nil {
		log.Error(err, "Unable to initialise Pingdom client")
		return ctrl.Result{}, err
	}

	// Initialise Pingdom resource reconciler and a finalizer manager for it
	reconciler := checkreconciler.New(&checkreconciler.Config{
		Logger:   log,
		PdClient: pdClient,
		Check:    &check,
	})
	finalizerMgr := finalizer.New(log, r.Client, reconciler)

	// Ensure finalizer is registered
	if err := finalizerMgr.EnsureAttached(ctx, &check); err != nil {
		return ctrl.Result{}, microerror.Maskf(err, "failure handling finalizer")
	}

	// Refresh internal representation of state of the external resource
	if err := reconciler.RefreshState(ctx); err != nil {
		return ctrl.Result{}, microerror.Maskf(
			err, "failure refreshing state of the check",
		)
	}

	// Ensure to update status subresource of the CRD in Kube when we're done
	// so it's visible to users; do it later so it picks up any changes the
	// reconciler makes to the object (like setting the ID initially)
	defer r.updateStatus(log, ctx, &check)

	// Ensure external resource (Pingdom check) matches desired spec
	if err := reconciler.EnsureState(ctx); err != nil {
		return ctrl.Result{}, microerror.Maskf(
			err, "failure reconciling external resource",
		)
	}

	// If object is being deleted and we got here without errors means external
	// resource is already gone and we can remove the finalizer.
	if !check.GetDeletionTimestamp().IsZero() {
		if err := finalizerMgr.EnsureDetached(ctx, &check); err != nil {
			return ctrl.Result{}, microerror.Maskf(err, "failure handling finalizer")
		}
	}

	// Schedule next run after some time.
	nextIn := "1m" // normally "5m" == default Pingdom resolution
	if reconciler.DidWork() {
		log.V(1).Info("reconciler.DidWork == true")
		nextIn = "10s"
	}
	log.Info("exiting Reconcile, scheduling next run", "nextIn", nextIn)
	delay, _ := time.ParseDuration(nextIn)
	return ctrl.Result{RequeueAfter: delay}, nil
}

/*
initPingdomClient reads the secrets and initialises a Pingdom API client.
*/
func (r *CheckReconciler) initPingdomClient(
	ctx context.Context,
	namespace, secretName string,
) (*pingdom.Client, error) {
	secretNsName := types.NamespacedName{
		Namespace: namespace,
		Name:      secretName,
	}
	secret := &corev1.Secret{}
	if err := r.Get(ctx, secretNsName, secret); err != nil {
		return nil, err
	}

	if secret.Data["user"] == nil {
		return nil, fmt.Errorf(
			"Pingdom API username not found in secret %v", secretNsName,
		)
	}
	if secret.Data["password"] == nil {
		return nil, fmt.Errorf(
			"Pingdom API password not found in secret %v", secretNsName,
		)
	}

	pdClient, err := pingdom.NewClientWithConfig(pingdom.ClientConfig{
		User:     string(secret.Data["user"]),
		Password: string(secret.Data["password"]),
		APIKey:   r.PdApiKey,
	})
	if err != nil {
		return nil, err
	}
	return pdClient, nil
}

/*
initCheckDefaults initializes defaults on the Check object

TODO: this should be done if a defaulting webhook but in the meantime we need
  it to be able to pass validation before the resource is created on Pingdom
  and its status reflected in Check's Status object
*/
func (r *CheckReconciler) initCheckDefaults(
	parentLog logr.Logger,
	req ctrl.Request,
	check *observabilityv1alpha1.Check,
) {
	log := parentLog.WithValues("stage", "initDefaults")
	if check.Spec.Name == nil {
		check.Spec.Name = &req.Name
		log.V(1).Info("using default Name", "name", req.Name)
	}
	if check.Status.CreatedTime.IsZero() {
		check.Status.CreatedTime = check.CreationTimestamp
		log.V(1).Info("using default CreatedTime", "created", check.Status.CreatedTime)
	}
	if check.Status.Name == nil {
		check.Status.Name = check.Spec.Name
	}
	if string(check.Status.Status) == "" {
		check.Status.Status = observabilityv1alpha1.Unknown
		log.V(1).Info("using default Status", "status", check.Status.Status)
	}
	if string(check.Status.Type) == "" {
		// TODO: technically there should be no default type; this only needs
		//   to be set in order to avoid failing validation before the resource
		//   is created and type is set correctly; this happens when setting up
		//   finalizers initially as we get an empty object and after adding
		//   finalizers and giving that object to Update() it blows up because
		//   the Status nested object is empty and so fails validation
		check.Status.Type = observabilityv1alpha1.CheckType("ping")
		log.V(1).Info("using default Type", "type", check.Status.Type)
	}
}

/*
updateStatus updates the status subresource on the API Check object.

It will skip the update if the object is being deleted as in this case the
external resource may not exist anymore so status would be out-of-date anyway.
Errors are only logged, not returned as a call to updateStatus should be
deferred so there would be no way of handling the error.
*/
func (r *CheckReconciler) updateStatus(
	parentLog logr.Logger,
	ctx context.Context,
	check *observabilityv1alpha1.Check,
) {
	log := parentLog.WithValues("action", "updateStatus")
	if !check.GetDeletionTimestamp().IsZero() {
		log.V(1).Info("skip object status update as object is deleted")
		return
	}
	if err := r.Status().Update(ctx, check); err != nil {
		log.Error(err, "unable to update object status")
		return
	}
	log.V(1).Info("updated object status")
}

func (r *CheckReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&observabilityv1alpha1.Check{}).
		Complete(r)
}
