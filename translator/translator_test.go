package translator

import (
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/example"
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/unsupported"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/genproto/googleapis/datastore/v1"
	"gotest.tools/assert"
	"log"
	"testing"
	"fmt"
	clientSDK "cloud.google.com/go/datastore"
	"reflect"
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

	assert.DeepEqual(t , srcProto, dstProto)
	assert.Equal(t, srcProto.GetStringKey(), dstProto.GetStringKey())
}

func TestFullyPopulatedModel(t *testing.T) {
	srcProto := &example.ExampleDBModel{
		StringKey: "some random string key for testing",
		BoolKey:   true,
		Int32Key:  int32(32),
		Int64Key:  64,
		FloatKey:  float32(10.14),
		DoubleKey: float64(10.2121),
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
		TimestampKey: ptypes.TimestampNow(),
	}
	entity, err := ProtoMessageToDatastoreEntity(srcProto, true)

	// make sure there is no error
	assert.NilError(t, err)
	log.Println(entity)
	dstProto := &example.ExampleDBModel{}

	err = DatastoreEntityToProtoMessage(&entity, dstProto, true)
	// make sure there is no error
	assert.NilError(t, err)

	assert.Equal(t, srcProto.GetStringKey(), dstProto.GetStringKey())
	assert.Equal(t, srcProto.GetBoolKey(), dstProto.GetBoolKey())
	assert.Equal(t, srcProto.GetInt32Key(), dstProto.GetInt32Key())
	assert.Equal(t, srcProto.GetInt64Key(), dstProto.GetInt64Key())
	assert.Equal(t, srcProto.GetFloatKey(), dstProto.GetFloatKey())
	assert.Equal(t, srcProto.GetDoubleKey(), dstProto.GetDoubleKey())
	//TODO BlobValue returns a string
	assert.DeepEqual(t, srcProto.GetBytesKey(), dstProto.GetBytesKey())
	//assert string array
	assert.DeepEqual(t, srcProto.GetStringArrayKey(), dstProto.GetStringArrayKey())
	//assert int32 array
	assert.DeepEqual(t, srcProto.Int32ArrayKey, dstProto.Int32ArrayKey)
	// enums are converted to int's in datastore
	assert.Equal(t, srcProto.GetEnumKey(), dstProto.GetEnumKey())
	//assert map[string]string
	//assert.DeepEqual(t, srcProto.GetMapStringString(), dstProto.GetMapStringString())
	//assert.DeepEqual(t, srcProto.GetMapStringInt32(), dstProto.GetMapStringInt32())

	//assert google.protobuf.Struct
	assert.DeepEqual(t, srcProto.GetStructKey(), dstProto.GetStructKey())
	//extra check to see if they are really equal
	assert.Equal(t, srcProto.GetStructKey().Fields["struct-key-string"].GetStringValue(), dstProto.GetStructKey().Fields["struct-key-string"].GetStringValue())

	//assert google.protobuf.timestamp
	assert.DeepEqual(t, srcProto.GetTimestampKey().Seconds, dstProto.GetTimestampKey().Seconds)

	//assert listvalues inside the struct
	assert.DeepEqual(t, srcProto.GetStructKey().Fields["struct-key-list"].GetListValue(), dstProto.GetStructKey().Fields["struct-key-list"].GetListValue())
}

func TestPartialModel(t *testing.T) {
	partialProto := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"struct-key-string": {Kind: &structpb.Value_StringValue{"some random string in proto.Struct"}},
			// not ready for this yet
			// "struct-key-list":   {Kind: &structpb.Value_ListValue{}},
			"struct-key-bool":   {Kind: &structpb.Value_BoolValue{true}},
			"struct-key-number": {Kind: &structpb.Value_NumberValue{float64(123456.12)}},
			"struct-key-null":   {Kind: &structpb.Value_NullValue{}},
		},
	}
	entity, err := ProtoMessageToDatastoreEntity(partialProto, true)
	assert.NilError(t, err, err)
	log.Println(entity)
	dstProto := &structpb.Struct{}
	err = DatastoreEntityToProtoMessage(&entity, dstProto, true)
	assert.NilError(t, err)

	//assert google.protobuf.Struct
	assert.DeepEqual(t, partialProto.Fields["struct-key-string"], dstProto.Fields["struct-key-string"])
}

func TestUnSupportedTypes(t *testing.T) {
	srcProto := &unsupported.Model{
		Uint32Key: uint32(10),
	}
	_, err := ProtoMessageToDatastoreEntity(srcProto, false)
	assert.Error(t, err, "[toDatastoreValue]: datatype[uint32] not supported")

	dstProto := &unsupported.Model{}
	datastoreEntity := &datastore.Entity{}
	err = DatastoreEntityToProtoMessage(datastoreEntity, dstProto, false)
	assert.Error(t, err, "datatype[uint32] not supported")
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
							"struct-key-string": {Kind: &structpb.Value_StringValue{"some random string in proto.Struct"}},
							// not ready for this yet
							// "struct-key-list":   {Kind: &structpb.Value_ListValue{}},
							"struct-key-bool":   {Kind: &structpb.Value_BoolValue{true}},
							"struct-key-number": {Kind: &structpb.Value_NumberValue{float64(123456.12)}},
							"struct-key-null":   {Kind: &structpb.Value_NullValue{}},
						},
					},
				},
			},
			datastoreValue: &datastore.Value{
				ValueType: &datastore.Value_EntityValue{
					EntityValue: &datastore.Entity{
						Properties: map[string]*datastore.Value{
							"struct-key-string": {ValueType: &datastore.Value_StringValue{"some random string in proto.Struct"}},
							// not ready for this yet
							// "struct-key-list":   {ValueType: &datastore.Value_ArrayValue{}},
							"struct-key-bool":   {ValueType: &datastore.Value_BooleanValue{true}},
							"struct-key-number": {ValueType: &datastore.Value_DoubleValue{float64(123456.12)}},
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
				StringKey: "string from second element",
			},
			{
				Int32Key:  1,
				StringKey: "string in second element",
			},
		},
	}
	_, err := ProtoMessageToDatastoreEntity(srcProto, false)
	assert.Error(t, err, "[toDatastoreValue]: datatype[*example.ExampleNestedModel] not supported")
}

func TestEmptyMessage(t *testing.T){
	dst := &unsupported.Child{}
	fmt.Println("Before: ", dst)
	clientSDK.EntityToStruct(dst,nil)
	fmt.Println("After: ",dst)
	dstValues := reflect.ValueOf(dst).Elem()
	for i:=0;i<dstValues.NumField();i++{
		name := dstValues.Type().Field(i).Name
		fmt.Println(name)
	}
}
