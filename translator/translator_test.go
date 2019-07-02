package translator

import (
	"testing"

	"fmt"

	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/example"
	"github.com/stretchr/testify/assert"
	//dbv2 "google.golang.org/api/datastore/v1"
	//datastore "cloud.google.com/go/datastore"
	//"github.com/golang/protobuf/ptypes/struct"
)

func TestProtoMessageToDatastoreEntitySimple(t *testing.T) {
	srcProto := &example.ExampleNestedModel{
		StringKey: "some random string",
		Int32Key:  22,
	}
	entity := ProtoMessageToDatastoreEntity(srcProto)
	dstProto := &example.ExampleNestedModel{}
	DEtoPM(entity, dstProto)
	assert.Equal(t, srcProto.GetStringKey(), dstProto.GetStringKey())
}

func TestProtoMessageToDatastoreEntityComplex(t *testing.T) {
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
		/*MapStringString: map[string]string{
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
			},
		},*/
	}
	entity := ProtoMessageToDatastoreEntity(srcProto)
	fmt.Println("++++")
	fmt.Println(entity)
	fmt.Println("++++")
	dstProto := &example.ExampleDBModel{}

	DEtoPM(entity, dstProto)
	assert.Equal(t, srcProto.GetStringKey(), dstProto.GetStringKey())
	assert.Equal(t, srcProto.GetBoolKey(), dstProto.GetBoolKey())
	assert.Equal(t, srcProto.GetInt32Key(), dstProto.GetInt32Key())
	assert.Equal(t, srcProto.GetInt64Key(), dstProto.GetInt64Key())
	assert.Equal(t, srcProto.GetFloatKey(), dstProto.GetFloatKey())
	assert.Equal(t, srcProto.GetDoubleKey(), dstProto.GetDoubleKey())
	//TODO BlobValue returns a string
	assert.Equal(t, srcProto.GetBytesKey(), dstProto.GetBytesKey())
	//assert string array
	assert.Equal(t, srcProto.GetStringArrayKey(), dstProto.GetStringArrayKey())
	//assert int32 array
	assert.Equal(t, srcProto.Int32ArrayKey, dstProto.Int32ArrayKey)
	// enums are converted to int's in datastore
	//assert.Equal(t, int64(1), GetProperty(entity.Properties, "EnumKey").Value)
	//assert map[string]string
	assert.Equal(t, srcProto.GetMapStringString(), dstProto.GetMapStringString())
	//assert.Equal(t, int64(1), entity.Properties["MapStringInt32"].EntityValue.Properties["int-key-1"].Value)
	//assert.Equal(t, int64(2), entity.Properties["MapStringInt32"].EntityValue.Properties["int-key-2"].Value)
	//assert google.protobuf.Struct
	/*assert.Equal(t, true, entity.Properties["StructKey"].EntityValue.Properties["struct-key-bool"].BooleanValue)
	assert.Equal(t, "some random string in proto.Struct", entity.Properties["StructKey"].EntityValue.Properties["struct-key-string"].StringValue)
	assert.Equal(t, float64(123456.12), entity.Properties["StructKey"].EntityValue.Properties["struct-key-number"].Value)
	assert.Equal(t, "", entity.Properties["StructKey"].EntityValue.Properties["struct-key-null"].NullValue)*/
}

/*
func TestDatastoreEntityToProtoMessage(t *testing.T) {
	properties := make(map[string]dbv2.Value)
	properties["StringKey"] = dbv2.Value{
		StringValue: "some random string key",
	}
	properties["Int64Key"] = dbv2.Value{
		IntegerValue: 64,
	}
	properties["DoubleKey"] = dbv2.Value{
		DoubleValue: float64(64),
	}
	properties["BoolKey"] = dbv2.Value{
		BooleanValue: false,
	}
	properties["EnumKey"] = dbv2.Value{
		IntegerValue: 2,
	}
	properties["MapStringString"] = dbv2.Value{
		EntityValue: &dbv2.Entity{
			Properties: map[string]dbv2.Value{
				"k1": {StringValue: "some-string-key-1"},
				"k2": {StringValue: "some-string-key-2"},
			},
		},
	}
	properties["MapStringInt32"] = dbv2.Value{
		EntityValue: &dbv2.Entity{
			Properties: map[string]dbv2.Value{
				"int-key-1": {IntegerValue: 10},
				"int-key-2": {IntegerValue: 20},
			},
		},
	}
	properties["StructKey"] = dbv2.Value{
		EntityValue: &dbv2.Entity{
			Properties: map[string]dbv2.Value{
				"struct-key-string": {StringValue: "apple inc"},
				"struct-key-number": {IntegerValue: 20},
				"struct-key-bool":   {BooleanValue: true},
				"struct-key-null":   {NullValue: ""},
			},
		},
	}

	entity := dbv2.Entity{
		Properties: properties,
	}
	dbModel := &example.ExampleDBModel{}
	DatastoreEntityToProtoMessage(entity, dbModel)
	assert.Equal(t, properties["StringKey"].StringValue, dbModel.GetStringKey())
	assert.Equal(t, entity.Properties["Int64Key"].IntegerValue, dbModel.GetInt64Key())
	assert.Equal(t, entity.Properties["BoolKey"].BooleanValue, dbModel.GetBoolKey())
	assert.Equal(t, entity.Properties["DoubleKey"].DoubleValue, dbModel.GetDoubleKey())
	assert.Equal(t, example.ExampleEnumModel_ENUM2, dbModel.GetEnumKey())
	//assert map[string]string
	assert.Equal(t, map[string]string{"k1": "some-string-key-1", "k2": "some-string-key-2"}, dbModel.GetMapStringString())
	//assert map[string]int32
	assert.Equal(t, map[string]int32{"int-key-1": 10, "int-key-2": 20}, dbModel.GetMapStringInt32())
	//assert google.protobuf.Struct
	fmt.Println(dbModel.GetStructKey())
}*/
