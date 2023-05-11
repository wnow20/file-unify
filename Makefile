build:
	go build -ldflags="-s -w" -o bin/unify main.go
	$(if $(shell command -v upx), upx goctl)

mac:
	GOOS=darwin go build -ldflags="-s -w" -o bin/unify-darwin main.go
	$(if $(shell command -v upx), upx goctl-darwin)

win:
	GOOS=windows go build -ldflags="-s -w" -o bin/unify.exe main.go
	$(if $(shell command -v upx), upx goctl.exe)

linux:
	GOOS=linux go build -ldflags="-s -w" -o bin/unify-linux main.go
	$(if $(shell command -v upx), upx goctl-linux)