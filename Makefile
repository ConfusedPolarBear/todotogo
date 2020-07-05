all: test build

build:
	go build -o todo

test:
	go test

clean:
	go clean
