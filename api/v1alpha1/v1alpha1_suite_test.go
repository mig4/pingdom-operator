package v1alpha1_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestV1alpha1(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "V1alpha1 Suite")
}

func ptrB(b bool) *bool {
	return &b
}

func ptrI32(i int32) *int32 {
	return &i
}

func ptrS(s string) *string {
	return &s
}
