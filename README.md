# ntsb

A CLI for investigating build failures in [Concourse](https://concourse-ci.org). Like `grep` but for
Concourse builds.

(inspired by Concourse aeronautics theme and
[this](https://en.wikipedia.org/wiki/National_Transportation_Safety_Board))

# Usage

```bash
$ ntsb --target my-team-target --pattern 'TLS handshake timeout'
The following build(s) contain text matching the pattern:
https://concourse.example.com/teams/my-team-1/pipelines/cool-pipeline/jobs/cool-job-a/builds/1
https://concourse.example.com/teams/my-team-1/pipelines/cool-pipeline/jobs/cool-job-z/builds/1
https://concourse.example.com/teams/my-team-4/pipelines/boring-pipeline/jobs/boring-job-a/builds/2
```

# Build and Testing

To run tests (requires [`ginkgo`](https://github.com/onsi/ginkgo)):
```bash
$ make test
```

To build:
```bash
$ make build
```

To install `ntsb` on $PATH:
```bash
$ make install
```
