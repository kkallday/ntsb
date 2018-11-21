package ntsb_test

import (
	"time"

	concoursefly "github.com/kkallday/ntsb/concourse"
	"github.com/kkallday/ntsb/ntsb"
	"github.com/kkallday/ntsb/ntsb/ntsbfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

/*
	fly -t raas-todd pipelines -a --json
		fly -t raas-todd bs -p build-tile-test-symlink-sink-resources --json
			fly -t raas-todd watch -b BID
*/

var _ = Describe("App", func() {
	var (
		concourse *ntsbfakes.FakeConcourse
		logger    *ntsbfakes.FakeLogger

		app ntsb.App
	)

	BeforeEach(func() {
		concourse = &ntsbfakes.FakeConcourse{}
		logger = &ntsbfakes.FakeLogger{}

		concourse.PipelinesReturns([]concoursefly.Pipeline{{Name: "pipeline-0"}, {Name: "pipeline-1"}, {Name: "pipeline-2"}, {Name: "pipeline-3"}}, nil)

		concourse.BuildsReturnsOnCall(0, []concoursefly.Build{
			{ID: 1111, Name: "1", Status: "failed", JobName: "job-1a", TeamName: "team-a"},
			{ID: 1112, Name: "2", Status: "failed", JobName: "job-1b", TeamName: "team-a"},
			{ID: 1113, Name: "3", Status: "succeeded", JobName: "job-1a", TeamName: "team-a"},
		}, nil)

		concourse.BuildsReturnsOnCall(1, []concoursefly.Build{
			{ID: 2221, Name: "1", Status: "failed", JobName: "job-2a", TeamName: "team-b"},
			{ID: 2222, Name: "6", Status: "succeeded", JobName: "job-2h", TeamName: "team-b"},
			{ID: 2223, Name: "7", Status: "failed", JobName: "job-2a", TeamName: "team-b"},
		}, nil)

		concourse.BuildsReturnsOnCall(2, []concoursefly.Build{
			{ID: 3331, Name: "1", Status: "succeeded", JobName: "job-3a", TeamName: "team-c"},
			{ID: 3332, Name: "3", Status: "succeeded", JobName: "job-3h", TeamName: "team-c"},
			{ID: 3333, Name: "4", Status: "succeeded", JobName: "job-3a", TeamName: "team-c"},
		}, nil)

		concourse.BuildOutputReturnsOnCall(0, "adfafkthis-is-a-patternafk8fahdjf", nil)
		concourse.BuildOutputReturnsOnCall(1, "adfafkthis-is-a-patternafk8fahdjf", nil)
		concourse.BuildOutputReturnsOnCall(2, "adfafkthis-is-a-patternafk8fahdjf", nil)
		concourse.BuildOutputReturnsOnCall(3, "doesnt-match", nil)

		concourse.TargetInfoReturns(concoursefly.TargetInfo{
			Name:   "my-team-target",
			URL:    "https://example.com",
			Expiry: time.Now(),
			Team:   "some-random-team",
		}, nil)

		app = ntsb.NewApp(concourse, logger)
	})

	Describe("Run", func() {
		It("finds builds that contain text matching given pattern", func() {
			err := app.Run("this-is-a-pattern")
			Expect(err).ToNot(HaveOccurred(), "Run returned an error")

			Expect(concourse.TargetInfoCallCount()).To(Equal(1), "TargetInfo call count")

			By("retrieving a list of all pipelines", func() {
				Expect(concourse.PipelinesCallCount()).To(Equal(1), "Pipelines call count")
			})

			By("retrieving a list of all builds", func() {
				Expect(concourse.BuildsCallCount()).To(Equal(4), "builds were not retrieved for each pipeline")

				pipelineName := concourse.BuildsArgsForCall(0)
				Expect(pipelineName).To(Equal("pipeline-0"), "builds retrieved from wrong pipeline on 1st call")

				pipelineName = concourse.BuildsArgsForCall(1)
				Expect(pipelineName).To(Equal("pipeline-1"), "builds retrieved from wrong pipeline on 2nd call")

				pipelineName = concourse.BuildsArgsForCall(2)
				Expect(pipelineName).To(Equal("pipeline-2"), "builds retrieved from wrong pipeline on 3rd call")

				pipelineName = concourse.BuildsArgsForCall(3)
				Expect(pipelineName).To(Equal("pipeline-3"), "builds retrieved from wrong pipeline on 4th call")
			})

			By("retrieving output of each build in each job of each pipeline", func() {
				Expect(concourse.BuildOutputCallCount()).To(Equal(4), "BuildOutput call count")

				bid := concourse.BuildOutputArgsForCall(0)
				Expect(bid).To(Equal(1111), "build output retrieved for wrong build ID on 1st call")

				bid = concourse.BuildOutputArgsForCall(1)
				Expect(bid).To(Equal(1112), "build output retrieved for wrong build ID on 2nd call")

				bid = concourse.BuildOutputArgsForCall(2)
				Expect(bid).To(Equal(2221), "build output retrieved for wrong build ID on 3rd call")

				bid = concourse.BuildOutputArgsForCall(3)
				Expect(bid).To(Equal(2223), "build output retrieved for wrong build ID on 4th call")
			})

			By("printing a list of each build containing the pattern", func() {
				Expect(logger.PrintlnCallCount()).To(Equal(4), "Println call count")

				args := logger.PrintlnArgsForCall(0)
				Expect(args).To(HaveLen(1), "incorrect number of args to Println")
				actualMsg, ok := args[0].(string)
				Expect(ok).To(BeTrue(), "failed to type assert 1st arg to Println")
				Expect(actualMsg).To(Equal("The following build(s) contain text matching the pattern:"), "did not print header message")

				args = logger.PrintlnArgsForCall(1)
				Expect(args).To(HaveLen(1), "incorrect number of args to Println on 1st call")
				actualMsg, ok = args[0].(string)
				Expect(ok).To(BeTrue(), "failed to type assert 1st arg to Println")
				Expect(actualMsg).To(Equal("https://example.com/teams/team-a/pipelines/pipeline-0/jobs/job-1a/builds/1"), "did not print 1st matching failed build")

				args = logger.PrintlnArgsForCall(2)
				Expect(args).To(HaveLen(1), "incorrect number of args to Println")
				actualMsg, ok = args[0].(string)
				Expect(ok).To(BeTrue(), "failed to type assert 1st arg to Println")
				Expect(actualMsg).To(Equal("https://example.com/teams/team-a/pipelines/pipeline-0/jobs/job-1b/builds/2"), "did not print 2nd matching failed build")

				args = logger.PrintlnArgsForCall(3)
				Expect(args).To(HaveLen(1), "incorrect number of args to Println")
				actualMsg, ok = args[0].(string)
				Expect(ok).To(BeTrue(), "failed to type assert 1st arg to Println")
				Expect(actualMsg).To(Equal("https://example.com/teams/team-b/pipelines/pipeline-1/jobs/job-2a/builds/1"), "did not print 3rd matching failed build")
			})
		})

	})
})
