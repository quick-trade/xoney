COVERAGE := testdata/coverage.txt
TESTPATH := ./...

setup :
	export PATH=$$PATH:$(shell go env GOPATH)/bin

lint :
	golangci-lint run --enable-all -D \
depguard,\
gci,\
varnamelen,\
gomnd,\
gofumpt,\
ifshort

fmt :
	go fmt ./... && \
	golangci-lint run --enable-all --fix

test-cov :
	go test $(TESTPATH) -v -coverprofile=$(COVERAGE) -coverpkg=./...

view-cov :
	go tool cover -html=$(COVERAGE)
