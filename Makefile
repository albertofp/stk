default: 
	@go run main.go $1

build:
	@go build -o bin/stk main.go
	@cp bin/stk ~/bin
