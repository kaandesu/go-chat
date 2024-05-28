
all: build test

run:
	@clear
	@echo "Running"
	@go run .

build:
	@echo "Building..."
	@go build .

test:
	@go test ./...

clean:
	@echo "Cleaning..."
	@if [ -f go-chat ]; then rm go-chat; fi
	@if [ -f chat.log ]; then rm chat.log; fi
	@if [ -f dist ]; then rm -rf dist; fi


local-release:
	@goreleaser release --snapshot --clean
