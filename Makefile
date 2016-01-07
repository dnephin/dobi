.PHONY: build binary shell test-unit


BUILD_ID ?= $(shell git rev-parse --short HEAD 2>/dev/null)
DOCKER_IMAGE := buildpipe-dev:$(BUILD_ID)

VOLUMES := \
	-v $(CURDIR):/go/src/github.com/dnephin/buildpipe \
	-v $(CURDIR)/dist/bin:/go/bin \
	-v $(CURDIR)/dist/pkg:/go/pkg

all: binary

build:
	docker build -t $(DOCKER_IMAGE) -f dockerfiles/Dockerfile.build .

dist:
	mkdir dist/

binary: dist build
	docker run --rm -e VERSION=$(BUILD_ID) $(VOLUMES) $(DOCKER_IMAGE)

shell: dist build
	docker run --rm -ti $(VOLUMES) $(DOCKER_IMAGE) bash

test-unit: build
	docker run --rm -ti $(VOLUMES) $(DOCKER_IMAGE) go test -v \$(glide novendor)

test: test-unit
