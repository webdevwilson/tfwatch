TEST?=$$(go list ./... | grep -v '/vendor/')
GOLANG_VERSION=1.8-alpine
WD=$(shell pwd)

# App information
APP:=tfwatch
VERSION:=$(shell cat VERSION)

# Docker variables
DOCKER_HOST:=webdevwilson
DOCKER_NAME:=$(DOCKER_HOST)/$(APP)
DOCKER_DIR:=/usr/local/go/src/github.com/webdevwilson/tfwatch
DOCKER_OPTS:=--rm -w $(DOCKER_DIR) \
		-v $(WD):$(DOCKER_DIR) \
		golang:$(GOLANG_VERSION)

default: test build

tools:
	go get -u github.com/kardianos/govendor
	go get -u github.com/pilu/fresh
	
test:
	@echo "[test] Running tests..."
	@echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=60s -parallel=4

clean:
	@echo "[clean] Removing build artifacts"
	rm -rf site/dist
	rm -f tfwatch

deps:
	@cd site && $(MAKE) deps

site/dist:
	@echo "[site] Building site"
	cd site && $(MAKE)

site: site/dist

tfwatch:
	@echo "[build] Building tfwatch"
	docker run $(DOCKER_OPTS) go build -o tfwatch main.go

build: tfwatch

install: build
	@echo "[install] Installing tfwatch to $(GOPATH)/bin"
	@mkdir -p $(GOPATH)/bin
	@cp tfwatch $(GOPATH)/bin/

serve:
	@echo "[serve]"
	docker run $(DOCKER_OPTS) -p 3000:3000 go run main.go $(DOCKER_DIR)/fixtures

docker: build site
	docker build -t $(DOCKER_NAME) .
	docker tag $(DOCKER_NAME) $(DOCKER_NAME):$(VERSION)

.PHONY: tools test clean site build install docker