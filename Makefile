
all: build test

run:
	echo "Running"
	@go run .

build:
	echo "Building..."
	@go build .

test:
	@go test ./...

clean:
	echo "Cleaning..."
	@rm go-chat chat.log
