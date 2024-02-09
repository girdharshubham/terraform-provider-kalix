default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: local
local:
	@sed 's#GOBIN#${GOPATH}/bin#g' sample-terraformrc > ~/.terraformrc
