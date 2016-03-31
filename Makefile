GO_CMD ?= go

default: build

clean:
	rm mibig-api

build:
	$(GO_CMD) build

update-deps:
	$(GO_CMD) get -u

deps:
	$(GO_CMD) get

migrate:
	./mibig-api migratedb

test:
	$(GO_CMD) test
