build: test
	go build -o $(GOPATH)/bin/wishgen main/main.go

test:
	go test ./... -v