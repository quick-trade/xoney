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
gci --fix

fmt :
	go fmt ./...
