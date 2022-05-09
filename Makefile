test:
	go test -v -race ./...
generate-sources:
	go run _support/generate-sources.go _support/sources.json > fact/sources.go
	