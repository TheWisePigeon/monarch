build:
	@go build -o monarch .

test:
	@go test ./... -v

clean-cache:
	@go clean -testcache
