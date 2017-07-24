install:
	go install

lint:
	gometalinter --config .linter.conf

test:
	go test -v ./... --cover

deps:
	go get golang.org/x/crypto/scrypt
	go get github.com/nats-io/nats
	go get github.com/labstack/echo
	go get github.com/dgrijalva/jwt-go
	go get github.com/nu7hatch/gouuid
	go get github.com/ghodss/yaml
	go get github.com/ernestio/ernest-config-client
	go get github.com/blang/semver
	go get golang.org/x/crypto/pbkdf2
	go get github.com/ernestio/crypto
	go get github.com/ernestio/crypto/aes
	go get github.com/Sirupsen/logrus
	go get gopkg.in/r3labs/graph.v2

dev-deps: deps
	go get github.com/smartystreets/goconvey
	go get github.com/alecthomas/gometalinter
	gometalinter --install
