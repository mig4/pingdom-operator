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

func (cr *checkReconciler) create() error {
	log := cr.log.WithValues("action", "create")
	log.Info("creating check resource on Pingdom")
	resp, err := cr.pdClient.Checks.Create(&cr.check.Spec)
	log.V(1).Info(
		"Pingdom Checks.Create() response", "response", resp, "error", err,
	)
	if err == nil {
		cr.check.Status.Id = int32(resp.ID)
		log.Info("created check resource on Pingdom", "id", resp.ID)
		cr.didWork = true
	}
	return err
}
