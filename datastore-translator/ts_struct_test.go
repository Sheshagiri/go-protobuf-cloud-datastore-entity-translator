package datastore_translator

import (
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/example"
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/unsupported"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/struct"
	"gotest.tools/assert"
	"log"
	"testing"
)

func TestAddTSSupport(t *testing.T) {
	src := &unsupported.TS{
		StartedOn: ptypes.TimestampNow(),
	}
	log.Println("Source: ", src)
	srcEntity, err := ProtoMessageToDatastoreEntity(src, false, nil)

	assert.NilError(t, err)
	log.Println("Source Datastore Entity: ", srcEntity)

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
	log.Println("Source: ", src)
	srcEntity, err := ProtoMessageToDatastoreEntity(src, false, nil)

	assert.NilError(t, err)
	log.Println("Source Datastore Entity: ", srcEntity)

	dst := &unsupported.StructMessage{}

	err = DatastoreEntityToProtoMessage(&srcEntity, dst, false)
	assert.NilError(t, err)
	log.Println("", dst)
	assert.DeepEqual(t, src, dst)
}

func TestSliceofNestedMessages(t *testing.T) {
	src := &example.ExampleDBModel{
		ComplexArrayKey: []*example.ExampleNestedModel{
			{
				StringKey: "string-1",
			},
			{
				StringKey: "string-2",
			},
		},
	}
	log.Println("Source: ", src)
	srcEntity, err := ProtoMessageToDatastoreEntity(src, false, nil)
	assert.NilError(t, err)

	log.Println("Source Datastore Entity: ", srcEntity)

	dst := &example.ExampleDBModel{}

	err = DatastoreEntityToProtoMessage(&srcEntity, dst, false)
	assert.NilError(t, err)
	log.Println("Destination: ", dst)
	assert.DeepEqual(t, src.GetComplexArrayKey(), dst.GetComplexArrayKey())

}

func TestNestedMessages(t *testing.T) {
	src := &unsupported.Child{
		Name: "Alex-II",
		Parent: &unsupported.Parent{
			Name: "Alex-I",
		},
	}
	log.Println("Source: ", src)
	srcEntity, err := ProtoMessageToDatastoreEntity(src, false, nil)
	assert.NilError(t, err)

	log.Println("Source Datastore Entity: ", srcEntity)

	dst := &unsupported.Child{}

	err = DatastoreEntityToProtoMessage(&srcEntity, dst, false)
	assert.NilError(t, err)
	log.Println("Destination: ", dst)
	assert.DeepEqual(t, src, dst)
}
