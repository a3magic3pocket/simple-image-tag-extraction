BINARY_NAME=extract-image-tag

build:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}_linux_amd64 main.go

run:
	./${BINARY_NAME}

build_and_run: build run

clean:
	go clean
	rm ${BINARY_NAME}_linux_amd64