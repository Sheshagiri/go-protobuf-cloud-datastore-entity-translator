package datastore_translator

import (
	"log"
	"testing"

	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/example"
	execution "github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/execution"
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/unsupported"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/struct"
	"gotest.tools/assert"
)

func TestAddTSSupport(t *testing.T) {
	src := &unsupported.TS{
		StartedOn: ptypes.TimestampNow(),
	}
	log.Println("Source: ", src)
	srcEntity, err := ProtoMessageToDatastoreEntity(src, false)

	assert.NilError(t, err)
	log.Println("Source Datastore Entity: ", srcEntity)

	dst, err := DatastoreEntityToProtoMessage(&srcEntity, &unsupported.TS{}, false)
	assert.NilError(t, err)
	assert.Equal(t, true, proto.Equal(src, dst))
}

func TestAddStructSupport(t *testing.T) {
	src := &unsupported.StructMessage{
		StructKey: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"struct-key-string": {Kind: &structpb.Value_StringValue{StringValue:"some random string in proto.Struct"}},
				"struct-key-bool":   {Kind: &structpb.Value_BoolValue{BoolValue:true}},
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
	}
	log.Println("Source: ", src)
	srcEntity, err := ProtoMessageToDatastoreEntity(src, false)

	assert.NilError(t, err)
	log.Println("Source Datastore Entity: ", srcEntity)

	dst, err := DatastoreEntityToProtoMessage(&srcEntity, &unsupported.StructMessage{}, false)
	assert.NilError(t, err)
	log.Println("", dst)
	assert.Equal(t, true, proto.Equal(src, dst))
}

func TestSliceNestedMessages(t *testing.T) {
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
	srcEntity, err := ProtoMessageToDatastoreEntity(src, false)
	assert.NilError(t, err)

	log.Println("Source Datastore Entity: ", srcEntity)

	dst, err := DatastoreEntityToProtoMessage(&srcEntity, &example.ExampleDBModel{}, false)
	assert.NilError(t, err)
	log.Println("Destination: ", dst)
	assert.Equal(t, true, proto.Equal(src, dst))

}

func TestNestedMessages(t *testing.T) {
	src := &unsupported.Child{
		Name: "Alex-II",
		Parent: &unsupported.Parent{
			Name: "Alex-I",
		},
	}
	log.Println("Source: ", src)
	srcEntity, err := ProtoMessageToDatastoreEntity(src, false)
	assert.NilError(t, err)

	log.Println("Source Datastore Entity: ", srcEntity)

	dst, err := DatastoreEntityToProtoMessage(&srcEntity, &unsupported.Child{}, false)
	assert.NilError(t, err)
	log.Println("Destination: ", dst)
	assert.Equal(t, true, proto.Equal(src, dst))
}

func TestStructInReferencedMessage(t *testing.T) {
	src := &execution.Execution{
		Name:                 "login",
		Action:               &execution.Action{
			Name:                 "ssh",
			Parameters:           &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"host": {Kind: &structpb.Value_StringValue{StringValue:"10.10.10.10"}},
					"port": {Kind: &structpb.Value_NumberValue{NumberValue:123456.12}},
				},
			},
		},
	}
	log.Println("Source: ", src)
	srcEntity, err := ProtoMessageToDatastoreEntity(src, false)
	assert.NilError(t, err)
	log.Println("Source Datastore Entity: ", srcEntity)

	dst, err := DatastoreEntityToProtoMessage(&srcEntity, &execution.Execution{}, false)
	assert.NilError(t, err)
	log.Println("Destination: ", dst)
	assert.Equal(t, true, proto.Equal(src, dst))
}