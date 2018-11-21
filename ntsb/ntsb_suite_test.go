package ntsb_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNtsb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ntsb Suite")
}
