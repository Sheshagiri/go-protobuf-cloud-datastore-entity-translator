package translator

import (
	"testing"
	"github.com/Sheshagiri/protobuf-struct/models"
	"github.com/golang/protobuf/ptypes"
	dbv2 "google.golang.org/appengine/datastore"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestTranslateToDatastore(t *testing.T) {
	p := execution.ExecutionRequest{
		Action:"a1",
		Uuid:"some-uuid",
	}
	TranslateToDatastore(p)
}

func TestProtoToEntity(t *testing.T) {
	p := &execution.ExecutionRequest{
		Action:"validate",
		Uuid:"some-random-uuid",
		StartedOn:ptypes.TimestampNow(),
	}
	entity := ProtoToEntity(p)
	assert.NotNil(t,entity.Properties)

	execRequest := &execution.ExecutionRequest{}
	EntityToProto(entity, execRequest)
	fmt.Println(execRequest)
}

func TestEntityToProto(t *testing.T) {
	entity := dbv2.Entity{}
	properties := make([]dbv2.Property,0)
	properties = append(properties, dbv2.Property{
		Name:"Action",
		Value:"action-1",
	})
	properties = append(properties, dbv2.Property{
		Name:"Uuid",
		Value:"some-random-uuid",
	})
	properties = append(properties, dbv2.Property{
		Name:"StartedOn",
		Value:ptypes.TimestampNow(),
	})
	entity.Properties = properties
	execRequest := &execution.ExecutionRequest{}
	EntityToProto(entity, execRequest)
	assert.Equal(t,"action-1",execRequest.GetAction())
	assert.Equal(t,"some-random-uuid",execRequest.GetUuid())
}