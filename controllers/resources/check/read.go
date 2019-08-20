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
	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	observabilityv1alpha1 "gitlab.com/mig4/pingdom-operator/api/v1alpha1"
)

func (cr *checkReconciler) read() error {
	status := &cr.check.Status
	checkId := status.Id
	log := cr.log.WithValues("action", "read", "id", checkId)

	log.V(1).Info("fetching check resource from Pingdom")
	pdCheck, err := cr.pdClient.Checks.Read(int(checkId))
	if err != nil {
		log.Error(err, "unable to fetch check resource from Pingdom")
		return microerror.Maskf(err, "unable to fetch check resource from Pingdom")
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
			err = microerror.New("check type is http but details not available")
			log.Error(err, "Pingdom API didn't return http check details")
			return err
		}
		status.Port = ptrI32(int32(pdCheck.Type.HTTP.Port))
		status.Url = &pdCheck.Type.HTTP.Url
		status.Encryption = &pdCheck.Type.HTTP.Encryption
	} else if pdCheck.Type.Name == string(observabilityv1alpha1.Tcp) {
		if pdCheck.Type.TCP == nil {
			err = microerror.New("check type is tcp but details not available")
			log.Error(err, "Pingdom API didn't return tcp check details")
			return err
		}
		status.Port = ptrI32(int32(pdCheck.Type.TCP.Port))
	}

	log.V(1).Info("populated Status object from Pingdom state")
	return nil
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
