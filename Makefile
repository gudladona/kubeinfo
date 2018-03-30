APP := kubeinfo
ROOT_PATH := $(PWD)
BUILD_PATH :=  $(PWD)/build
COV_PATH :=  $(PWD)/coverage
SHELL := /bin/bash
TAG_NAME := latest

.PHONY: build-setup docker

build: clean build-setup docker

push: push-to-docker

deploy: build push deploy-kube

destroy: delete-kube

#Clean the build directory
clean:
	@echo "===> Running Cleanup"
	rm -rf $(BUILD_PATH)/

#Setup the build directory
build-setup:
	@echo "--- Cleaned setting up $(BUILD_PATH)"
	mkdir -p $(BUILD_PATH)/
	cp .kube/* $(BUILD_PATH)/

#Test Go code and generate coverage reports
test:
	rm -rf $(COV_PATH)/
	mkdir -p $(COV_PATH)/
	go get ./...
	gocov test ./... -v > $(COV_PATH)/coverage.json && \
	gocov-html $(COV_PATH)/coverage.json > $(COV_PATH)/coverage.html

#Build Docker image into the build directory
docker:
	docker build -f Dockerfile -t $(APP):$(TAG_NAME) .

#Delete Kube deployments
delete-kube:
	kubectl delete -f build/deploy.yml
	kubectl delete -f build/service.yml

#Apply Kube deployments
deploy-kube:
	kubectl create configmap kubeinfo-config --dry-run -o yaml --from-file=build/config.cfg | kubectl replace -f -
	kubectl apply -f build/deploy.yml
	kubectl apply -f build/service.yml

#Push to Public Docker repo (https://hub.docker.com/r/gudladona87/kubeinfo/)
push-to-docker:
	echo "Pushing docker image for $(APP)"
	docker tag $(APP):$(TAG_NAME) $(DOCKER_USER)/$(APP):$(TAG_NAME)
	docker push $(DOCKER_USER)/$(APP):$(TAG_NAME)
