.PHONY: restore run-tests cover vet lint

COVERALLS_TOKEN := "DJljOZEYgGX0vkvoPCVNxyIupLJ6VscoO"
GITHUB_API_TOKEN := ""
VERSION :=""

all: restore run-tests cover vet lint
test: restore run-tests
cover-remote: restore run-cover-remote

restore:
	go get -u github.com/golang/lint/golint
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

run-cover-remote:
	go get -u github.com/mattn/goveralls
	go test -covermode=count -coverprofile=cover.tmp
	goveralls -service travis-ci -coverprofile cover.tmp

run-tests:
	go test -cover `go list ./... | grep -v /vendor/`

cover:
	go test -cover `go list ./... | grep -v /vendor/`

lint:
	golint `go list ./... | grep -v /vendor/`

vet:
	go vet `go list ./... | grep -v /vendor/`
