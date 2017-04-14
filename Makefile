TEST?=$$(go list ./...)

default: test build

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=60s -parallel=4

build:
	go build -o terraform-ci main.go

.PHONY: test