include contrib/build/makefiles/pkg/base/base.mk
include contrib/build/makefiles/pkg/string/string.mk
include contrib/build/makefiles/pkg/color/color.mk
include contrib/build/makefiles/pkg/functions/functions.mk
include contrib/build/makefiles/target/buildenv/buildenv.mk
include contrib/build/makefiles/target/go/go.mk
THIS_FILE := $(firstword $(MAKEFILE_LIST))
SELF_DIR := $(dir $(THIS_FILE))
.PHONY: test build clean run temp-clean
.SILENT: test build clean run temp-clean
PORT:=8080
RPC_ENDPOINT:=rpc
build: 
	- $(call print_running_target)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) go-build
	- $(call print_completed_target)
run: kill
	- $(call print_running_target)
	- bin$(PSEP)dare daemon --api-addr=127.0.0.1:${PORT} > $(PWD)/server.log 2>&1 &
	- $(call print_completed_target)
clean:
	- $(call print_running_target)
	- @$(MAKE) --no-print-directory -f $(THIS_FILE) go-clean
	- $(call print_completed_target)
kill : temp-clean
	- $(call print_running_target)
	- $(RM) $(PWD)/server.log
	- for pid in $(shell ps  | grep "dare" | awk '{print $$1}'); do kill -9 "$$pid"; done
	- $(call print_completed_target)
temp-clean:
	- $(call print_running_target)
	- $(RM) /tmp/go-build*
	- $(call print_completed_target)
