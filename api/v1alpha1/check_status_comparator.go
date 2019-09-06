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

package v1alpha1

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

/*
NeedsUpdate returns true if the external resource coresponding to this check
needs to be updated because its Spec differs from its Status.

The Spec reflects requested state of the resource while Status reflects its
state in Kubernetes.

Optional fields that aren't set in the Spec do not affect the comparison (in
this case the field in Status may still be set as it reflects the default value
returned from Pingdom API, but that's of no consequence).
*/
func (this *Check) NeedsUpdate() bool {
	spec := &this.Spec
	status := &this.Status

	if spec.Paused != nil {
		if *spec.Paused && status.Status != Paused {
			return true
		} else if !*spec.Paused && status.Status == Paused {
			return true
		}
	}

	opts := make(cmp.Options, 0)
	unspecifiedFields := make([]string, 0)

	if spec.Name == nil {
		unspecifiedFields = append(unspecifiedFields, "Name")
	}
	if spec.Port == nil {
		unspecifiedFields = append(unspecifiedFields, "Port")
	}
	if spec.ResolutionMinutes == nil {
		unspecifiedFields = append(unspecifiedFields, "ResolutionMinutes")
	}
	if spec.UserIds == nil {
		unspecifiedFields = append(unspecifiedFields, "UserIds")
	}
	if spec.Url == nil {
		unspecifiedFields = append(unspecifiedFields, "Url")
	}
	if spec.Encryption == nil {
		unspecifiedFields = append(unspecifiedFields, "Encryption")
	}

	if len(unspecifiedFields) > 0 {
		opts = append(
			opts,
			cmpopts.IgnoreFields(CheckParameters{}, unspecifiedFields...),
		)
	}

	return !cmp.Equal(spec.CheckParameters, status.CheckParameters, opts)
}
