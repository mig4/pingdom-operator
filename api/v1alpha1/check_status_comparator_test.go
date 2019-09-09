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
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"

	. "gitlab.com/mig4/pingdom-operator/api/v1alpha1"
)

var _ = Describe("CheckStatusComparator", func() {
	DescribeTable("checking if update is needed",
		func(check *Check, match types.GomegaMatcher) {
			Expect(check.NeedsUpdate()).To(match)
		},
		Entry("with no difference", &Check{
			Spec: CheckSpec{CheckParameters: CheckParameters{
				Name: ptrS("foo"), Host: "foo", Type: Ping,
			}},
			Status: CheckStatus{ID: 1, CheckParameters: CheckParameters{
				Name: ptrS("foo"), Host: "foo", Type: Ping,
			}},
		}, BeFalse()),
		Entry("with no diff but more details in status", &Check{
			Spec: CheckSpec{CheckParameters: CheckParameters{
				Name: ptrS("bar"), Host: "bar", Type: HTTP,
			}},
			Status: CheckStatus{ID: 2, CheckParameters: CheckParameters{
				Name: ptrS("bar"), Host: "bar", Type: HTTP,
				ResolutionMinutes: ptrI32(1), Port: ptrI32(80), URL: ptrS("/"),
			}},
		}, BeFalse()),
		Entry("with more values in spec", &Check{
			Spec: CheckSpec{CheckParameters: CheckParameters{
				Name: ptrS("baz"), Host: "baz", Type: HTTP,
				Port: ptrI32(8081), URL: ptrS("/path"),
			}},
			Status: CheckStatus{ID: 3, CheckParameters: CheckParameters{
				Name: ptrS("baz"), Host: "baz", Type: HTTP,
			}},
		}, BeTrue()),
		Entry("with different host", &Check{
			Spec: CheckSpec{CheckParameters: CheckParameters{
				Name: ptrS("qux"), Host: "qux", Type: Ping,
			}},
			Status: CheckStatus{ID: 4, CheckParameters: CheckParameters{
				Name: ptrS("qux"), Host: "foo", Type: Ping,
			}},
		}, BeTrue()),
		Entry("with different port", &Check{
			Spec: CheckSpec{CheckParameters: CheckParameters{
				Name: ptrS("quux"), Host: "quux", Type: HTTP, Port: ptrI32(8080),
			}},
			Status: CheckStatus{ID: 5, CheckParameters: CheckParameters{
				Name: ptrS("quux"), Host: "quux", Type: HTTP, Port: ptrI32(80),
			}},
		}, BeTrue()),
		Entry("with spec paused but status up", &Check{
			Spec: CheckSpec{CheckParameters: CheckParameters{
				Name: ptrS("quuz"), Host: "quuz", Type: HTTP,
			}, Paused: ptrB(true)},
			Status: CheckStatus{ID: 6, CheckParameters: CheckParameters{
				Name: ptrS("quuz"), Host: "quuz", Type: HTTP,
			}, Status: Up},
		}, BeTrue()),
		Entry("with spec not paused but status paused", &Check{
			Spec: CheckSpec{CheckParameters: CheckParameters{
				Name: ptrS("corge"), Host: "corge", Type: HTTP,
			}, Paused: ptrB(false)},
			Status: CheckStatus{ID: 7, CheckParameters: CheckParameters{
				Name: ptrS("corge"), Host: "corge", Type: HTTP,
			}, Status: Paused},
		}, BeTrue()),
		Entry("with different resolutions", &Check{
			Spec: CheckSpec{CheckParameters: CheckParameters{
				Name: ptrS("dres"), Host: "dres", Type: Ping, ResolutionMinutes: ptrI32(1),
			}},
			Status: CheckStatus{ID: 7, CheckParameters: CheckParameters{
				Name: ptrS("dres"), Host: "dres", Type: Ping, ResolutionMinutes: ptrI32(5),
			}},
		}, BeTrue()),
		Entry("with different User IDs", &Check{
			Spec: CheckSpec{CheckParameters: CheckParameters{
				Name: ptrS("grault"), Host: "grault", Type: Ping, UserIds: &[]int{22},
			}},
			Status: CheckStatus{ID: 7, CheckParameters: CheckParameters{
				Name: ptrS("grault"), Host: "grault", Type: Ping, UserIds: &[]int{2},
			}},
		}, BeTrue()),
		Entry("with different URL", &Check{
			Spec: CheckSpec{CheckParameters: CheckParameters{
				Name: ptrS("garply"), Host: "garply", Type: HTTP, URL: ptrS("/text"),
			}},
			Status: CheckStatus{ID: 8, CheckParameters: CheckParameters{
				Name: ptrS("garply"), Host: "garply", Type: HTTP, URL: ptrS("/"),
			}},
		}, BeTrue()),
		Entry("with different encryption setting", &Check{
			Spec: CheckSpec{CheckParameters: CheckParameters{
				Name: ptrS("waldo"), Host: "waldo", Type: HTTP, Encryption: ptrB(true),
			}},
			Status: CheckStatus{ID: 9, CheckParameters: CheckParameters{
				Name: ptrS("waldo"), Host: "waldo", Type: HTTP, Encryption: ptrB(false),
			}},
		}, BeTrue()),
		Entry("with no difference with all parameters", &Check{
			Spec: CheckSpec{CheckParameters: CheckParameters{
				Name: ptrS("fred"), Host: "fred", Type: HTTP, Port: ptrI32(443),
				ResolutionMinutes: ptrI32(1), UserIds: &[]int{42, 24},
				URL: ptrS("/text"), Encryption: ptrB(true),
			}, Paused: ptrB(true)},
			Status: CheckStatus{ID: 10, CheckParameters: CheckParameters{
				Name: ptrS("fred"), Host: "fred", Type: HTTP, Port: ptrI32(443),
				ResolutionMinutes: ptrI32(1), UserIds: &[]int{42, 24},
				URL: ptrS("/text"), Encryption: ptrB(true),
			}, Status: Paused},
		}, BeFalse()),
	)
})
