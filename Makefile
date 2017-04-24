TEST?=$$(go list ./...)

default: test build

tools:
	go get -u github.com/kardianos/govendor
	go get -u github.com/mjibson/esc
	
test:
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=60s -parallel=4

clean:
	rm routes/site_static.go

site/dist:
	cd site; npm run build

routes/site_static.go: site/dist
	go generate routes/site.go

build: routes/site_static.go
	go build -o terraform-ci main.go

run: routes/site_static.go
	go run main.go

.PHONY: test tools build