TEST?=$$(go list ./... | grep -v '/vendor/')

APP:=tfwatch
VERSION:=0.0.2

# Docker variables
DOCKER_HOST:=webdevwilson
DOCKER_NAME:=$(DOCKER_HOST)/$(APP)

docker_builder:
	@echo [docker_builder] Building docker
	docker build -file docker/Dockerfile-build -tag tfwatch-builder .

default: test build

tools:
	go get -u github.com/kardianos/govendor
	go get -u github.com/pilu/fresh
	
test:
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=60s -parallel=4

clean:
	rm -rf site/dist

site/dist:
	cd site && \
	npm install && \
	npm run build

tfwatch: site/dist
	go build -o tfwatch main.go

build: tfwatch
	
install: build
	mkdir -p $(GOPATH)/bin
	cp tfwatch $(GOPATH)/bin/

docker-build: site/dist
	docker build -t $(DOCKER_NAME) .
	docker tag $(DOCKER_NAME) $(DOCKER_NAME):$(VERSION)

.PHONY: test tools build install docker-build