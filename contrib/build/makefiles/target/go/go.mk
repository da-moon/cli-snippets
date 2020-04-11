ifeq ($(OS),Windows_NT)
    GO_PATH = $(subst \,/,${GOPATH})
else
    GO_PATH = ${GOPATH}
endif
fatal=fatal: No names found, cannot describe anything.
THIS_FILE := $(lastword $(MAKEFILE_LIST))
SELF_DIR := $(dir $(THIS_FILE))
GO_TARGET = $(notdir $(patsubst %/,%,$(dir $(wildcard ./cmd/*/.))))
GO_BUILD_WINDOWS_TARGETS = $(GO_TARGET:%=go-build-windows-%)
GO_BUILD_DARWIN_TARGETS = $(GO_TARGET:%=go-build-darwin-%)
GO_BUILD_LINUX_TARGETS = $(GO_TARGET:%=go-build-linux-%)
GO_BUILD_TARGETS = $(GO_TARGET:%=go-build-%)
GO_BUILD_OS_TARGETS = $(GO_BUILD_WINDOWS_TARGETS) $(GO_BUILD_DARWIN_TARGETS) $(GO_BUILD_LINUX_TARGETS)
.PHONY: $(GO_BUILD_TARGETS) $(GO_BUILD_OS_TARGETS)
.SILENT: $(GO_BUILD_TARGETS) $(GO_BUILD_OS_TARGETS) 
CGO=0
GO_ARCHITECTURE=amd64
GO_IMAGE=golang:buster
MOD=on
GO_PKG=github.com/da-moon/cli-snippets


.PHONY: go-build full-build build-darwin build-linux build-windows go-clean go-dependancy go-targets 
.SILENT: go-build full-build build-darwin build-linux build-windows go-clean go-dependancy go-targets
go-targets:
	- $(info  $(GO_BUILD_TARGETS) $(GO_BUILD_OS_TARGETS))

go-build: go-clean go-dependancy
	- $(call print_running_target)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) $(GO_BUILD_TARGETS)
	- $(call print_completed_target)
full-build: go-clean go-dependancy
	- $(CLEAR)
	- $(call print_running_target)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) build-linux
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) build-windows
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) build-darwin
	- $(call print_completed_target)

build-linux:  
	- $(call print_running_target)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) $(GO_BUILD_LINUX_TARGETS)
	- $(call print_completed_target)

build-windows:
	- $(call print_running_target)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) $(GO_BUILD_WINDOWS_TARGETS)
	- $(call print_completed_target)

build-mac-os:
	- $(call print_running_target)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) $(GO_BUILD_DARWIN_TARGETS)
	- $(call print_completed_target)

go-dependancy:
	- $(call print_running_target)
    ifeq ($(DOCKER_ENV),true)
ifeq (${MOD},on)
ifeq ($(wildcard ./go.mod),)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell docker_image="${GO_IMAGE}" container_name="go_builder_container" mount_point="/go/src/${GO_PKG}" cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go mod init"
endif
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell docker_image="${GO_IMAGE}" container_name="go_builder_container" mount_point="/go/src/${GO_PKG}" cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go mod tidy"
endif
ifeq (${MOD},off)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell docker_image="${GO_IMAGE}" container_name="go_builder_container" mount_point="/go/src/${GO_PKG}" cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go get -v -d ./..."
endif
    endif
    ifeq ($(DOCKER_ENV),false)
ifeq (${MOD},on)
ifeq ($(wildcard ./go.mod),)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go mod init"
endif
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go mod tidy"
endif
ifeq (${MOD},off)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell cmd="GO111MODULE=${MOD} \
    CGO_ENABLED=${CGO} \
    GOARCH=${GO_ARCHITECTURE} \
    go get -v -d ./..."
endif
    endif
	- $(call print_completed_target)



$(GO_BUILD_TARGETS): 
	- $(call print_running_target)
	- $(eval name=$(@:go-build-%=%))
	- $(eval command= GO111MODULE=${MOD})
	- $(eval command= ${command} CGO_ENABLED=${CGO})
	- $(eval command= ${command} GOARCH=${GO_ARCHITECTURE})
	- $(eval command= ${command} go build -a -installsuffix cgo)
	- $(eval command= ${command} -ldflags '-X ${GO_PKG}/version.Branch=${BRANCH}' )
	- $(eval command= ${command} -ldflags '-X ${GO_PKG}/version.BuildUser=${BUILDUSER}' )
	- $(eval command= ${command} -ldflags '-X ${GO_PKG}/version.BuildDate=${BUILDTIME}' )
	- $(eval command= ${command} -ldflags '-X ${GO_PKG}/version.Revision=${REVISION}' )
ifneq (${VERSION}, )
	- $(eval command= ${command} -ldflags '-X ${GO_PKG}/version.Version=${VERSION}' )
endif
	- for pid in $(shell ps  | grep "${name}" | awk '{print $$1}'); do kill -9 "$$pid"; done
	- $(eval command= ${command} -o .$(PSEP)bin$(PSEP)${name} .$(PSEP)cmd$(PSEP)${name} )
    ifeq ($(DOCKER_ENV),true)
	- @$(MAKE) --no-print-directory \
	 -f $(THIS_FILE) shell \
	 docker_image="${GO_IMAGE}" \
	 container_name="go_builder_container" \
	 mount_point="/go/src/${GO_PKG}" \
	 cmd="${command}"
    endif
    ifeq ($(DOCKER_ENV),false)
	- @$(MAKE) --no-print-directory \
	 -f $(THIS_FILE) shell cmd="${command}"
    endif
	- $(call print_completed_target)

.PHONY: $(GO_BUILD_OS_TARGETS)
.SILENT: $(GO_BUILD_OS_TARGETS)
$(GO_BUILD_OS_TARGETS): 
	- $(call print_running_target)
	- $(eval trimmed=$(@:go-build-%=%))
	- $(eval GOOS := $(firstword $(subst -, ,$(trimmed))))
	- $(eval name=$(@:go-build-$(GOOS)-%=%))
	- $(eval command= GO111MODULE=${MOD})
	- $(eval command= ${command} CGO_ENABLED=${CGO})
	- $(eval command= ${command} GOARCH=${GO_ARCHITECTURE})
	- $(eval command= ${command} GOARCH=${GO_ARCHITECTURE})
	- $(eval command= ${command} GOOS=${GOOS})
	- $(eval command= ${command} go build -a -installsuffix cgo \
			-o $(PWD)$(PSEP)bin$(PSEP)$(GOOS)$(PSEP)${name} $(PWD)$(PSEP)cmd$(PSEP)${name} \
		)
	- for pid in $(shell ps  | grep "${name}" | awk '{print $$1}'); do kill -9 "$$pid"; done
    ifeq ($(DOCKER_ENV),true)
	- @$(MAKE) --no-print-directory \
	 -f $(THIS_FILE) shell \
	 docker_image="${GO_IMAGE}" \
	 container_name="go_builder_container" \
	 mount_point="/go/src/${GO_PKG}" \
	 cmd="${command}"
    endif
    ifeq ($(DOCKER_ENV),false)
	- @$(MAKE) --no-print-directory \
	 -f $(THIS_FILE) shell cmd="${command}"
    endif
	- $(call print_completed_target)
	


go-clean:
	- $(CLEAR)
	- $(call print_running_target)
    ifeq ($(DOCKER_ENV),true)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) shell docker_image="${GO_IMAGE}" container_name="go_builder_container" mount_point="/go/src/${GO_PKG}" cmd="rm -rf /go/src/${GO_PKG}/bin/"
    else
	- $(RM) ./bin/
    endif
	- $(call print_completed_target)
