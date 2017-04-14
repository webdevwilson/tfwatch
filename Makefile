TEST?=$$(go list ./...)

default: test

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=60s -parallel=4

.PHONY: test