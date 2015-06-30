default : all check 

all:
	go build

check:
	go vet
	go test -coverprofile=cover.out -covermode=count
#	golint


generate-test-data:
	make -C test-data-fetcher
	mkdir -p data
	./test-data-fetcher/test-data-fetcher -k $(API_KEY) -o data/go-lol_testdata.json
	go-bindata -pkg="lol" -o testdata_test.go data
