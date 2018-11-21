package main_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var pathToMain string

func TestNtsb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ntsb Acceptance Suite")
}

var _ = BeforeSuite(func() {
	var err error
	pathToMain, err = gexec.Build("github.com/kkallday/ntsb")
	Expect(err).NotTo(HaveOccurred())

	pathToFakeFly, err := gexec.Build("github.com/kkallday/ntsb/fakes/fly")
	Expect(err).NotTo(HaveOccurred())

	currPath := os.Getenv("PATH")
	newPath := fmt.Sprintf("%s:%s", filepath.Dir(pathToFakeFly), currPath)
	err = os.Setenv("PATH", newPath)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
