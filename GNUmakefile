default: testacc

# Run acceptance tests
.PHONY: testacc
testacc: export UNLEASH_API_URL=http://localhost:4242/api/
testacc: export UNLEASH_AUTH_TOKEN=token
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
