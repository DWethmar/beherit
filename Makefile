BINARY_NAME=beherit

build:
	go build -o bin/${BINARY_NAME} main.go

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm bin/${BINARY_NAME}