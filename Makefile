STAGE ?= dev
BRANCH ?= master
APPNAME ?= exitus
DOMAIN ?= exitus
WHITELIST_DOMAIN ?= wolfe.id.au
S3BUCKET ?= versent-innovation-2019-lambda-${AWS_REGION}

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
	@aws s3 mb s3://$(S3BUCKET) || true
	@aws cloudformation package \
		--template-file sam/app/cognito.yml \
		--s3-bucket $(S3BUCKET) \
		--output-template-file dist/cognito.out.yml
.PHONY: package

deploy:
	@echo "--- deploy cognito stack to aws"
	@aws cloudformation deploy \
		--template-file dist/cognito.out.yml \
		--capabilities CAPABILITY_NAMED_IAM \
		--stack-name cognito-$(APPNAME)-$(STAGE)-$(BRANCH) \
		--parameter-overrides AppName=$(APPNAME) Stage=$(STAGE) Branch=$(BRANCH) Domain=$(DOMAIN) WhitelistDomain=$(WHITELIST_DOMAIN)
.PHONY: deploy

deployci:
	@echo "--- deploy cognito stack to aws"
	@aws cloudformation deploy \
		--template-file sam/ci/template.yaml \
		--capabilities CAPABILITY_NAMED_IAM \
		--stack-name serverless-cognito-auth-ci
.PHONY: deployci
