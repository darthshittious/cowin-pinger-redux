clean:
	@echo "Getting Dependencies..."
	go get -d ./...

osx_64: clean
	@echo "Building "
	GOOS=darwin GOARCH=amd64 go build -o bin/main-osx ./main/main.go


osx_32: clean
	@echo "Port removed from Go. See https://github.com/golang/go/issues/37610"

linux_64: clean
	@echo "Building "
	GOOS=linux GOARCH=amd64 go build -o bin/main-linux ./main/main.go

linux_32: clean
	@echo "Building "
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 ./main/main.go

windows_64: clean
	@echo "Building "
	GOOS=windows GOARCH=amd64 go build -o bin/main-windows.exe ./main/main.go

windows_32: clean
	@echo "Building "
	GOOS=windows GOARCH=386 go build -o bin/main-windows-386.exe ./main/main.go

run:
	go run main/main.go
