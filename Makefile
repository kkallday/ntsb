all: test build
test:
	dep ensure
	ginkgo -r -race -randomizeSuites --randomizeAllSpecs
build:
	dep ensure
	rm -rf ./output
	mkdir ./output
	go build -o output/ntsb
install:
	go install .
