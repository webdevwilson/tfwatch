TEST?=$$(go list ./... | grep -v '/vendor/')

APP:=terraform-ci
VERSION:=0.0.1

# Docker variables
DOCKER_HOST:=webdevwilson
DOCKER_NAME:=$(DOCKER_HOST)/$(APP)

default: test build

tools:
	go get -u github.com/kardianos/govendor
	go get -u github.com/pilu/fresh
	
test:
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=60s -parallel=4

clean:
	rm routes/site_static.go
	rm -rf site/dist

site/dist:
	cd site && \
	npm install && \
	npm run build

terraform-ci: site/dist
	go build -o terraform-ci main.go

build: terraform-ci
	
install: build
	mkdir -p $(GOPATH)/bin
	cp terraform-ci $(GOPATH)/bin/

docker-build: site/dist
	docker build -t $(DOCKER_NAME) .
	docker tag $(DOCKER_NAME) $(DOCKER_NAME):$(VERSION)

.PHONY: test tools build install docker-build