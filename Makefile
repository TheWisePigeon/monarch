build:
	@go build -o monarch .

test:
	@go test ./... -v
