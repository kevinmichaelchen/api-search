.PHONY: all
all:
	$(MAKE) gen-proto

.PHONY: gen-proto
gen-proto:
	buf generate idl
