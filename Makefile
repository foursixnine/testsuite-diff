build:
	go build -o testsuite-diff ./cmd/testsuite-diff

test:
	go fmt -x ./... && go test -v ./...