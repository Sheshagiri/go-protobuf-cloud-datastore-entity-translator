package translator

import (
	"fmt"
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/unsupported"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/struct"
	"gotest.tools/assert"
	"testing"
)

func TestAddTSSupport(t *testing.T) {
	src := &unsupported.TS{
		StartedOn: ptypes.TimestampNow(),
	}
	fmt.Println("Source: ", src)
	srcEntity, err := ProtoMessageToDatastoreEntity(src, false)

	assert.NilError(t, err)
	fmt.Println("Source Datastore Entity: ", srcEntity)

	dst := &unsupported.TS{}

	err = DatastoreEntityToProtoMessage(&srcEntity, dst, false)
	assert.NilError(t, err)
	assert.DeepEqual(t, src, dst)
}

func TestAddStructSupport(t *testing.T) {
	src := &unsupported.StructMessage{
		StructKey: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"struct-key-string": {Kind: &structpb.Value_StringValue{"some random string in proto.Struct"}},
				"struct-key-bool":   {Kind: &structpb.Value_BoolValue{true}},
				"struct-key-number": {Kind: &structpb.Value_NumberValue{float64(123456.12)}},
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
	}
	fmt.Println("Source: ", src)
	srcEntity, err := ProtoMessageToDatastoreEntity(src, false)

	assert.NilError(t, err)
	fmt.Println("Source Datastore Entity: ", srcEntity)

	dst := &unsupported.StructMessage{}

	err = DatastoreEntityToProtoMessage(&srcEntity, dst, false)
	assert.NilError(t, err)
	fmt.Println("",dst)
	assert.DeepEqual(t, src, dst)
}
