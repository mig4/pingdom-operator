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

import "github.com/russellcardullo/go-pingdom/pingdom"

// An error the Pingdom API returns when a given ID is not found
var invalidIdentifierError = pingdom.PingdomError{
	StatusCode: 403,
	StatusDesc: "Forbidden",
	Message:    "Invalid check identifier",
}

// IsInvalidIdentifierError returns true if given error returned by the Pingdom
// API indicates a given ID was not found.
func IsInvalidIdentifierError(err error) bool {
	if err == nil {
		return false
	}
	switch t := err.(type) {
	case *pingdom.PingdomError:
		return (t.StatusCode == invalidIdentifierError.StatusCode &&
			t.StatusDesc == invalidIdentifierError.StatusDesc &&
			t.Message == invalidIdentifierError.Message)
	default:
		return false
	}
}
