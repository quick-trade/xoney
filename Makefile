setup :
	export PATH=$$PATH:$(shell go env GOPATH)/bin

lint :
	golangci-lint run --enable-all -D \
deadcode,\
golint,\
scopelint,\
ifshort,\
nosnakecase,\
maligned,\
exhaustivestruct,\
interfacer,\
varcheck,\
structcheck,\
depguard,\
gci,\
gofumpt,\
wrapcheck,\
varnamelen

fmt :
	go fmt ./... && golangci-lint run --enable-all --fix

test-cov :
	go test ./test/... -v -coverprofile="testdata/coverage.txt" -coverpkg=./...

view-cov :
	go tool cover -html="testdata/coverage.txt"
