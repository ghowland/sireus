BINARY_NAME=sireus

run:
	@go build -o build/${BINARY_NAME} code/sireus.go || (echo "Build failed: $$?"; exit 1)
	@./build/${BINARY_NAME}

# Force everything to rebuild
force-build:
	@go build -a -o build/${BINARY_NAME} code/sireus.go || (echo "Build failed: $$?"; exit 1)


test:
	@go test ./code/...

cov:
	@go test -coverprofile=coverage.txt -covermode=atomic ./code/...

coverage:
	$(MAKE) cov
	@go tool cover -html=coverage.txt

vet:
	@go vet ./code/

clean:
	go clean
	rm -f ./build/${BINARY_NAME}
	rm -f coverage.txt
	@echo Clean done

