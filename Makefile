run:
	go run cmd/devcon/devcon.go

build:
	GOOS=windows GOARCH=386 go1.10.7 build -v -o out/devcon.exe cmd/devcon/devcon.go