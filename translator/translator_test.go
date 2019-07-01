package translator

import (
	"testing"

	"fmt"

	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/example"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/datastore/v1"
)

func TestProtoMessageToDatastoreEntitySimple(t *testing.T) {
	example := &example.ExampleNestedModel{
		StringKey: "some random string",
		Int32Key:  22,
	}
	entity := ProtoMessageToDatastoreEntity(example)

	assert.Equal(t, "some random string", entity.Properties["StringKey"].StringValue)
	assert.Equal(t, 22, int(entity.Properties["Int32Key"].IntegerValue))
}

func TestProtoMessageToDatastoreEntityComplex(t *testing.T) {
	float32 := float32(10.1)
	e1 := &example.ExampleDBModel{
		StringKey: "some random string key for testing",
		BoolKey:   true,
		Int32Key:  32,
		Int64Key:  64,
		FloatKey:  float32,
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
			},
		},
	}
	entity := ProtoMessageToDatastoreEntity(e1)

	assert.Equal(t, "some random string key for testing", entity.Properties["StringKey"].StringValue)
	assert.Equal(t, true, entity.Properties["BoolKey"].BooleanValue)
	assert.Equal(t, int64(32), entity.Properties["Int32Key"].IntegerValue)
	assert.Equal(t, int64(64), entity.Properties["Int64Key"].IntegerValue)
	assert.Equal(t, float64(float32), entity.Properties["FloatKey"].DoubleValue)
	assert.Equal(t, float64(10.2121), entity.Properties["DoubleKey"].DoubleValue)
	//TODO BlobValue returns a string
	assert.Equal(t, string([]byte("this is a byte array")), entity.Properties["BytesKey"].BlobValue)
	//assert string array
	assert.Equal(t, "element-1", entity.Properties["StringArrayKey"].ArrayValue.Values[0].StringValue)
	assert.Equal(t, "element-2", entity.Properties["StringArrayKey"].ArrayValue.Values[1].StringValue)
	//assert int32 array
	assert.Equal(t, int64(1), entity.Properties["Int32ArrayKey"].ArrayValue.Values[0].IntegerValue)
	assert.Equal(t, int64(3), entity.Properties["Int32ArrayKey"].ArrayValue.Values[2].IntegerValue)
	assert.Equal(t, int64(5), entity.Properties["Int32ArrayKey"].ArrayValue.Values[4].IntegerValue)
	assert.Equal(t, int64(6), entity.Properties["Int32ArrayKey"].ArrayValue.Values[5].IntegerValue)
	// enums are converted to int's in datastore
	assert.Equal(t, int64(1), entity.Properties["EnumKey"].IntegerValue)
	//assert map[string]string
	assert.Equal(t, "k1", entity.Properties["MapStringString"].EntityValue.Properties["string-key-1"].StringValue)
	assert.Equal(t, "k2", entity.Properties["MapStringString"].EntityValue.Properties["string-key-2"].StringValue)
	assert.Equal(t, int64(1), entity.Properties["MapStringInt32"].EntityValue.Properties["int-key-1"].IntegerValue)
	assert.Equal(t, int64(2), entity.Properties["MapStringInt32"].EntityValue.Properties["int-key-2"].IntegerValue)
	//assert google.protobuf.Struct
	assert.Equal(t, true, entity.Properties["StructKey"].EntityValue.Properties["struct-key-bool"].BooleanValue)
	assert.Equal(t, "some random string in proto.Struct", entity.Properties["StructKey"].EntityValue.Properties["struct-key-string"].StringValue)
	assert.Equal(t, float64(123456.12), entity.Properties["StructKey"].EntityValue.Properties["struct-key-number"].DoubleValue)
	assert.Equal(t, "", entity.Properties["StructKey"].EntityValue.Properties["struct-key-null"].NullValue)
}

func TestDatastoreEntityToProtoMessage(t *testing.T) {
	properties := make(map[string]datastore.Value)
	properties["StringKey"] = datastore.Value{
		StringValue: "some random string key",
	}
	properties["Int64Key"] = datastore.Value{
		IntegerValue: 64,
	}
	properties["DoubleKey"] = datastore.Value{
		DoubleValue: float64(64),
	}
	properties["BoolKey"] = datastore.Value{
		BooleanValue: false,
	}
	properties["EnumKey"] = datastore.Value{
		IntegerValue: 2,
	}
	properties["MapStringString"] = datastore.Value{
		EntityValue: &datastore.Entity{
			Properties: map[string]datastore.Value{
				"k1": {StringValue: "some-string-key-1"},
				"k2": {StringValue: "some-string-key-2"},
			},
		},
	}
	properties["MapStringInt32"] = datastore.Value{
		EntityValue: &datastore.Entity{
			Properties: map[string]datastore.Value{
				"int-key-1": {IntegerValue: 10},
				"int-key-2": {IntegerValue: 20},
			},
		},
	}
	properties["StructKey"] = datastore.Value{
		EntityValue: &datastore.Entity{
			Properties: map[string]datastore.Value{
				"struct-key-string": {StringValue: "apple inc"},
				"struct-key-number": {IntegerValue: 20},
				"struct-key-bool":   {BooleanValue: true},
				"struct-key-null":   {NullValue: ""},
			},
		},
	}

	entity := datastore.Entity{
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
}
