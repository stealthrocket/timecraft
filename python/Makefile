.PHONY: docker-build docker-buildx

TIMECRAFT_VERSION ?= "main-1687814584-4c9ed9ce"

ifeq ($(GITHUB_BRANCH_NAME),)
	branch := $(shell git rev-parse --abbrev-ref HEAD)-
else
	branch := $(GITHUB_BRANCH_NAME)-
endif
commit_timestamp := $(shell git show --no-patch --format=%ct)-
ifeq ($(GITHUB_SHA),)
	commit := $(shell git rev-parse --short=8 HEAD)
else
	commit := $(shell echo $(GITHUB_SHA) | cut -c1-8)
endif
container.version ?= $(if $(RELEASE_TAG),$(RELEASE_TAG),$(shell git describe --tags || echo '$(subst /,-,$(branch))$(commit_timestamp)$(commit)'))
container.image ?= ghcr.io/stealthrocket/timecraft-python

docker-build: 
	docker build -f Dockerfile.timecraft -t $(container.image):$(container.version) .

docker-buildx:
	docker buildx build --push \
		--platform=linux/amd64,linux/arm64 \
		--build-arg TIMECRAFT_VERSION=$(TIMECRAFT_VERSION) \
		-f Dockerfile.timecraft \
		-t $(container.image):$(container.version) .