test:
	go test -v -race ./...
generate-sources:
	rm -rf _support/sources/
	echo "Cloning last sources"
	git clone git@github.com:joinself/sources.git _support/sources
	echo "Generating go code"
	go run _support/generate-sources.go _support/sources/sources.json > fact/sources.go
	