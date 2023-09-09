

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin

UNAME := $(shell uname)
ARCH := $(shell uname -m)

## goreleaser tool binary
GORELEASER="$(LOCALBIN)/goreleaser"
GORELEASER_VERSION="1.20.0"

# Setup release variables
ifeq ($(UNAME), Linux)
GORELEASER_RELEASE="Linux_x86_64"
endif
ifeq ($(UNAME), Darwin)
ifeq ($(ARCH), amd64)
GORELEASER_RELEASE="Darwin_x86_64"
else
GORELEASER_RELEASE="Darwin_${ARCH}"
endif
endif

######################################################################
##@ General
######################################################################

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo ""
	
	@echo -e "\033[0;33mVariable values:\033[0m"
ifndef TAG
	@echo -e "  \033[0;1mTAG\033[0m: <NOT_SET> (e.g. example TAG=1.0.0)"
else
	@echo -e "  \033[0;1mTAG\033[0m: $(TAG)"
endif
ifndef COMMENT
	@echo -e "  \033[0;1mCOMMENT\033[0m: <NOT_SET> (e.g. example COMMENT=\"Fix for some problem\")"
else
	@echo -e "  \033[0;1mCOMMENT\033[0m: $(COMMENT)"
endif
ifndef GITHUB_TOKEN
	@echo -e "  \033[0;1mGITHUB_TOKEN\033[0m: <NOT_SET>"
else
	@echo -e "  \033[0;1mGITHUB_TOKEN\033[0m: $(GITHUB_TOKEN)"
endif
	@echo ""
	@echo -e "\033[0;33mNotes about variables:\033[0m"
	@echo -e "- In target \"\033[0;36mtag\033[0m\", variables \033[0;1mTAG\033[0m and \033[0;1mCOMMENT\033[0m are required."
	@echo -e "- In target \"\033[0;36mrelease\033[0m\", variable \033[0;1mGITHUB_TOKEN\033[0m is required."
	@echo ""

######################################################################
##@ Development
######################################################################

snapshot: ## Build a snapshot release
	@$(GORELEASER) release --snapshot --clean

######################################################################
##@ Release new version
######################################################################

tag: ## Create a tag
ifeq ($(TAG),)
	$(error variable TAG must be specified)
endif
ifeq ($(COMMENT),)
	$(error variable COMMENT must be specified)
endif
	git tag -a v$(TAG) -m "$(COMMENT)"
	git push origin v$(TAG)

release: ## Use goreleaser to release a new version
ifeq ($(GITHUB_TOKEN),)
	$(error variable GITHUB_TOKEN must be specified)
endif
	@$(GORELEASER) release --clean

######################################################################
##@ Install Tools
######################################################################

$(LOCALBIN):
	mkdir -p $(LOCALBIN)

goreleaser: $(GORELEASER) ## Download kustomize locally if necessary.
$(GORELEASER): $(LOCALBIN)
	$(call install-goreleaser,$(GORELEASER),$(GORELEASER_VERSION),${GORELEASER_RELEASE})

install-tools: goreleaser ## Install all needed tools

######################################################################
## Defines
######################################################################

# install-goreleaser: takes 3 params:
# $1: target install dir (e.g. $PROJECT_DIR/bin/goreleaser)
# $2: the goreleaser version to download, e.g. 1.13.0
# $3: the goreleaser release (based on OS), e.g. "Darwin_x86_64" for mac or "Linux_x86_64" for linux
#
# It will download a goreleaser release and version and install it in the tools folder
#https://github.com/goreleaser/goreleaser/releases/download/v$(GORELEASER_VERSION)/goreleaser_$(GORELEASER_RELEASE).tar.gz
define install-goreleaser
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
echo "Downloading https://github.com/goreleaser/goreleaser/releases/download/v$(2)/goreleaser_$(3).tar.gz" ;\
echo "Downloading https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv$(2)/kustomize_v$(2)_$(3).tar.gz" ;\
echo "Goreleaser version: $(2), release: $(3)" ;\
wget -q https://github.com/goreleaser/goreleaser/releases/download/v$(2)/goreleaser_$(3).tar.gz -O goreleaser.tar.gz ;\
tar xzvf goreleaser.tar.gz ;\
cp ./goreleaser $(1) ;\
rm -rf $$TMP_DIR ;\
}
endef
