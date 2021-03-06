STAGE ?= dev
BRANCH ?= master
APPNAME ?= exitus
DOMAIN ?= exitus
WHITELIST_DOMAIN ?= wolfe.id.au

PACKAGE_BUCKET ?= versent-innovation-2019-lambda-$(AWS_REGION)

default: clean prepare test build archive package deploy
.PHONY: default

ci: clean test build-linux archive package deploy
.PHONY: ci

LDFLAGS := -ldflags="-s -w"

clean:
	@echo "--- clean all the things"
	@rm -rf dist
.PHONY: clean

prepare:
	@echo "--- prepare all the things"
	@go mod download
	@mkdir -p dist
.PHONY: prepare

test:
	@echo "--- test all the things"
	@go test -v -cover ./...
.PHONY: test

build:
	@docker run --rm \
		-v $$(pwd):/src/$$(basename $$(pwd)) \
		-v $$(go env GOPATH)/pkg/mod:/go/pkg/mod \
		-w /src/$$(basename $$(pwd)) -it golang:1.13 make build-linux
.PHONY: build

build-linux:
	@echo "--- build all the things"
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/triggers ./cmd/triggers
.PHONY: linux

archive:
	@echo "--- build an archive"	
	@cd dist && zip -X -9 -r ./handler.zip ./triggers
.PHONY: archive

package:
	@echo "--- package cognito stack to aws"
	@aws cloudformation package \
		--template-file sam/app/cognito.yml \
		--s3-bucket $(PACKAGE_BUCKET) \
		--output-template-file dist/packaged-template.yaml
.PHONY: package

packagetest:
	@echo "--- package test stack to aws"
	@aws cloudformation package \
		--template-file sam/testing/template.yaml \
		--s3-bucket $(PACKAGE_BUCKET) \
		--output-template-file dist/test-packaged-template.yaml
.PHONY: packagetest

deploytest:
	@echo "--- deploy cognito stack to aws"
	@aws cloudformation deploy \
		--template-file dist/test-packaged-template.yaml \
		--capabilities CAPABILITY_NAMED_IAM CAPABILITY_AUTO_EXPAND \
		--stack-name cognito-$(APPNAME)-$(STAGE)-$(BRANCH) \
		--parameter-overrides AppName=$(APPNAME) Stage=$(STAGE) Branch=$(BRANCH) Domain=$(DOMAIN) WhitelistDomain=$(WHITELIST_DOMAIN)
.PHONY: deploytest

deployci:
	@echo "--- deploy cognito stack to aws"
	@aws cloudformation deploy \
		--template-file sam/ci/template.yaml \
		--capabilities CAPABILITY_NAMED_IAM CAPABILITY_IAM CAPABILITY_AUTO_EXPAND \
		--stack-name serverless-cognito-auth-ci
.PHONY: deployci
