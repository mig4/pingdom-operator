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

/*
sliceContains returns true if given slice contains the given string, false
otherwise
*/
func sliceContains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

/*
sliceRemove removes given element from slice, returns the updated slice
*/
func sliceRemove(slice []string, element string) (result []string) {
	for _, item := range slice {
		if item == element {
			continue
		}
		result = append(result, item)
	}
	return
}
