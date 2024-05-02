build:
	@go build -o ./bin/monarch .

test:
	@go test ./... -v
