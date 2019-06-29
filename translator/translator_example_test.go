package translator

import (
	"fmt"
	"github.com/Sheshagiri/protobuf-struct/models/example"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/datastore/v1"
	"testing"
)

func TestProtoToEntity2(t *testing.T) {
	exampleProto := &example.ExampleDBModel{
		StringKey:    "sample string key",
		BoolKey:      true,
		TimestampKey: ptypes.TimestampNow(),
	}
	entityPB := ProtoToEntity(exampleProto)
	fmt.Println(entityPB)
}

func TestEntityToProto2(t *testing.T) {
	entityPB := datastore.Entity{}
	properties := make(map[string]datastore.Value)
	properties["StringKey"] = datastore.Value{
		StringValue: "sample string key",
	}
	properties["BoolKey"] = datastore.Value{
		BooleanValue: true,
	}
	properties["TimestampKey"] = datastore.Value{
		TimestampValue: ptypes.TimestampString(ptypes.TimestampNow()),
	}
	entityPB.Properties = properties
	exampleProto := &example.ExampleDBModel{}
	EntityToProto(entityPB, exampleProto)
	fmt.Println(exampleProto)
}

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
	}
	/*StringArrayKey:[]string{
		"element-1",
		"element-2",
	},*/
	/*MapStringInt32:map[string]int32{
		"int-key-1":1,
		"int-key-2":2,
	},
	MapStringString:map[string]string{
		"string-key-1":"k1",
		"string-key-2":"k2",
	},
	Int32ArrayKey:[]int32{
		1,2,3,4,5,6,
	},*/
	entity := ProtoMessageToDatastoreEntity(e1)

	assert.Equal(t, "some random string key for testing", entity.Properties["StringKey"].StringValue)
	assert.Equal(t, true, entity.Properties["BoolKey"].BooleanValue)
	assert.Equal(t, int64(32), entity.Properties["Int32Key"].IntegerValue)
	assert.Equal(t, int64(64), entity.Properties["Int64Key"].IntegerValue)
	assert.Equal(t, float64(float32), entity.Properties["FloatKey"].DoubleValue)
	assert.Equal(t, float64(10.2121), entity.Properties["DoubleKey"].DoubleValue)
}
