package main

/*
  Got stuck at:
      The following builds contain text matching the pattern:
      https://example.com/teams//pipelines/cool-pipeline/jobs/cool-job-a/builds/1
      https://example.com/teams//pipelines/boring-pipeline/jobs/boring-job-a/builds/1

  Waiting for:
      The following build(s) contain the pattern:
      https://example.com/teams/my-team-1/pipelines/my-pipeline-a/jobs/fun-job/builds/420
      https://example.com/teams/my-team-2/pipelines/my-pipeline-a/jobs/boring-job/builds/817
      https://example.com/teams/my-team-3/pipelines/my-pipeline-b/jobs/cool-job/builds/999
*/
const pipelinesJSON = `[
	{
		"id": 450,
		"name": "cool-pipeline"
	},
	{
		"id": 943,
		"name": "boring-pipeline"
	}
]`

const coolBuildsJSON = `[
	{
		"id": 8719,
		"name": "2",
		"status": "succeeded",
		"job_name": "cool-job-a",
		"team_name": "my-team-1"
	},
	{
		"id": 9817,
		"name": "1",
		"status": "failed",
		"job_name": "cool-job-a",
		"team_name": "my-team-1"
	},
	{
		"id": 3781,
		"name": "1",
		"status": "failed",
		"job_name": "cool-job-z",
		"team_name": "my-team-1"
	}
]`

const boringBuildsJSON = `[
	{
		"id": 7816,
		"name": "2",
		"status": "succeeded",
		"job_name": "boring-job-c",
		"team_name": "my-team-4"
	},
	{
		"id": 47187,
		"name": "2",
		"status": "failed",
		"job_name": "boring-job-a",
		"team_name": "my-team-4"
	},
	{
		"id": 7788,
		"name": "5",
		"status": "failed",
		"job_name": "boring-job-a",
		"team_name": "my-team-4"
	}
]`
