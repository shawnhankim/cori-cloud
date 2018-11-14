#
# Makefile to build cori-cloud which is sample code to handle multiple cloud provider and Kubernetes
#
# PACKAGE refers the main directory of the project and MAIN is the entry point.
#

# PACKAGE refers the main directory of the project and MAIN is the entry point.
PACKAGE   = cori-cloud
MAIN 	  = main.go
BIN       = bin
OUT_DIR  ?= $(BIN)
GO_GCFLAGS="-newexport=0"

# tools and shortcuts
GO      = go
GODOC   = godoc
GOFMT   = gofmt
GO_PKGS = $(shell $(GO) list ./... | grep -v /vendor/ )
DEP   	= dep
TIMEOUT = 15
V 		= 1
Q 		= $(if $(filter 1,$V),,@)
M 		= $(shell printf "\033[34;1m▶\033[0m")
CD		= $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

define update_dependencies
	@echo $(GOPATH)
	$Q cd $(BASE) && pwd && GOPATH=$(GOPATH) $(DEP) ensure -update
endef

export GOPATH 

.PHONY: all fmt clean version build build-container package semver upload

init:
	go get -d -u github.com/golang/dep
	go get -u github.com/hairyhenderson/gomplate
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install
	go get -u github.com/tebeka/go2xunit
	go get -u github.com/axw/gocov/...
	go get -u github.com/AlekSi/gocov-xml
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/rs/xid
	go get -u github.com/aws/aws-sdk-go/...
	go get -u github.com/aws/aws-sdk-go
	go get -u github.com/golang/protobuf/proto
	go get -u github.com/golang/protobuf/protoc-gen-go


# for now just build until pipeline adds hooks for test/coverage
all: test lint-crit build templates package build-container; $(info $(M) building all...) @

deps:
	dep ensure -v

build:; $(info $(M) building executable image of cori-cloud...) @
	$Q CGO_ENABLED=0 go build -a -installsuffix cgo \
		-tags 'release netgo' \
		-o $(OUT_DIR)/$(PACKAGE) $(MAIN)

build-linux:; $(info $(M) building executable image of cori-cloud...) @
	$Q CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
		-tags 'release netgo' \
		-o $(OUT_DIR)/$(PACKAGE) $(MAIN)

templates: ; $(info $(M) processing templates..) @
	rm -rf deploy
	cp -R deploy-templates deploy
	gomplate --input-dir=deploy-templates --output-dir=deploy

package: ; $(info $(M) creating archives..) @
	mkdir -p $(OUT_DIR) && tar czf $(ARTIFACT) -C deploy .

upload:
	bot artifactory-uploader $(ARTIFACT)

test: fmt; $(info $(M) running tests...) @
	$Q $(GO) vet $(GO_PKGS)
	echo "mode: set" > coverage-all.out
	$(foreach pkg,$(GO_PKGS),\
		$(GO) test -v -race -coverprofile=coverage.out $(pkg) | tee -a test-results.out || exit 1;\
		tail -n +2 coverage.out >> coverage-all.out || exit 1;)
	$(GO) tool cover -func=coverage-all.out


fmt: ; $(info $(M) running gofmt...) @ ## Run gofmt on all source files
	@ret=0 && for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GOFMT) -l -s -w $$d/*.go || ret=$$? ; \
	 done ; exit $$ret

# Dependency management  - requires local tools to be available
vendor: ; $(info $(M) retrieving dependencies...)
	$Q $(DEP) ensure

# Misc

clean: ; $(info $(M) cleaning...)	@ ## Cleanup everything
	rm -rf bin ./build_version coverage* ./deploy test-results.out unit_tests.xml full-lint-report.xml

version:
	@echo $(VERSION)

# Running non-critical checks
.PHONY: lint
lint:
	gometalinter --skip=vendor --disable-all --enable=gocyclo --enable=gas ./...

.PHONY: lint-crit
lint-crit:
	# Running critical checks
	gometalinter \
	--enable-gc \
	--skip=vendor \
	--skip=bindings \
	--disable-all \
	--enable=golint \
	--enable=misspell \
	--enable=vetshadow \
	--enable=gotype \
	--enable=vet \
	--enable=goconst \
	--enable=ineffassign \
	--enable=staticcheck \
	--deadline=500s \
	./...;

