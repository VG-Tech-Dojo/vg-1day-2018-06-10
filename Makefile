nickname=
repository_name=$(shell basename $(PWD))

DOCKER_IMAGE   := $(repository_name)
DOCKER_WORKDIR := /go/src/github.com/VG-Tech-Dojo/vg-1day-2018-06-10

.PHONY: setup/* docker/*

setup/mac: $(nickname)
	$(MAKE) setup/bsd

setup/bsd: $(nickname) ## for mac
	sed -i '' -e 's/original/$(nickname)/g' ./$(nickname)/*.go
	sed -i '' -e 's/original/$(nickname)/g' ./$(nickname)/**/*.go
	sed -i '' -e 's/vg-1day-2018-06-10/$(repository_name)/g' ./$(nickname)/*.go
	sed -i '' -e 's/vg-1day-2018-06-10/$(repository_name)/g' ./$(nickname)/**/*.go

setup/gnu: $(nickname) ## for linux
	sed --in-place 's/original/$(nickname)/g' ./$(nickname)/*.go
	sed --in-place 's/original/$(nickname)/g' ./$(nickname)/**/*.go
	sed --in-place 's/vg-1day-2018-06-10/$(repository_name)/g' ./$(nickname)/*.go
	sed --in-place 's/vg-1day-2018-06-10/$(repository_name)/g' ./$(nickname)/**/*.go

$(nickname):
	cp -rf original $(nickname)
	echo "vendor/*" >> $(nickname)/.gitignore

docker/build:
	docker build -t $(DOCKER_IMAGE) .

docker/deps: docker/deps/original

docker/run: docker/run/original

docker/test: docker/test/original

docker/deps/%: $(@F)
	docker run --rm -v $(CURDIR):$(DOCKER_WORKDIR) -it $(DOCKER_IMAGE) -C $(@F) deps

docker/run/%: $(@F)
	docker run --rm -v $(CURDIR):$(DOCKER_WORKDIR) -it -p 8080:8080 $(DOCKER_IMAGE) -C $(@F) run

docker/test/%: $(@F)
	docker run --rm -v $(CURDIR):$(DOCKER_WORKDIR) -it $(DOCKER_IMAGE) -C $(@F) test
