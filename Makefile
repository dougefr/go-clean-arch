#!make

path=./...
GOPATH=$(shell go env GOPATH)
min_coverage=1

setup: .make.setup
.make.setup:
	GO111MODULE=off go get -u github.com/kyoh86/richgo
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get -u github.com/golang/mock/mockgen
	GO111MODULE=off go get -u golang.org/x/tools/cmd/cover
	GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports
	touch .make.setup

deps:
	go get -v -t -d $(path)

fmt: setup
	go fmt $(path)
	find . -name \*.go -exec goimports -w {} \;

lint: setup mock
	$(GOPATH)/bin/golint -set_exit_status -min_confidence 0.9 $(path)
	@echo "Golint found no problems on your code!"
	go vet $(path)

test: mock
	$(GOPATH)/bin/richgo test $(path) $(args)

fullcover:
	go test -coverprofile=cover.out $(path)
	go tool cover -func=cover.out || true

run:
	go run ./cmd/user-api/main.go

mock:
	mockgen -source=./core/usecase/igateway/user.go -destination=./core/usecase/igateway/mock_igateway/user.go
	mockgen -source=./core/usecase/interactor/createuser.go -destination=./core/usecase/interactor/mock_interactor/createuser.go
	mockgen -source=./interface/iinfra/database.go -destination=./interface/iinfra/mock_iinfra/database.go
	mockgen -source=./interface/iinfra/logprovider.go -destination=./interface/iinfra/mock_iinfra/logprovider.go