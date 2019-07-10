package translator

import (
	"testing"
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/unsupported"
	"github.com/golang/protobuf/ptypes"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestAddTSSupport(t *testing.T) {
	src := &unsupported.TS{
		StartedOn:ptypes.TimestampNow(),
	}
	fmt.Println("Source: ",src)
	srcEntity,err := ProtoMessageToDatastoreEntity(src, false)

	assert.NoError(t, err)
	fmt.Println("Source Datastore Entity: ",srcEntity)

	dst := &unsupported.TS{}

	err = DatastoreEntityToProtoMessage(&srcEntity, dst, false)
	assert.NoError(t, err)
	fmt.Println("Destination: ",dst.GetStartedOn())
}
