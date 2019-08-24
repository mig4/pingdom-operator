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

package check

import (
	"context"
)

var Name = "check-resource"

func (cr *checkReconciler) RefreshState(ctx context.Context) error {
	status := &cr.check.Status
	checkId := status.Id
	log := cr.log.WithValues("action", "refreshState", "id", checkId)
	if cr.check.Status.Id == 0 {
		log.Info("Pingdom resource doesn't exist yet, nothing to refresh")
		return nil
	}
	return cr.read()
}

func (cr *checkReconciler) EnsureState(ctx context.Context) (err error) {
	log := cr.log.WithValues("action", "ensureState")
	log.Info("entered reconciling external resource state")

	if !cr.check.GetDeletionTimestamp().IsZero() {
		err = cr.delete()
	} else if cr.check.Status.Id == 0 {
		err = cr.create()
	} else if cr.check.NeedsUpdate() {
		err = cr.update()
	} else {
		log.V(1).Info(
			"check is up-to-date with regards to its spec", "id", cr.check.Status.Id,
		)
	}

	log.Info("finished reconciling external resource state")
	return err
}

func (cr *checkReconciler) FinalizerName() *string {
	return &Name
}
