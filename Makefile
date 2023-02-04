run:
	@go build -o build/sireus code/sireus.go || (echo "Build failed: $$?"; exit 1)
	./build/sireus

test:
	@go test $(go list ./... | grep -v /example/)

cov:
	@go test -coverprofile=coverage.txt -covermode=atomic $(go list ./... | grep -v /example/)

coverage:
	$(MAKE) cov
	@go tool cover -html=coverage.txt
