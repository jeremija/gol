VENDOR_DIR := ${PWD}/_vendor

ifdef GOPATH
	GOPATH := ${VENDOR_DIR}:${GOPATH}
else
	GOPATH := ${VENDOR_DIR}
endif

export GOPATH

.PHONY: default
default: build

.PHONY: build
build: vet

	go build -v -o ./bin/gol ./

.PHONY: doc
doc:

	godoc -http=:6060 -index

.PHONY: fmt
fmt:

	go fmt ./

.PHONY: lint
lint:

	golint ./src

.PHONY: run
run: build

	./bin/gol

.PHONY: test
test:

	go test ./ ./dispatchers

.PHONY: test-verbose
test-verbose:

	go test -v ./

.PHONY: vendor-clean
vendor-clean:

	rm -dRf ./_vendor/src

.PHONY: vendor-get
vendor-get: vendor-clean

	export GOPATH=${VENDOR_DIR}
	go get github.com/hpcloud/tail
	go get github.com/BurntSushi/toml
	go get github.com/influxdata/influxdb

.PHONY: vendor-update
vendor-update: vendor-get

	rm -rf `find ${VENDOR_DIR}/src -type d -name .git`
	rm -rf `find ${VENDOR_DIR}/src -type d -name .hg`
	rm -rf `find ${VENDOR_DIR}/src -type d -name .bzr`
	rm -rf `find ${VENDOR_DIR}/src -type d -name .svn`

.PHONY: vet
vet:

	go vet ./src/...
