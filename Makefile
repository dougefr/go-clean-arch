#!make

path=./...

setup:
	GO111MODULE=off go get -u github.com/kyoh86/richgo
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=on go get github.com/golang/mock/mockgen@latest

fmt:
	go fmt $(path)

lint:
	$(GOPATH)/bin/golint -set_exit_status -min_confidence 0.9 $(path)
	@echo "Golint found no problems on your code!"

run:
	go run ./cmd/user-api/main.go

mock:
	mockgen -source=./core/igateway/user.go -destination=./core/igateway/mock_igateway/user.go
	mockgen -source=./interface/iinfra/database.go -destination=./interface/iinfra/mock_iinfra/database.go
	mockgen -source=./interface/iinfra/logprovider.go -destination=./interface/iinfra/mock_iinfra/logprovider.go