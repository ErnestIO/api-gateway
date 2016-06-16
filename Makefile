install:
	go install

lint:
	golint ./...
	go vet ./...

test:
	go test -v ./... --cover

deps: dev-deps
	go get golang.org/x/crypto/scrypt
	go get github.com/nats-io/nats
	go get github.com/labstack/echo
	go get github.com/dgrijalva/jwt-go

dev-deps:
	go get github.com/smartystreets/goconvey
	go get -u github.com/golang/lint/golint
