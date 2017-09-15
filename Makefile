GITHUB_API_TOKEN := ""
VERSION :=""

run-tests:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure
	go test -cover -v

cover:
	go test -coverprofile=cover.tmp && go tool cover -html=cover.tmp
