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

var _ = PDescribe("CheckStatusComparator", func() {
	DescribeTable("checking if update is needed",
		func(check *Check, match types.GomegaMatcher) {
			Expect(check.NeedsUpdate()).To(match)
		},
		Entry("with no difference", &Check{
			Spec:   CheckSpec{Name: ptrS("foo"), Host: "foo", Type: Ping},
			Status: CheckStatus{Id: 1, Name: ptrS("foo"), Host: "foo", Type: Ping},
		}, BeFalse()),
		Entry("with no diff but more details in status", &Check{
			Spec: CheckSpec{Name: ptrS("bar"), Host: "bar", Type: Http},
			Status: CheckStatus{
				Id: 2, Name: ptrS("bar"), Host: "bar", Type: Http,
				Port: ptrI32(80), Url: ptrS("/"),
			},
		}, BeFalse()),
		Entry("with more values in spec", &Check{
			Spec: CheckSpec{
				Name: ptrS("baz"), Host: "baz", Type: Http,
				Port: ptrI32(8081), Url: ptrS("/path"),
			},
			Status: CheckStatus{Name: ptrS("baz"), Host: "baz", Type: Http},
		}, BeTrue()),
		// TODO: check if PD API allows the type to be changed
		Entry("with different host", &Check{
			Spec:   CheckSpec{Name: ptrS("qux"), Host: "qux", Type: Ping},
			Status: CheckStatus{Id: 3, Name: ptrS("qux"), Host: "foo", Type: Ping},
		}, BeTrue()),
		Entry("with different port", &Check{
			Spec: CheckSpec{
				Name: ptrS("quux"), Host: "quux", Type: Http, Port: ptrI32(8080),
			},
			Status: CheckStatus{
				Name: ptrS("quux"), Host: "quux", Type: Http, Port: ptrI32(80),
			},
		}, BeTrue()),
		Entry("with different paused setting", &Check{
			Spec: CheckSpec{
				Name: ptrS("quuz"), Host: "quuz", Type: Http, Paused: ptrB(true),
			},
			Status: CheckStatus{
				Name: ptrS("quuz"), Host: "quuz", Type: Http, Paused: ptrB(false),
			},
		}, BeTrue()),
		Entry("with different URL", &Check{
			Spec: CheckSpec{
				Name: ptrS("corge"), Host: "corge", Type: Http, Url: ptrS("/text"),
			},
			Status: CheckStatus{
				Name: ptrS("corge"), Host: "corge", Type: Http, Url: ptrS("/"),
			},
		}, BeTrue()),
		Entry("with different encryption setting", &Check{
			Spec: CheckSpec{
				Name: ptrS("grault"), Host: "grault", Type: Http, Encryption: ptrB(true),
			},
			Status: CheckStatus{
				Name: ptrS("grault"), Host: "grault", Type: Http, Encryption: ptrB(false),
			},
		}, BeTrue()),
	)
})
