VENDOR_DIR := ${PWD}/_vendor

ifdef GOPATH
	GOPATH := ${VENDOR_DIR}:${GOPATH}
else
	GOPATH := ${VENDOR_DIR}
endif

PROJECT := github.com/jeremija/gol
DEPS := \
	github.com/hpcloud/tail \
	github.com/BurntSushi/toml \
	github.com/influxdata/influxdb/client/v2

export GOPATH

.PHONY: default
default: build

.PHONY: build
build: vet

	go build -v -o ./bin/gol ./app

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

	@for dep in ${DEPS}; do \
		echo getting dependency $$dep; \
		go get -u $$dep; \
	done

	mkdir -p ${VENDOR_DIR}/${PROJECT}
	rmdir ${VENDOR_DIR}/${PROJECT}
	ln -s ../../../.. ${VENDOR_DIR}/src/${PROJECT}

.PHONY: vendor-update
vendor-update: vendor-get

	rm -rf `find ${VENDOR_DIR}/src -type d -name .git`
	rm -rf `find ${VENDOR_DIR}/src -type d -name .hg`
	rm -rf `find ${VENDOR_DIR}/src -type d -name .bzr`
	rm -rf `find ${VENDOR_DIR}/src -type d -name .svn`

vendor-install:

	export GOPATH=${VENDOR_DIR}

	@for dep in ${DEPS}; do \
		echo installing dependency $$dep; \
		go install $$dep; \
	done

install:

	export GOPATH=${VENDOR_DIR}
	go install ${PROJECT}

.PHONY: vet
vet:

	go vet ./ ./dispatchers
