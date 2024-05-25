
all: build test

run:
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
