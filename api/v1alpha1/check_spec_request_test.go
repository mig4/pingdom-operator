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

package v1alpha1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "gitlab.com/mig4/pingdom-operator/api/v1alpha1"
)

var _ = Describe("CheckSpecRequest Adapter", func() {
	var (
		spec   CheckSpec
		params map[string]string
	)

	// Common specifications logic
	AssertSuccess := func(paramsDesc string) {
		Specify("PutParams returns a "+paramsDesc, func() {
			// PutParams cannot include `Type` parameter, otherwise Pingdom API
			// returns error 400 (bad request)
			delete(params, "type")
			Expect(spec.PutParams()).To(Equal(params))
		})

		Specify("PostParams returns a "+paramsDesc, func() {
			Expect(spec.PostParams()).To(Equal(params))
		})

		Specify("Valid returns no error", func() {
			Expect(spec.Valid()).To(Succeed())
		})
	}

	AssertFailure := func(errSubstring string) {
		Specify("PutParams returns an empty map", func() {
			// PutParams cannot include `Type` parameter, otherwise Pingdom API
			// returns error 400 (bad request)
			delete(params, "type")
			Expect(spec.PutParams()).To(Equal(params))
		})

		Specify("PostParams returns an empty map", func() {
			Expect(spec.PostParams()).To(Equal(params))
		})

		Specify("Valid returns an error", func() {
			Expect(spec.Valid()).Should(MatchError(ContainSubstring(errSubstring)))
		})
	}

	Context("When the Spec is correct", func() {
		BeforeEach(func() {
			spec = CheckSpec{
				CheckParameters: CheckParameters{
					Name:       ptrS("foo"),
					Host:       "foo.example.com",
					Type:       Http,
					Paused:     ptrB(false),
					Port:       ptrI32(443),
					Url:        ptrS("/text"),
					Encryption: ptrB(true),
				},
			}
			params = map[string]string{
				"name":       "foo",
				"host":       "foo.example.com",
				"type":       "http",
				"paused":     "false",
				"port":       "443",
				"url":        "/text",
				"encryption": "true",
			}
		})

		AssertSuccess("map with all required properties")
	})

	Context("When the Spec contains unset optional values", func() {
		BeforeEach(func() {
			spec = CheckSpec{CheckParameters: CheckParameters{
				Name: ptrS("bar"),
				Host: "bar.example.com",
				Type: Ping,
			}}
			params = map[string]string{
				"name": "bar",
				"host": "bar.example.com",
				"type": "ping",
			}
		})

		AssertSuccess("map with specified properties only")
	})

	Context("When the Spec contains missing required values", func() {
		BeforeEach(func() {
			spec = CheckSpec{}
			params = map[string]string{}
		})

		AssertFailure("Check `Name` must be set")
	})

	Context("When the Spec contains invalid values", func() {
		BeforeEach(func() {
			spec = CheckSpec{CheckParameters: CheckParameters{
				Name: ptrS("bad-type-port"),
				Host: "bad.example.com",
				Type: CheckType("pong"),
				Port: ptrI32(10000),
			}}
			params = map[string]string{}
		})

		AssertFailure("Check `Type` must be one of")
	})
})
