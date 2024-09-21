build:
	@go build -o bin/iam_service

run: build
	./bin/iam_service

