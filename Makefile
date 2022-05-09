test:
	go test -v -race ./...
generate-sources:
	./_support/generate-sources.sh > fact/sources.go
	