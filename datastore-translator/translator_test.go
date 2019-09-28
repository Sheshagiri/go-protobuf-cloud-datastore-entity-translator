package datastore_translator

import (
	"log"
	"testing"

	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/example"
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/execution"
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/unsupported"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/datastore/v1"
	"gotest.tools/assert"
)

func TestNestedModel(t *testing.T) {
	srcProto := &example.ExampleNestedModel{
		StringKey: "some random string",
		Int32Key:  22,
	}
	entity, err := ProtoMessageToDatastoreEntity(srcProto, true)
	// make sure there is no error
	assert.NilError(t, err)
	dstProto := &example.ExampleNestedModel{}
	err = DatastoreEntityToProtoMessage(&entity, dstProto, true)
	// make sure there is no error
	assert.NilError(t, err)

	assert.Equal(t, true, proto.Equal(srcProto, dstProto), "before and after translation proto messages should be equal")
}

func TestProtoMessageToDatastoreEntityWithExcludeFieldsFromIndex(t *testing.T) {
	srcProto := &example.ExampleDBModel{
		StringKey: "some random string key for testing",
		BoolKey:   true,
		Int32Key:  int32(32),
		Int64Key:  64,
		FloatKey:  float32(10.14),
		DoubleKey: 10.2121,
		BytesKey:  []byte("this is a byte array"),
	}

	// No exclude from index fields specified
	entity, err := ProtoMessageToDatastoreEntity(srcProto, true)
	assert.NilError(t, err)
	assert.Equal(t, entity.Properties["string_key"].ExcludeFromIndexes, false)
	assert.Equal(t, entity.Properties["bytes_key"].ExcludeFromIndexes, false)
	assert.Equal(t, entity.Properties["float_key"].ExcludeFromIndexes, false)

	assert.Equal(t, entity.Properties["bool_key"].ExcludeFromIndexes, false)
	assert.Equal(t, entity.Properties["int32_key"].ExcludeFromIndexes, false)
	assert.Equal(t, entity.Properties["double_key"].ExcludeFromIndexes, false)
}

func TestFullyPopulatedModel(t *testing.T) {
	srcProto := &example.ExampleDBModel{
		StringKey: "some random string key for testing",
		BoolKey:   true,
		Int32Key:  int32(32),
		Int64Key:  64,
		FloatKey:  float32(10.14),
		DoubleKey: 10.2121,
		BytesKey:  []byte("this is a byte array"),
		StringArrayKey: []string{
			"element-1",
			"element-2",
		},
		Int32ArrayKey: []int32{
			1, 2, 3, 4, 5, 6,
		},
		EnumKey: example.ExampleEnumModel_ENUM1,
		MapStringString: map[string]string{
			"string-key-1": "k1",
			"string-key-2": "k2",
		},
		MapStringInt32: map[string]int32{
			"int-key-1": 1,
			"int-key-2": 2,
		},
		StructKey: &structpb.Struct{
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
		TimestampKey: ptypes.TimestampNow(),
	}
	entity, err := ProtoMessageToDatastoreEntity(srcProto, true)

	// make sure there is no error
	assert.NilError(t, err)
	log.Println(entity)
	dstProto := &example.ExampleDBModel{}

	err = DatastoreEntityToProtoMessage(&entity, dstProto, true)
	// make sure there is no error
	require.NoError(t, err)

	assert.Equal(t, true, proto.Equal(srcProto, dstProto), "proto messages should be equal")
}

func TestPartialModel(t *testing.T) {
	partialProto := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"struct-key-string": {Kind: &structpb.Value_StringValue{StringValue: "some random string in proto.Struct"}},
			// not ready for this yet
			// "struct-key-list":   {Kind: &structpb.Value_ListValue{}},
			"struct-key-bool":   {Kind: &structpb.Value_BoolValue{BoolValue: true}},
			"struct-key-number": {Kind: &structpb.Value_NumberValue{NumberValue: 123456.12}},
			"struct-key-null":   {Kind: &structpb.Value_NullValue{}},
		},
	}
	entity, err := ProtoMessageToDatastoreEntity(partialProto, true)
	assert.NilError(t, err, err)
	log.Println(entity)
	dstProto := &structpb.Struct{}
	err = DatastoreEntityToProtoMessage(&entity, dstProto, true)
	assert.NilError(t, err)
	// assert google.protobuf.Struct
	assert.DeepEqual(t, partialProto.Fields["struct-key-string"], dstProto.Fields["struct-key-string"])
}

func TestUnSupportedTypes(t *testing.T) {
	srcProto := &unsupported.Model{
		Uint32Key: uint32(10),
	}
	_, err := ProtoMessageToDatastoreEntity(srcProto, false)
	assert.Error(t, err, "[toDatastoreValue]: datatype[uint32] not supported")
}

func TestPMtoDE(t *testing.T) {
	srcProto := &example.ExampleNestedModel{
		StringKey: "some random string",
		Int32Key:  22,
	}
	entity, err := ProtoMessageToDatastoreEntity(srcProto, true)
	assert.NilError(t, err)
	log.Println(entity)
}

func TestStructValueDatastoreValue(t *testing.T) {
	tests := []struct {
		structValue    *structpb.Value
		datastoreValue *datastore.Value
	}{
		{
			structValue: &structpb.Value{
				Kind: &structpb.Value_StringValue{
					StringValue: "some random string key for testing.",
				},
			},
			datastoreValue: &datastore.Value{
				ValueType: &datastore.Value_StringValue{
					StringValue: "some random string key for testing.",
				},
			},
		},
		{
			structValue: &structpb.Value{
				Kind: &structpb.Value_BoolValue{
					BoolValue: true,
				},
			},
			datastoreValue: &datastore.Value{
				ValueType: &datastore.Value_BooleanValue{
					BooleanValue: true,
				},
			},
		},
		{
			structValue: &structpb.Value{
				Kind: &structpb.Value_NumberValue{
					NumberValue: 15,
				},
			},
			datastoreValue: &datastore.Value{
				ValueType: &datastore.Value_DoubleValue{
					DoubleValue: float64(15),
				},
			},
		},
		{
			structValue: &structpb.Value{
				Kind: &structpb.Value_NullValue{},
			},
			datastoreValue: &datastore.Value{
				ValueType: &datastore.Value_NullValue{},
			},
		},
		{
			structValue: &structpb.Value{
				Kind: &structpb.Value_ListValue{
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
			datastoreValue: &datastore.Value{
				ValueType: &datastore.Value_ArrayValue{
					ArrayValue: &datastore.ArrayValue{
						Values: []*datastore.Value{
							{ValueType: &datastore.Value_DoubleValue{DoubleValue: 10}},
							{ValueType: &datastore.Value_StringValue{StringValue: "hello, world"}},
							{ValueType: &datastore.Value_BooleanValue{BooleanValue: true}},
							{ValueType: &datastore.Value_DoubleValue{DoubleValue: 200}},
						},
					},
				},
			},
		},
		{
			structValue: &structpb.Value{
				Kind: &structpb.Value_StructValue{
					StructValue: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"struct-key-string": {Kind: &structpb.Value_StringValue{StringValue: "some random string in proto.Struct"}},
							// not ready for this yet
							// "struct-key-list":   {Kind: &structpb.Value_ListValue{}},
							"struct-key-bool":   {Kind: &structpb.Value_BoolValue{BoolValue: true}},
							"struct-key-number": {Kind: &structpb.Value_NumberValue{NumberValue: 123456.12}},
							"struct-key-null":   {Kind: &structpb.Value_NullValue{}},
						},
					},
				},
			},
			datastoreValue: &datastore.Value{
				ValueType: &datastore.Value_EntityValue{
					EntityValue: &datastore.Entity{
						Properties: map[string]*datastore.Value{
							"struct-key-string": {ValueType: &datastore.Value_StringValue{StringValue: "some random string in proto.Struct"}},
							// not ready for this yet
							// "struct-key-list":   {ValueType: &datastore.Value_ArrayValue{}},
							"struct-key-bool":   {ValueType: &datastore.Value_BooleanValue{BooleanValue: true}},
							"struct-key-number": {ValueType: &datastore.Value_DoubleValue{DoubleValue: 123456.12}},
							"struct-key-null":   {ValueType: &datastore.Value_NullValue{}},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		actualDatastoreValue := fromStructValueToDatastoreValue(test.structValue)
		assert.DeepEqual(t, test.datastoreValue, actualDatastoreValue)
		// test the other way around now

		actualStructValue := fromDatastoreValueToStructValue(test.datastoreValue)
		assert.DeepEqual(t, test.structValue, actualStructValue)
	}
}

func TestProtoWithCustomImport(t *testing.T) {
	srcProto := &example.ExampleDBModel{
		ComplexArrayKey: []*example.ExampleNestedModel{
			{
				Int32Key:  0,
				StringKey: "string in first element",
			},
			{
				Int32Key:  1,
				StringKey: "string in second element",
			},
		},
	}

	srcEntity := &datastore.Entity{
		Properties: map[string]*datastore.Value{
			"ComplexArrayKey": {
				ValueType: &datastore.Value_ArrayValue{
					ArrayValue: &datastore.ArrayValue{
						Values: []*datastore.Value{
							{
								ValueType: &datastore.Value_EntityValue{
									EntityValue: &datastore.Entity{
										Properties: map[string]*datastore.Value{
											"Int32Key":  {ValueType: &datastore.Value_IntegerValue{IntegerValue: 0}},
											"StringKey": {ValueType: &datastore.Value_StringValue{StringValue: "string in first element"}},
										},
									},
								},
							},
							{
								ValueType: &datastore.Value_EntityValue{
									EntityValue: &datastore.Entity{
										Properties: map[string]*datastore.Value{
											"Int32Key":  {ValueType: &datastore.Value_IntegerValue{IntegerValue: 1}},
											"StringKey": {ValueType: &datastore.Value_StringValue{StringValue: "string in second element"}},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	dstEntity, err := ProtoMessageToDatastoreEntity(srcProto, false)
	assert.NilError(t, err)
	// our interest here only to compare the ComplexArrayKey
	assert.DeepEqual(t, srcEntity.GetProperties()["ComplexArrayKey"], dstEntity.GetProperties()["ComplexArrayKey"])
}

func TestSlicedMessages(t *testing.T) {
	tests := []proto.Message{
		&execution.Execution{
			Name: "a",
		},
		&execution.Execution{
			Name: "b",
		},
	}
	want := make([]proto.Message, 0)
	model := &execution.Execution{}
	for _, test := range tests {
		dsEntity, err := ProtoMessageToDatastoreEntity(test, true)
		require.NoError(t, err)
		m := proto.Clone(model)
		err = DatastoreEntityToProtoMessage(&dsEntity, m, true)
		require.NoError(t, err)
		want = append(want, m)
	}
	assert.DeepEqual(t, tests, want)
}
