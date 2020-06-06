#!make

path=./...

setup:
	GO111MODULE=off go get -u github.com/kyoh86/richgo
	GO111MODULE=off go get -u golang.org/x/lint/golint

fmt:
	go fmt $(path)

lint:
	$(GOPATH)/bin/golint -set_exit_status -min_confidence 0.9 $(path)
	@echo "Golint found no problems on your code!"

run:
	go run ./cmd/user-api/main.go