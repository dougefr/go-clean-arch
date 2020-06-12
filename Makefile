#!make

path=./...

setup: .make.setup
.make.setup:
	GO111MODULE=off go get -u github.com/kyoh86/richgo
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get -u github.com/golang/mock/mockgen
	GO111MODULE=off go get -u golang.org/x/tools/cmd/cover
	touch .make.setup

fmt:
	go fmt $(path)

lint:
	$(GOPATH)/bin/golint -set_exit_status -min_confidence 0.9 $(path)
	@echo "Golint found no problems on your code!"

test: mock
	$(GOPATH)/bin/richgo test $(path) $(args)

fullcover:
	go test -coverprofile=cover.out $(path)
	go tool cover -func=cover.out || true
	rm cover.out

run:
	go run ./cmd/user-api/main.go

mock:
	mockgen -source=./core/usecase/igateway/user.go -destination=./core/usecase/igateway/mock_igateway/user.go
	mockgen -source=./interface/iinfra/database.go -destination=./interface/iinfra/mock_iinfra/database.go
	mockgen -source=./interface/iinfra/logprovider.go -destination=./interface/iinfra/mock_iinfra/logprovider.go