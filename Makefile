.PHONY: build-server
build-server:
	sh hack/build.sh ./cmd/pando-server

.PHONY: build-worker
build-worker:
	sh hack/build.sh ./cmd/pando-worker

TAG = $(shell git describe --tags --abbrev=0)

.PHONY: pando/worker
pando/worker:
	docker build -t pando/worker:${TAG} -t pando/worker:latest -f ./docker/Dockerfile.worker .

.PHONY: pando/server
pando/server:
	docker build -t pando/server:${TAG} -t pando/server:latest -f ./docker/Dockerfile.server .

clean:
	@rm -rf ./builds
