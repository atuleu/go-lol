default: all check

all:
	go fmt
	go build


check:
	go vet
	go test -coverprofile=cover.out
	golint
