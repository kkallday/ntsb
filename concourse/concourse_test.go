package concourse_test

import (
	"time"

	concoursefly "github.com/kkallday/ntsb/concourse"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Concourse", func() {
	Describe("BuildOutput", func() {
		var (
			concourse concoursefly.Concourse
		)

		BeforeEach(func() {
			concourse = concoursefly.New(pathToFakeFly, "some-target")
		})

		It("returns output for a build", func() {
			buildOutput, err := concourse.BuildOutput(971)
			Expect(err).NotTo(HaveOccurred())
			Expect(buildOutput).To(Equal("something\nsomething else\neven more\n"))
		})
	})

	Describe("Builds", func() {
		var (
			concourse concoursefly.Concourse
		)

		BeforeEach(func() {
			concourse = concoursefly.New(pathToFakeFly, "some-target")
		})

		It("returns builds for a pipeline", func() {
			builds, err := concourse.Builds("pipeline-a")
			Expect(err).NotTo(HaveOccurred())
			Expect(builds).To(Equal([]concoursefly.Build{
				{ID: 1234, Name: "2", Status: "succeeded", JobName: "job-a"},
				{ID: 4321, Name: "3", Status: "failed", JobName: "job-b"},
			}))
		})
	})

	Describe("Pipelines", func() {
		var (
			concourse concoursefly.Concourse
		)

		BeforeEach(func() {
			concourse = concoursefly.New(pathToFakeFly, "some-target")
		})

		It("returns pipelines", func() {
			pipelines, err := concourse.Pipelines()
			Expect(err).NotTo(HaveOccurred())
			Expect(pipelines).To(Equal([]concoursefly.Pipeline{
				{
					ID:       450,
					Name:     "pipeline-a",
					TeamName: "some-team-1",
				},
				{
					ID:       451,
					Name:     "pipeline-b",
					TeamName: "some-team-1",
				},
				{
					ID:       943,
					Name:     "pipeline-c",
					TeamName: "some-team-2",
				},
			}))
		})

	})

	Describe("TargetInfo", func() {
		var (
			concourse concoursefly.Concourse
		)

		BeforeEach(func() {
			concourse = concoursefly.New(pathToFakeFly, "target-1")
		})

		It("returns target info", func() {
			ti, err := concourse.TargetInfo()
			Expect(err).NotTo(HaveOccurred())

			expectedTime, err := time.Parse(time.RFC1123, "Thu, 25 Oct 2018 06:22:59 UTC")
			Expect(err).NotTo(HaveOccurred())

			Expect(ti).To(Equal(concoursefly.TargetInfo{
				Name:   "target-1",
				URL:    "https://example.com",
				Team:   "some-team",
				Expiry: expectedTime,
			}))
		})
	})
})
