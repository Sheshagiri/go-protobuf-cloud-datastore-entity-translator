deps:
	go mod vendor

protobuf:
	protoc --proto_path=proto/ --go_out=models/execution proto/execution.proto
	protoc --proto_path=proto/ --go_out=models/example proto/example.proto
	protoc --proto_path=proto/ --go_out=models/unsupported proto/unsupported.proto

unit-tests: deps
	@echo "Running Unit Tests"
	go test -v ./... -mod=vendor -race -coverprofile=coverage.txt -covermode=atomic

integration-tests: deps
	@echo "Running Integration Tests"
	go test -v ./... -mod=vendor -race -tags=integration
