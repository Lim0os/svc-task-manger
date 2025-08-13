.PHONY: build run test

build:
	go build -o bin/svc-task_master main.go

run:
	go run main.go

test:
	go test ./...