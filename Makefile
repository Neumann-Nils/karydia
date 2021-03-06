# Copyright 2019 Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

IMAGE_REPO=karydia
IMAGE_NAME=karydia

KUBERNETES_SERVER ?= ""
KUBECONFIG_PATH ?= "$(HOME)/.kube/config"

VERSION=$(shell git describe --tags --always --dirty)

.PHONY: all
all: build

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/karydia \
		-ldflags "-s -X github.com/karydia/karydia.Version=$(VERSION)" \
		cli/main.go

.PHONY: container
container:
	docker build -t $(IMAGE_REPO)/$(IMAGE_NAME) .

.PHONY: codegen
codegen:
	hack/update-codegen.sh

.PHONY: test
test:
	go test $(shell go list ./... | grep -v /vendor/ | grep -v /tests/)

.PHONY: e2e-test
e2e-test:
	go test -v ./tests/e2e/... --server $(KUBERNETES_SERVER) --kubeconfig $(KUBECONFIG_PATH)
