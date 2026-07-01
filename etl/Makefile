# build-ExtractFunction:
# 	go mod tidy
# 	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap cmd/extract/main.go
# 	@echo "Build complete. Binary size:"
# 	ls -lh $(ARTIFACTS_DIR)/bootstrap

build-TransformFunction:
	go mod tidy
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap cmd/transform/main.go

build-LoadFunction:
	go mod tidy
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap cmd/load/main.go

build-FetchFunction:
	go mod tidy
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap cmd/fetch/main.go
	@echo "Build complete. Binary size:"
	ls -lh $(ARTIFACTS_DIR)/bootstrap


build-DownloadFunction:
	go mod tidy
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap cmd/download/main.go
	@echo "Build complete. Binary size:"
	ls -lh $(ARTIFACTS_DIR)/bootstrap