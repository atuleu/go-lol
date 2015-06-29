default : all check 

all:
	go build

check:
	go vet
	go test -coverprofile=cover.out -covermode=count
#	golint
