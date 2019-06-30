package translator

import (
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/example"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/datastore/v1"
	"testing"
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
	}
	/*MapStringInt32:map[string]int32{
		"int-key-1":1,
		"int-key-2":2,
	},
	MapStringString:map[string]string{
		"string-key-1":"k1",
		"string-key-2":"k2",
	},*/
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
	entity := datastore.Entity{
		Properties: properties,
	}
	example := &example.ExampleDBModel{}
	DatastoreEntityToProtoMessage(entity, example)
	assert.Equal(t, properties["StringKey"].StringValue, example.GetStringKey())
	assert.Equal(t, entity.Properties["Int64Key"].IntegerValue, example.GetInt64Key())
	assert.Equal(t, entity.Properties["BoolKey"].BooleanValue, example.GetBoolKey())
	assert.Equal(t, float64(64), entity.Properties["DoubleKey"].DoubleValue)
}
