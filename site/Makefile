WD=$(shell pwd)
NODE_VERSION=8.1.1

deps:
	docker run -v $(WD):/usr/src/app -w /usr/src/app node:$(NODE_VERSION) npm install

dist: deps
	docker run -v $(WD):/usr/src/app -w /usr/src/app node:$(NODE_VERSION) npm run build

default: dist

.PHONY: deps default