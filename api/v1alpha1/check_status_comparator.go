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

import "reflect"

/*
NeedsUpdate returns true if this check resource needs to be updated in Kube
because its Spec differs from its Status.

The Spec reflects requested state of the resource while Status reflects its
state in Kubernetes.
*/
func (this *Check) NeedsUpdate() bool {
	if fieldNeedsUpdate(this.Spec.Name, this.Status.Name) {
		return true
	}
	if this.Spec.Host != this.Status.Host {
		return true
	}
	if this.Spec.Type != this.Status.Type {
		return true
	}
	if fieldNeedsUpdate(this.Spec.Paused, this.Status.Paused) {
		return true
	}
	if fieldNeedsUpdate(this.Spec.Port, this.Status.Port) {
		return true
	}
	if fieldNeedsUpdate(this.Spec.Url, this.Status.Url) {
		return true
	}
	if fieldNeedsUpdate(this.Spec.Encryption, this.Status.Encryption) {
		return true
	}
	return false
}

/*
fieldNeedsUpdate returns true if status of a field doesn't match what's
expected according to its spec.

This is determined as follows:

- if spec is a nil pointer, return false (it means no specific value was
  requested, status may still reflect a default but it is of no consequence)
- if both are nil return false
- if their dereferenced values are equal return false
- otherwise return true

Note: only works on pointers to primitive types; use on individual fields
*/
func fieldNeedsUpdate(specField, statusField interface{}) bool {
	if specField == nil {
		return false
	} else if statusField == nil {
		return true
	}
	return !equal(reflect.ValueOf(specField), reflect.ValueOf(statusField))
}

func equal(va, vb reflect.Value) bool {
	if !va.IsValid() || !vb.IsValid() {
		return false
	}
	if va.Type() != vb.Type() {
		return false
	}

	vaKind := va.Kind()
	vbKind := vb.Kind()
	derefA := vaKind == reflect.Ptr || vaKind == reflect.Interface
	derefB := vbKind == reflect.Ptr || vbKind == reflect.Interface

	if derefA || derefB {
		if derefA {
			va = va.Elem()
		}
		if derefB {
			vb = vb.Elem()
		}
		return equal(va, vb)
	}

	switch vaKind {
	case reflect.Float32, reflect.Float64:
		return va.Float() == vb.Float()
	case reflect.Bool:
		return va.Bool() == vb.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return va.Int() == vb.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return va.Uint() == vb.Uint()
	case reflect.String:
		return va.String() == vb.String()
	default:
		return false
	}
}
