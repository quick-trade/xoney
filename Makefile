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
gci

fmt :
	go fmt ./... && golangci-lint run --enable-all --fix
