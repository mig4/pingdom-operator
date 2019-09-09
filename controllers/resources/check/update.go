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

func (cr *checkReconciler) update() error {
	log := cr.log.WithValues("action", "update", "id", cr.check.Status.ID)
	log.Info("updating check resource on Pingdom")
	resp, err := cr.pdClient.Checks.Update(int(cr.check.Status.ID), &cr.check.Spec)
	log.V(1).Info("Pingdom Checks.Update() response", "response", resp, "error", err)
	if err == nil {
		log.Info("updated check resource on Pingdom", "message", resp.Message)
		cr.didWork = true
	}
	return err
}
