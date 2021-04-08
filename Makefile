TEST?=$(shell go list ./...)
COUNT?=1

.PHONY: testacc

# testacc runs acceptance tests
testacc: # Run acceptance tests
	@echo "WARN: Acceptance tests will take a long time to run and may cost money. Ctrl-C if you want to cancel."
	PACKER_ACC=1 go test -count $(COUNT) -v $(TEST) $(TESTARGS) -timeout=10m