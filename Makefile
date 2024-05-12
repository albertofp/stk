default: 
	@go run main.go $1

build:
	@go build -o bin/main main.go
