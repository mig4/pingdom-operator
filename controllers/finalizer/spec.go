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

	"k8s.io/apimachinery/pkg/runtime"
)

/*
FinalizerManager manages resistration and de-registration of a finalizer on an
API object in an idempotent way.
*/
type FinalizerManager interface {
	/*
	   EnsureAttached ensures a finalizer for the configured reconciler is
	   registered on the API object.

	   If it's not it will add it and update the object.
	*/
	EnsureAttached(context.Context, runtime.Object) error

	/*
	   EnsureDetached ensures a finalizer for the configured reconciler is
	   not registered on the API object.

	   If it is it will remove it and update the object.
	*/
	EnsureDetached(context.Context, runtime.Object) error
}
