protobuf:
	protoc -I=$(GOPATH)/src:. --go_out=models/execution/ proto/execution.proto
	protoc -I=$(GOPATH)/src:. --go_out=models/example/ proto/example.proto
	protoc -I=$(GOPATH)/src:. --go_out=models/unsupported/ proto/unsupported.proto