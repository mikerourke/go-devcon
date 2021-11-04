run:
	go run cmd/devcon/devcon.go

.PHONY: build-x86
build-x86:
	GOOS=windows GOARCH=386 go1.10.7 build -v -o out/devcon-test-32.exe cmd/devcon/devcon.go

.PHONY: build-x64
build-x64:
	GOOS=windows GOARCH=amd64 go build -v -o out/devcon-test-64.exe cmd/devcon/devcon.go

format:
	gofmt -w .

lint:
	golint .
	golangci-lint run

vet:
	go vet .