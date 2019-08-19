package controllers

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
