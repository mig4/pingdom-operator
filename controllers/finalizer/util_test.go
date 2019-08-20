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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("Util", func() {
	WithMatchingElements := func(expected interface{}) TableEntry {
		return Entry(
			"with matching elements",
			[]string{"foo", "bar", "foo", "baz"}, "foo", expected,
		)
	}
	WithoutMatchingElements := func(expected interface{}) TableEntry {
		return Entry(
			"without matching elements",
			[]string{"foo", "bar"}, "baz", expected,
		)
	}
	WithEmptySlice := func(expected interface{}) TableEntry {
		return Entry(
			"with empty slice",
			[]string{}, "foo", expected,
		)
	}
	WithNilSlice := func(expected interface{}) TableEntry {
		return Entry(
			"with nil slice",
			[]string(nil), "foo", expected,
		)
	}

	Describe("sliceContains", func() {
		DescribeTable("checking if element is in slice",
			func(slice []string, element string, match types.GomegaMatcher) {
				Expect(sliceContains(slice, element)).To(match)
			},
			WithMatchingElements(BeTrue()),
			WithoutMatchingElements(BeFalse()),
			WithEmptySlice(BeFalse()),
			WithNilSlice(BeFalse()),
		)
	})

	Describe("sliceRemove", func() {
		DescribeTable("removing an element from slice",
			func(slice []string, element string, expected []string) {
				Expect(sliceRemove(slice, element)).To(Equal(expected))
			},
			WithMatchingElements([]string{"bar", "baz"}),
			WithoutMatchingElements([]string{"foo", "bar"}),
			WithEmptySlice([]string(nil)),
			WithNilSlice([]string(nil)),
		)
	})
})
