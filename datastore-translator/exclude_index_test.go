package datastore_translator

import (
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/execution"
	"github.com/golang/protobuf/ptypes"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"gotest.tools/assert"

	"testing"
)

func TestExclude(t *testing.T) {
	er := &execution.ExecutionRequest{
		Uuid:      "some uuid",
		Action:    "some random action",
		StartedOn: ptypes.TimestampNow(),
		Parameters: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"struct-key-string": {Kind: &structpb.Value_StringValue{StringValue: "some random string in proto.Struct"}},
				"struct-key-bool":   {Kind: &structpb.Value_BoolValue{BoolValue: true}},
				"struct-key-number": {Kind: &structpb.Value_NumberValue{NumberValue: 123456.12}},
				"struct-key-null":   {Kind: &structpb.Value_NullValue{}},
				"struct-key-list": {Kind: &structpb.Value_ListValue{
					ListValue: &structpb.ListValue{
						Values: []*structpb.Value{
							{Kind: &structpb.Value_NumberValue{NumberValue: 10}},
							{Kind: &structpb.Value_StringValue{StringValue: "hello, world"}},
							{Kind: &structpb.Value_BoolValue{BoolValue: true}},
							{Kind: &structpb.Value_NumberValue{NumberValue: 200}},
						},
					},
				},
				},
			},
		},
		Result: []byte("some large message, this will not be indexed in cloud datastore"),
	}
	entity, err := ProtoMessageToDatastoreEntity(er, true)
	assert.NilError(t, err)
	assert.Equal(t, entity.Properties["result"].ExcludeFromIndexes, true)
	assert.Equal(t, entity.Properties["parameters"].ExcludeFromIndexes, true)
}
