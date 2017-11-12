all: dep test build docker-image

dep:
	go get -t ./...

test:
	go test -v ./...

integration-test:
	go test -v -tags integration ./...

run: build
	@export `cat ${mkfile_path}.env | xargs`; ./simplesurance-api

build:
	go build -o simplesurance-api cmd/simplesurance-api/main.go

build-static:
	CGO_ENABLED=0 go build -v -a -installsuffix cgo -o simplesurance-api cmd/simplesurance-api/main.go

docker-build:
	docker build -t guilherme-santos/simplesurance:build . -f Dockerfile.build
	docker create --name simplesurance-builded guilherme-santos/simplesurance:build
	docker cp simplesurance-builded:/go/src/github.com/guilherme-santos/simplesurance/simplesurance-api .
	docker cp simplesurance-builded:/etc/ssl/certs/ca-certificates.crt .
	docker rm -f simplesurance-builded

docker-image: docker-build
	docker build -t guilherme-santos/simplesurance .
