protobuf:
	protoc --proto_path=proto/ --go_out=models/execution proto/execution.proto
	protoc --proto_path=proto/ --go_out=models/example proto/example.proto
	protoc --proto_path=proto/ --go_out=models/unsupported proto/unsupported.proto
