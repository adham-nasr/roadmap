build-ExtractFunction:
	go mod tidy
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap cmd/extract/main.go
