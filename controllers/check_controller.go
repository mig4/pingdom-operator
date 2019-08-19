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
	"errors"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/russellcardullo/go-pingdom/pingdom"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	observabilityv1alpha1 "gitlab.com/mig4/pingdom-operator/api/v1alpha1"
)

// CheckReconciler reconciles a Check object
type CheckReconciler struct {
	client.Client
	Log      logr.Logger
	PdApiKey string
	pdClient *pingdom.Client
	justChk  time.Time
}

// +kubebuilder:rbac:groups=observability.pingdom.mig4.gitlab.io,resources=checks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=observability.pingdom.mig4.gitlab.io,resources=checks/status,verbs=get;update;patch

func (r *CheckReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("check", req.NamespacedName)

	// TODO: CheckReconciler is globally scoped which means we need to be
	//   careful about any state we store on it; it could be useful to have
	//   a per-request scoped handler we delegate to so that we can store
	//   things like logger with values, initialised Pingdom API client, etc.

	log.Info("entered Reconcile", "justChk", r.justChk)

	// Retrieve our Check object from Kube
	var check observabilityv1alpha1.Check
	if err := r.Get(ctx, req.NamespacedName, &check); err != nil {
		log.Error(err, "Check not found")
		// Ignore NotFound errors, we'll get a new notification when the
		// object exists.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Initialise Pingdom client
	if err := r.InitClient(
		ctx, check.GetNamespace(), check.Spec.CredentialsSecret.Name,
	); err != nil {
		log.Error(err, "Unable to initialise Pingdom client")
		return ctrl.Result{}, err
	}

	// TODO: move to a defaulting webhook
	r.initCheckDefaults(log, req, &check)

	// Handle deletion
	if ok := r.reconcileFinalizationDeletion(ctx, log, &check); !ok {
		// TODO: failure to delete or update the object should be indicated
		//   so it's retried
		return ctrl.Result{}, nil
	}

	// Read the check from Pingdom API here and populate `check.Status`
	if err := r.updateStatusObj(log, &check); err != nil {
		log.Error(err, "Updating Check.Status from Pingdom state failed, attempting to continue")
	}

	// Update status of the CRD so it's visible to users
	// TODO: shouldn't we do this later, so it covers newly created checks?
	if err := r.Status().Update(ctx, &check); err != nil {
		log.Error(err, "unable to update Check status")
		return ctrl.Result{}, err
	}

	// Create/update the resource in Pingdom
	if err := r.updatePingdom(log, &check); err != nil {
		log.Error(err, "Unable to update Pingdom")
		return ctrl.Result{}, err
	}
	// Read the check from Pingdom API here and update `check.Status`
	// TODO: do we need to do it twice?
	if err := r.updateStatusObj(log, &check); err != nil {
		log.Error(err, "Updating Check.Status from Pingdom state failed, attempting to continue")
	}

	// Update status of the CRD so it's visible to users
	if err := r.Status().Update(ctx, &check); err != nil {
		log.Error(err, "unable to update Check status")
		return ctrl.Result{}, err
	}

	r.justChk = time.Now()
	log.Info("exiting Reconcile, next run in 5m", "justChk", r.justChk)
	// TODO: implement a webhook receiver and configure a webhook on Pingdom
	//   side to notify of any changes so we don't need to guess when to
	//   reconcile next to update the status again
	// Update after the check executes on Pingdom next (based on default check
	// resolution).
	delay, _ := time.ParseDuration("5m")
	return ctrl.Result{RequeueAfter: delay}, nil
}

/*
InitClient reads the secrets and initialises a Pingdom API client.
*/
func (r *CheckReconciler) InitClient(ctx context.Context, namespace, secretName string) error {
	secretNsName := types.NamespacedName{
		Namespace: namespace,
		Name:      secretName,
	}
	secret := &corev1.Secret{}
	if err := r.Get(ctx, secretNsName, secret); err != nil {
		return err
	}

	if secret.Data["user"] == nil {
		return fmt.Errorf(
			"Pingdom API username not found in secret %v", secretNsName,
		)
	}
	if secret.Data["password"] == nil {
		return fmt.Errorf(
			"Pingdom API password not found in secret %v", secretNsName,
		)
	}

	pdClient, err := pingdom.NewClientWithConfig(pingdom.ClientConfig{
		User:     string(secret.Data["user"]),
		Password: string(secret.Data["password"]),
		APIKey:   r.PdApiKey,
	})
	if err != nil {
		return err
	}
	r.pdClient = pdClient
	return nil
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
		check.Status.Status = observabilityv1alpha1.CheckResult("unknown")
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
reconcileFinalizationDeletion handles logic related to object deletion.

I.e. setting up or removing finalizers, reacting to object deletion, etc.

It returns a bool indicating if the reconciliation should proceed.
*/
func (r *CheckReconciler) reconcileFinalizationDeletion(
	ctx context.Context,
	parentLog logr.Logger,
	check *observabilityv1alpha1.Check,
) (proceed bool) {
	finalizerName := "pingdom-cleanup." + observabilityv1alpha1.GroupVersion.Group
	log := parentLog.WithValues(
		"stage", "finalizationDeletion",
		"finalizer", finalizerName,
	)
	doUpdate := false
	proceed = true
	log.V(1).Info("handling finalizer registration / deletion")

	if !check.GetDeletionTimestamp().IsZero() {
		log.Info("check object is being deleted")

		if sliceContains(check.GetFinalizers(), finalizerName) {
			log.Info("finalizer present, deleting Pingdom resource")
			if err := r.deletePingdomCheck(parentLog, check); err == nil {
				// remove the finalizer to indicate the object can now be deleted too
				log.Info("removing the finalizer")
				check.SetFinalizers(sliceRemove(check.GetFinalizers(), finalizerName))
				doUpdate = true
			}
			proceed = false
		}
	} else if !sliceContains(check.GetFinalizers(), finalizerName) {
		log.Info("registering the finalizer")
		check.SetFinalizers(append(check.GetFinalizers(), finalizerName))
		doUpdate = true
	}

	if doUpdate {
		if err := r.Update(ctx, check); err != nil {
			log.Error(err, "unable to update the check object")
			proceed = false
		} else {
			log.Info("check object updated")
		}
	}

	return
}

/*
updateStatusObj updates the state of the check stored in CheckStatus object on
the given Check, based on state of the resource retrieved from Pingdom.
*/
func (r *CheckReconciler) updateStatusObj(
	parentLog logr.Logger,
	check *observabilityv1alpha1.Check,
) error {
	status := &check.Status
	checkId := status.Id
	log := parentLog.WithValues("stage", "updateStatusObj", "checkId", checkId)
	if checkId == 0 {
		log.Info("Pingdom resource doesn't exist yet, nothing to update")
		return nil
	}
	log.V(1).Info("fetching check resource from Pingdom")
	pdCheck, err := r.pdClient.Checks.Read(int(checkId))
	if err != nil {
		log.Error(err, "unable to fetch check from Pingdom")
	}

	log.V(1).Info("populating Status")
	status.Name = &pdCheck.Name
	status.Type = observabilityv1alpha1.CheckType(pdCheck.Type.Name)
	status.Host = pdCheck.Hostname
	status.Status = observabilityv1alpha1.CheckResult(pdCheck.Status)
	status.Paused = &pdCheck.Paused
	status.LastErrorTime = parsePdTime(pdCheck.LastErrorTime)
	status.LastTestTime = parsePdTime(pdCheck.LastTestTime)
	status.LastResponseTimeMilis = &pdCheck.LastResponseTime
	status.CreatedTime = *parsePdTime(pdCheck.Created)

	if pdCheck.Type.Name == string(observabilityv1alpha1.Http) {
		if pdCheck.Type.HTTP == nil {
			err = errors.New("check type is http but details not available")
			log.Error(err, "Pingdom API didn't return http check details")
			return err
		}
		status.Port = ptrI32(int32(pdCheck.Type.HTTP.Port))
		status.Url = &pdCheck.Type.HTTP.Url
		status.Encryption = &pdCheck.Type.HTTP.Encryption
	} else if pdCheck.Type.Name == string(observabilityv1alpha1.Tcp) {
		if pdCheck.Type.TCP == nil {
			err = errors.New("check type is tcp but details not available")
			log.Error(err, "Pingdom API didn't return tcp check details")
			return err
		}
		status.Port = ptrI32(int32(pdCheck.Type.TCP.Port))
	}

	log.V(1).Info("populated Status object from Pingdom state")
	return nil
}

/*
updatePingdom updates the state of the check resource on Pingdom, creating or
updating it as necessary to bring it in line with the specification in
CheckSpec.
*/
func (r *CheckReconciler) updatePingdom(
	parentLog logr.Logger,
	check *observabilityv1alpha1.Check,
) (err error) {
	log := parentLog.WithValues("stage", "updatePingdom")
	if check.Status.Id == 0 {
		log.Info("creating check resource on Pingdom")
		pdResp, err := r.pdClient.Checks.Create(&check.Spec)
		log.V(1).Info("Pingdom Checks.Create() response", "response", pdResp, "error", err)
		if err == nil {
			check.Status.Id = int32(pdResp.ID)
			log.Info("created check resource on Pingdom", "id", pdResp.ID)
		}
	} else {
		// TODO: only update if necessary
		log = log.WithValues("id", check.Status.Id)
		log.Info("updating check resource on Pingdom")
		pdResp, err := r.pdClient.Checks.Update(int(check.Status.Id), &check.Spec)
		log.V(1).Info("Pingdom Checks.Update() response", "response", pdResp, "error", err)
		if err == nil {
			log.Info("updated check resource on Pingdom", "message", pdResp.Message)
		}
	}
	return
}

/*
deletePingdomCheck deletes the check resource from Pingdom if any.

If Status.Id is not set (0) - check doesn't exist on Pingdom so do nothing.
*/
func (r *CheckReconciler) deletePingdomCheck(
	parentLog logr.Logger,
	check *observabilityv1alpha1.Check,
) error {
	log := parentLog.WithValues("stage", "deletion")
	log.V(1).Info("deleting check resource from Pingdom")
	if check.Status.Id == 0 {
		log.Info("Pingdom resource doesn't exist, nothing to delete")
		return nil
	}
	pdResp, err := r.pdClient.Checks.Delete(int(check.Status.Id))
	if err == nil {
		log.Info(
			"deleted check resource from Pingdom",
			"id", check.Status.Id, "message", pdResp.Message,
		)
	} else {
		// TODO: check if error is 404 in which case do nothing
		log.Error(err, "unable to delete the Pingdom check resource")
	}
	return err
}

/*
prtI32 returns a pointer to a given int32 value
*/
func ptrI32(i int32) *int32 {
	return &i
}

/*
parsePdTime parses unix time value returned from Pingdom API into a metav1.Time
object.
*/
func parsePdTime(pdTime int64) *metav1.Time {
	t := metav1.Unix(pdTime, 0)
	return &t
}

func (r *CheckReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&observabilityv1alpha1.Check{}).
		Complete(r)
}
