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

package resources

import "context"

/*
ResourceReconciler can ensure an external resource matches desired state
(Spec) and inspect the external resource and update internal representation
(Status).

Note these should be request-scoped (new instances created per-request) in
order to avoid state leaking between requests.
*/
type ResourceReconciler interface {
	// RefreshState reads the status of the external resource and updates the
	// internal Status resource on the API object.
	RefreshState(context.Context) error

	// Inspect the difference between Spec and current Status and update the
	// external resource as necessary (creating, updating or deleting it).
	EnsureState(context.Context) error

	// FinalizerName returns a base name of the finalizer for this reconciler
	// or nil if this reconciler doesn't do finalization
	FinalizerName() *string

	// DidWork indicates if the reconciler actually did any work in its
	// `EnsureState` method. The reconciler may determine that it has nothing
	// to do, e.g. if the external resource already matches the desired state,
	// in which case this will return false.
	DidWork() bool
}
