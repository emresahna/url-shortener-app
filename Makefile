build: build-http build-async

build-http:
	go build -o bin/http_server ./cmd/http_server

build-async:
	go build -o bin/async_server ./cmd/async_server

test:
	go test -v ./...

run-http:
	./bin/http_server

run-async:
	./bin/async_server

clean:
	rm -rf bin/

lint:
	golangci-lint run