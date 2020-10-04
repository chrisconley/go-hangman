.PHONY: test
test:
	go test ./...

stest:
	go test -short ./...

vtest:
	go test -v ./...

rtest:
	go test -v -run $(test) ./...

build:
	go build -o ./player ./cmd/player.go