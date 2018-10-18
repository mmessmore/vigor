all: build

build: vigor.exe vigor.darwin vigor.linux vigor.freebsd

vigor.exe:
	GOOS=windows GOARCH=amd64 go build

vigor.linux:
	GOOS=linux GOARCH=amd64 go build
	mv vigor vigor.linux

vigor.darwin:
	GOOS=darwin GOARCH=amd64 go build
	mv vigor vigor.darwin

vigor.freebsd:
	GOOS=darwin GOARCH=amd64 go build
	mv vigor vigor.freebsd

clean:
	rm -f vigor vigor.linux vigor.darwin vigor.freebsd vigor.exe

@PHONY: build clean all
