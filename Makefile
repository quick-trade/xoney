COVERAGE := testdata/coverage.txt
TESTPATH := ./...

setup :
	export PATH=$$PATH:/home/vlad/go/bin

lint :
	golangci-lint run --enable-all -D \
depguard,\
gci,\
varnamelen,\
gomnd,\
gofumpt,\
ifshort,\
wrapcheck,\
paralleltest,\
ireturn

fmt :
	go fmt ./... && \
	golangci-lint run --enable-all --fix

test-cov :
	go test $(TESTPATH) -v -coverprofile=$(COVERAGE) -coverpkg=./...

view-cov :
	go tool cover -html=$(COVERAGE)

build :
	go build -v ./...

doc:
	@echo "Generating godoc documentation..."
	@godoc -http=:6060 &
