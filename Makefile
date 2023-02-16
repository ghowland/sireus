BINARY_NAME=sireus

buildonly:
	@go get 2>/dev/null || (echo "Go packages up to date")
	@go build -o build/${BINARY_NAME} code/sireus.go || (echo "Build failed: $$?"; exit 1)
	@echo Build done: build/${BINARY_NAME}

run:
	@go build -o build/${BINARY_NAME} code/sireus.go || (echo "Build failed: $$?"; exit 1)
	@./build/${BINARY_NAME}

# Force everything to rebuild
force-build:
	@go build -a -o build/${BINARY_NAME} code/sireus.go || (echo "Build failed: $$?"; exit 1)

test:
	@go test ./code/...

format:
	@gofmt -l -s -w code

doc:
	@cd code ; ~/go/bin/gomarkdoc ./... > ../documentation/godoc.md

cov:
	@go test -coverprofile=coverage.txt -covermode=atomic ./code/...

coverage:
	$(MAKE) cov
	@go tool cover -html=coverage.txt

vet:
	@go vet ./code/

dockerbuild:
	docker build --build-arg PROJECT_VERSION=`cat version` -t ghowland/sireus:latest .

dockerrunsh:
	docker run -it --entrypoint /usr/bin/bash ghowland/sireus:latest

clean:
	go clean
	rm -f ./build/${BINARY_NAME}
	rm -f coverage.txt
	@echo Clean done

