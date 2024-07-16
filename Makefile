lint:
	golangci-lint run

vet:
	go vet ./...

# test:
# 	sudo go mod vendor
# 	sudo go test -v ./...