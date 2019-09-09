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
func (cs *CheckSpec) PutParams() map[string]string {
	params := cs.PostParams()
	delete(params, "type")
	return params
}

/*
PostParams returns a map of parameters to send in POST requests (for creating
a check) to the Pingdom API.

Only non-nil properties of the CheckSpec are included in the resulting map.

Implements pingdom.Check interface.
*/
func (cs *CheckSpec) PostParams() map[string]string {
	if err := cs.Valid(); err != nil {
		return map[string]string{}
	}

	params := map[string]string{
		"name": *cs.Name,
		"type": string(cs.Type),
		"host": cs.Host,
	}

	if cs.Paused != nil {
		params["paused"] = strconv.FormatBool(*cs.Paused)
	}

	if cs.Port != nil {
		params["port"] = strconv.FormatInt(int64(*cs.Port), 10)
	}

	if cs.ResolutionMinutes != nil {
		params["resolution"] = strconv.FormatInt(int64(*cs.ResolutionMinutes), 10)
	}

	if cs.UserIds != nil {
		params["userids"] = intSliceToCommaSep(*cs.UserIds)
	}

	if cs.URL != nil {
		params["url"] = *cs.URL
	}

	if cs.Encryption != nil {
		params["encryption"] = strconv.FormatBool(*cs.Encryption)
	}

	return params
}

/*
Valid checks if this CheckSpec is valid in terms of parameters and values
accepted by the Pingdom API.

Implements pingdom.Check interface.
*/
func (cs *CheckSpec) Valid() error {
	if cs.Name == nil || *cs.Name == "" {
		return fmt.Errorf("check `Name` must be set and not empty")
	}

	if cs.Host == "" {
		return fmt.Errorf("check `Host` must not be empty")
	}

	switch cs.Type {
	case HTTP, HTTPCustom, TCP, Ping, DNS, UDP, SMTP, POP3, IMAP:
	default:
		return fmt.Errorf(
			"check `Type` must be one of: http, httpcustom, tcp, ping, dns, udp, smtp, pop3, imap",
		)
	}

	if cs.Port != nil && (*cs.Port < 1 || *cs.Port > 65535) {
		return fmt.Errorf("check `Port` must be between 1-65535")
	}

	if cs.ResolutionMinutes != nil && !isValidResolution(*cs.ResolutionMinutes) {
		return fmt.Errorf("check `ResolutionMinutes` must be one of 1, 5, 15, 30 or 60")
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
