package translator

import (
	"fmt"
	"github.com/Sheshagiri/protobuf-struct/models/execution"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	dbv2 "google.golang.org/api/datastore/v1"
	"testing"
)

func TestTranslateToDatastore(t *testing.T) {
	p := execution.ExecutionRequest{
		Action: "a1",
		Uuid:   "some-uuid",
	}
	TranslateToDatastore(p)
}

func TestProtoToEntity(t *testing.T) {
	p := &execution.ExecutionRequest{
		Action:    "validate",
		Uuid:      "some-random-uuid",
		StartedOn: ptypes.TimestampNow(),
	}
	entity := ProtoToEntity(p)
	assert.NotNil(t, entity.Properties)
	fmt.Println(entity)
}

func TestEntityToProto(t *testing.T) {
	entity := dbv2.Entity{}
	var properties map[string]dbv2.Value

	properties["Action"] = dbv2.Value{
		StringValue: "action-1",
	}
	properties["Uuid"] = dbv2.Value{
		StringValue: "some-random-uuid",
	}
	properties["StartedOn"] = dbv2.Value{
		TimestampValue: ptypes.TimestampString(ptypes.TimestampNow()),
	}
	entity.Properties = properties
	execRequest := &execution.ExecutionRequest{}
	EntityToProto(entity, execRequest)
	assert.Equal(t, "action-1", execRequest.GetAction())
	assert.Equal(t, "some-random-uuid", execRequest.GetUuid())
}
