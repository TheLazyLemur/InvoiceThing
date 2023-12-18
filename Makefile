run:
	@templ generate
	@go run cmd/main.go

build: 
	@templ generate
	@go build -o bin/invoicething cmd/main.go
