build:
	@go build -o bin/iam_service

run: build
	./bin/iam_service

release:
	CGO_ENABLED=0 go build -ldflags "-s -w" -o /release/iam
