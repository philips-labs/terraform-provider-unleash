default: testacc

# Run acceptance tests
.PHONY: testacc
testacc: export TF_ACC=1
testacc:
	go test ./... -v $(TESTARGS) -timeout 120m
