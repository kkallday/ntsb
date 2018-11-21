package main_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("ntsb", func() {
	It("prints a list of builds that contain text matching provided regex", func() {
		cmd := exec.Command(pathToMain,
			"--target", "target-1",
			"--pattern", "this-is-a-pattern",
		)

		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		Expect(string(session.Out.Contents())).To(Equal(`The following build(s) contain text matching the pattern:
https://example.com/teams/my-team-1/pipelines/cool-pipeline/jobs/cool-job-a/builds/1
https://example.com/teams/my-team-1/pipelines/cool-pipeline/jobs/cool-job-z/builds/1
https://example.com/teams/my-team-4/pipelines/boring-pipeline/jobs/boring-job-a/builds/2
`))
	})
})
