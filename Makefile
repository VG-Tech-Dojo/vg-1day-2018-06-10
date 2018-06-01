nickname=
repository_name=$(shell basename $(PWD))

DOCKER_IMAGE        := $(repository_name)
DOCKER_WORKDIR_BASE := /go/src/github.com/VG-Tech-Dojo/vg-1day-2018

.PHONY: setup/* docker/*

setup/mac: $(nickname)
	$(MAKE) setup/bsd

setup/bsd: $(nickname) ## for mac
	sed -i '' -e 's/original/$(nickname)/g' ./$(nickname)/*.go
	sed -i '' -e 's/original/$(nickname)/g' ./$(nickname)/**/*.go
	sed -i '' -e 's/vg-1day-2018/$(repository_name)/g' ./$(nickname)/*.go
	sed -i '' -e 's/vg-1day-2018/$(repository_name)/g' ./$(nickname)/**/*.go

setup/gnu: $(nickname) ## for linux
	sed --in-place 's/original/$(nickname)/g' ./$(nickname)/*.go
	sed --in-place 's/original/$(nickname)/g' ./$(nickname)/**/*.go
	sed --in-place 's/vg-1day-2018/$(repository_name)/g' ./$(nickname)/*.go
	sed --in-place 's/vg-1day-2018/$(repository_name)/g' ./$(nickname)/**/*.go

$(nickname):
	cp -rf original $(nickname)

docker/build:
	docker build -t $(DOCKER_IMAGE) .

docker/deps:
	docker run --rm -v $(CURDIR):$(DOCKER_WORKDIR_BASE) -it --workdir $(DOCKER_WORKDIR_BASE)/original $(DOCKER_IMAGE) deps

docker/run:
	docker run --rm -v $(CURDIR):$(DOCKER_WORKDIR_BASE) -it --workdir $(DOCKER_WORKDIR_BASE)/original -p 8080:8080 $(DOCKER_IMAGE)

docker/test:
	docker run --rm -v $(CURDIR):$(DOCKER_WORKDIR_BASE) -it --workdir $(DOCKER_WORKDIR_BASE)/original $(DOCKER_IMAGE) test

docker/deps/%: $(@F)
	docker run --rm -v $(CURDIR):$(DOCKER_WORKDIR_BASE) -it --workdir $(DOCKER_WORKDIR_BASE)/$(@F) $(DOCKER_IMAGE) deps

docker/run/%: $(@F)
	docker run --rm -v $(CURDIR):$(DOCKER_WORKDIR_BASE) -it --workdir $(DOCKER_WORKDIR_BASE)/$(@F) -p 8080:8080 $(DOCKER_IMAGE)

docker/test/%: $(@F)
	docker run --rm -v $(CURDIR):$(DOCKER_WORKDIR_BASE) -it --workdir $(DOCKER_WORKDIR_BASE)/$(@F) $(DOCKER_IMAGE) test
