package translator

import (
	"testing"
	"github.com/Sheshagiri/protobuf-struct/models/example"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/api/datastore/v1"
	"github.com/stretchr/testify/assert"
)

func TestProtoToEntity2(t *testing.T) {
	exampleProto := &example.ExampleDBModel{
		StringKey:"sample string key",
		BoolKey:true,
		TimestampKey:ptypes.TimestampNow(),
	}
	entityPB := ProtoToEntity(exampleProto)
	fmt.Println(entityPB)
}

func TestEntityToProto2(t *testing.T) {
	entityPB := datastore.Entity{}
	properties := make(map[string]datastore.Value)
	properties["StringKey"] = datastore.Value{
		StringValue:"sample string key",
	}
	properties["BoolKey"] = datastore.Value{
		BooleanValue:true,
	}
	properties["TimestampKey"] = datastore.Value{
		TimestampValue:ptypes.TimestampString(ptypes.TimestampNow()),
	}
	entityPB.Properties = properties
	exampleProto := &example.ExampleDBModel{}
	EntityToProto(entityPB,exampleProto)
	fmt.Println(exampleProto)
}

func TestProtoMessageToDatastoreEntity(t *testing.T) {
	example := &example.ExampleNestedModel{
		StringKey:"some random string",
		Int32Key:22,
	}
	entity := ProtoMessageToDatastoreEntity(example)

	assert.Equal(t, "some random string", entity.Properties["string_key"].StringValue)
	assert.Equal(t, 22, int(entity.Properties["int32_key"].IntegerValue))
}
