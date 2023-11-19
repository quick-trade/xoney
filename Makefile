COVERAGE := testdata/coverage.txt
TESTPATH := ./test/...

setup :
	export PATH=$$PATH:$(shell go env GOPATH)/bin

lint :
	golangci-lint run --enable-all -D \
depguard,\
gci,\
varnamelen,\
gomnd,\
gofumpt

fmt :
	go fmt -s && \
	golangci-lint run --enable-all --fix

test-cov :
	go test $(TESTPATH) -v -coverprofile=$(COVERAGE) -coverpkg=./...

view-cov :
	go tool cover -html=$(COVERAGE)
