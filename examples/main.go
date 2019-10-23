package main

import (
	"context"
	"cloud.google.com/go/datastore"
	"log"
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/execution"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/golang/protobuf/ptypes/struct"
	translator "github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/datastore-translator"
)

func main() {
	ctx := context.Background()

	// 1. create datastore client
	dsClient, err := datastore.NewClient(ctx,"st2-saas-prototype-dev")
	if err != nil {
		log.Fatalf("unable to connect to datastore, error: %v", err)
	}

	// 2. create a protobuf message
	execReq := &execution.ExecutionRequest{
		StartedOn:ptypes.TimestampNow(),
		Uuid:uuid.New().String(),
		Action:"create_vm",
		Parameters: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"vm-name": {Kind: &structpb.Value_StringValue{StringValue:"sheshagiri-vm-1"}},
				"allow-root-login":   {Kind: &structpb.Value_BoolValue{BoolValue:true}},
				"network-interfaces": {Kind: &structpb.Value_NumberValue{NumberValue:2}},
			},
		},
	}

	// 3. translate the protobuf message to the format that datastore understands
	entity, err := translator.ProtoMessageToDatastoreEntity(execReq, true)
	if err != nil {
		log.Fatalf("unable to translate execution request to datastore format, error: %v", err)
	}

	// 4. create a key where we would like to store the message, think of this as a primary key
	parentKey := datastore.NameKey("ExecutionRequest",execReq.GetAction(), nil)
	childKey := datastore.NameKey(execReq.GetAction(),execReq.GetUuid(), parentKey)

	// 5. save it to datastore against the key
	_, err = dsClient.PutEntity(ctx, childKey, &entity)
	if err != nil {
		log.Fatalf("unable to translate execution request to datastore format, error: %v", err)
	}
	log.Printf("key %v is saved to datastore",childKey)

	// 6. Try to retrieve the key
	dsEntity, err := dsClient.GetEntity(ctx, childKey)
	if err != nil {
		log.Fatalf("unable to get %v from datastore", childKey)
	}

	// 7. create an empty protobuf
	dsExecReq := &execution.ExecutionRequest{}

	// 8. convert the value fetched from datastore to protobuf
	err = translator.DatastoreEntityToProtoMessage(dsEntity,dsExecReq, true)
	if err != nil {
		log.Fatalf("error while converting to proto message, %v", err)
	}

	// 9. simply log it :)
	log.Println(dsExecReq)
}
