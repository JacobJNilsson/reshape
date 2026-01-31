.PHONY: test format

test:
	go test ./...

format:
	gofmt -w ./cli ./internal
