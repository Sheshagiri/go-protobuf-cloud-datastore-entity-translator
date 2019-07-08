package translator

import (
	"testing"
	"google.golang.org/genproto/googleapis/datastore/v1"
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/unsupported"
)

func TestDEtoPM(t *testing.T) {
	childDE := datastore.Entity{
		Properties:map[string]*datastore.Value{
			"parent":{
				ValueType:&datastore.Value_EntityValue{
					EntityValue: &datastore.Entity{
						Properties:map[string]*datastore.Value{
							"name":{ValueType:&datastore.Value_StringValue{StringValue:"Parent"}},
						},
					} ,
				},
			},
			"name":{ ValueType:&datastore.Value_StringValue{StringValue:"Child"}},
		},
	}
	childPM := unsupported.Child{}
	DEtoPM(&childDE, &childPM)
}