build: scriptulus.go
	go build -ldflags "-X main.root '`pwd`'" -o scriptulus scriptulus.go
