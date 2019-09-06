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
	"fmt"
	"strconv"
	"strings"
)

/*
PutParams returns a map of parameters to send in PUT requests (for updating
a check) to the Pingdom API.

Only non-nil properties of the CheckSpec are included in the resulting map.

Implements pingdom.Check interface.

TODO: moving these methods to Check would allow e.g. PUT to be smarter and
  only output fields that have changed, or handling clearing of fields
*/
func (this *CheckSpec) PutParams() map[string]string {
	params := this.PostParams()
	delete(params, "type")
	return params
}

/*
PostParams returns a map of parameters to send in POST requests (for creating
a check) to the Pingdom API.

Only non-nil properties of the CheckSpec are included in the resulting map.

Implements pingdom.Check interface.
*/
func (this *CheckSpec) PostParams() map[string]string {
	if err := this.Valid(); err != nil {
		return map[string]string{}
	}

	params := map[string]string{
		"name": *this.Name,
		"type": string(this.Type),
		"host": this.Host,
	}

	if this.Paused != nil {
		params["paused"] = strconv.FormatBool(*this.Paused)
	}

	if this.Port != nil {
		params["port"] = strconv.FormatInt(int64(*this.Port), 10)
	}

	if this.ResolutionMinutes != nil {
		params["resolution"] = strconv.FormatInt(int64(*this.ResolutionMinutes), 10)
	}

	if this.UserIds != nil {
		params["userids"] = intSliceToCommaSep(*this.UserIds)
	}

	if this.Url != nil {
		params["url"] = *this.Url
	}

	if this.Encryption != nil {
		params["encryption"] = strconv.FormatBool(*this.Encryption)
	}

	return params
}

/*
Valid checks if this CheckSpec is valid in terms of parameters and values
accepted by the Pingdom API.

Implements pingdom.Check interface.
*/
func (this *CheckSpec) Valid() error {
	if this.Name == nil || *this.Name == "" {
		return fmt.Errorf("Check `Name` must be set and not empty")
	}

	if this.Host == "" {
		return fmt.Errorf("Check `Host` must not be empty")
	}

	switch this.Type {
	case Http, HttpCustom, Tcp, Ping, Dns, Udp, Smtp, Pop3, Imap:
	default:
		return fmt.Errorf(
			"Check `Type` must be one of: http, httpcustom, tcp, ping, dns, udp, smtp, pop3, imap",
		)
	}

	if this.Port != nil && (*this.Port < 1 || *this.Port > 65535) {
		return fmt.Errorf("Check `Port` must be between 1-65535")
	}

	if this.ResolutionMinutes != nil && !isValidResolution(*this.ResolutionMinutes) {
		return fmt.Errorf("Check `ResolutionMinutes` must be one of 1, 5, 15, 30 or 60")
	}

	return nil
}

func intSliceToCommaSep(intSlice []int) string {
	return strings.Join(intSliceToStrSlice(intSlice), ",")
}

func intSliceToStrSlice(intSlice []int) (result []string) {
	for _, item := range intSlice {
		result = append(result, strconv.Itoa(item))
	}
	return
}

func isValidResolution(res int32) bool {
	return res == 1 || res == 5 || res == 15 || res == 30 || res == 60
}
