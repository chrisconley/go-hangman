.PHONY: test
test:
	go clean -testcache && go test ./...

stest:
	go clean -testcache && go test -short ./...

vtest:
	go clean -testcache && go test -v ./...

build:
	go build -o ./player ./cmd/player.go