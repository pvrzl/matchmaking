# Load environment variables
BUILDTAG=$(shell git tag --points-at HEAD)
VERSION=$(shell git rev-parse HEAD)
BUILDDATE=$(shell date '+%F-%T%z')
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
ifeq ($(VERSION),)
	VERSION="unknown"
endif
ifeq ($(BUILDTAG),)
	BUILDTAG="unknown"
endif

.PHONY: echo-env
echo-env:
	echo ${BUILDTAG}
	echo ${VERSION}
	echo ${BUILDDATE}
	echo ${COMMIT}
	echo ${BRANCH}

.PHONY: build-http generator-repository generator-usecase generator-interface migrate-create migrate-up migrate-down migrate-fix
build-http:
	go build -ldflags="-X 'main.BuildTag=${BUILDTAG}' -X 'main.Version=${VERSION}' -X 'main.BuildDate=${BUILDDATE}' -X 'main.Commit=${COMMIT}' -X 'main.Branch=${BRANCH}'" -v -o app-http cmd/server/http/*.go

generator-repository:
	go run cmd/generator/repository/main.go

generator-interface:
	 go run cmd/generator/interface/main.go

generator-usecase:
	 go run cmd/generator/usecase/main.go

migrate-create:
	 go run cmd/migrate/* -command=create

migrate-up:
	 go run cmd/migrate/* -command=up

migrate-down:
	 go run cmd/migrate/* -command=down

migrate-fix:
	 go run cmd/migrate/* -command=fix
