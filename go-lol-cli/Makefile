all: build check

build: 
	go build

check:
	go vet
	go test -coverprofile=cover.out -covermode=count


