build:
	@go build -o ./yumfoods ./cmd/.

test:
	@go test -v ./...

run: build
	@./yumfoods
