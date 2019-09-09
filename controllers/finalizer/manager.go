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

package finalizer

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	observabilityv1alpha1 "gitlab.com/mig4/pingdom-operator/api/v1alpha1"
	"gitlab.com/mig4/pingdom-operator/controllers/resources"
)

type finalizerManager struct {
	log        logr.Logger
	client     client.Client
	reconciler resources.ResourceReconciler

	haveFinalizer bool
	finalizerName string
}

// New creates a new instance of finalizer Manager
func New(
	log logr.Logger,
	client client.Client,
	reconciler resources.ResourceReconciler,
) Manager {
	baseName := reconciler.FinalizerName()
	name := ""
	haveFinalizer := baseName != nil
	if haveFinalizer {
		name = finalizerName(baseName)
	}

	return &finalizerManager{
		log:        log.WithName("finalizer-manager").WithValues("finalizer", name),
		client:     client,
		reconciler: reconciler,

		haveFinalizer: haveFinalizer,
		finalizerName: name,
	}
}

func (fm *finalizerManager) EnsureAttached(ctx context.Context, obj runtime.Object) error {
	log := fm.log.WithValues("action", "ensureAttached")
	if !fm.haveFinalizer {
		log.V(1).Info("reconciler has no finalizer, skipping registration")
		return nil
	}

	objMeta, err := fm.accessor(obj)
	if err != nil {
		return microerror.Maskf(err, "cannot access ObjectMeta")
	}

	if sliceContains(objMeta.GetFinalizers(), fm.finalizerName) {
		return nil
	}

	log.Info("registering finalizer")
	objMeta.SetFinalizers(append(objMeta.GetFinalizers(), fm.finalizerName))
	if err := fm.updateObj(ctx, obj); err != nil {
		return microerror.Maskf(err, "unable to update object")
	}
	return nil
}

func (fm *finalizerManager) EnsureDetached(ctx context.Context, obj runtime.Object) error {
	log := fm.log.WithValues("action", "ensureDetached")
	if !fm.haveFinalizer {
		log.V(1).Info("reconciler has no finalizer, skipping de-registration")
		return nil
	}

	objMeta, err := fm.accessor(obj)
	if err != nil {
		return microerror.Maskf(err, "cannot access ObjectMeta")
	}

	if !sliceContains(objMeta.GetFinalizers(), fm.finalizerName) {
		return nil
	}

	log.Info("removing the finalizer")
	objMeta.SetFinalizers(sliceRemove(objMeta.GetFinalizers(), fm.finalizerName))
	if err := fm.updateObj(ctx, obj); err != nil {
		return microerror.Maskf(err, "unable to update object")
	}
	return nil
}

func (fm *finalizerManager) accessor(obj runtime.Object) (metav1.Object, error) {
	objMeta, err := meta.Accessor(obj)
	if err != nil {
		fm.log.Error(err, "invalid object, cannot access ObjectMeta", "obj", obj)
		return nil, err
	}
	return objMeta, nil
}

func (fm *finalizerManager) updateObj(ctx context.Context, obj runtime.Object) error {
	if err := fm.client.Update(ctx, obj); err != nil {
		fm.log.Error(err, "unable to update object", "obj", obj)
		return err
	}
	fm.log.Info("object updated", "obj", obj)
	return nil
}

func finalizerName(baseName *string) string {
	return observabilityv1alpha1.GroupVersion.Group + "/" + *baseName
}
