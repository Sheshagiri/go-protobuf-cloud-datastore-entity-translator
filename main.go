package main

import (
	"github.com/Sheshagiri/protobuf-struct/models"
	"github.com/golang/protobuf/ptypes"
	"log"
	//stpb "github.com/golang/protobuf/ptypes/struct"
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"os"
	"reflect"
)

var db *datastore.Client
var ctx context.Context

func setupDatastore() {
	var err error
	ctx = context.Background()
	projectID := os.Getenv("DATASTORE_PROJECT_ID")
	if projectID == "" {
		log.Fatal(`set the environment variable "DATASTORE_PROJECT_ID"`)
	}

	if db, err = datastore.NewClient(ctx, projectID); err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

}

func main() {
	setupDatastore()
	execRequest1 := execution.ExecutionRequest{
		Action:    "a1",
		Uuid:      "uuuu-aaaa",
		StartedOn: ptypes.TimestampNow(),
		//Parameters:[]byte("{\"key-3\": \"value-3\",\"key-4\": 4}"),
	}
	log.Printf("execution request1: %v\n", execRequest1)
	jsonm := &jsonpb.Marshaler{Indent: " "}
	jsonm.Marshal(os.Stdout, &execRequest1)

	/*key1 := datastore.NameKey("st2-saas-scheduler",execRequest1.Action, nil)
	log.Println(key1)
	if _, err := db.Put(ctx, key1, &execRequest1); err != nil {
		log.Printf("unable to save, %v", err)
	}*/

	json := `{"uuid": "uuuu-bbbb","action": "a2","startedOn": "2019-06-25T05:21:49.077578Z","parameters": {"key-3": "value-3","key-4": 4}}`
	// execRequest2 := execution.ExecutionRequest{}
	execRequest2 := execution.ExecutionRequest{}
	if err := jsonpb.UnmarshalString(json, &execRequest2); err != nil {
		log.Printf("Unmarshalling failed, %v\n", err)
	}
	log.Printf("execution request2: %v\n", execRequest2)
	jsonm.Marshal(os.Stdout, &execRequest2)
	key2 := datastore.NameKey("st2-saas-scheduler", "a2", nil)

	props, err := datastore.SaveStruct(&execRequest2)
	if err != nil {
		log.Printf("unable to save struct, %v", err)
	}
	fmt.Printf("properties are: %v", props)
	if _, err := db.Put(ctx, key2, props); err != nil {
		log.Printf("unable to save, %v", err)
	}
}

func proto_to_entity(src proto.Message) {
	//entity := datastore.Entity{}
	t := reflect.TypeOf(src)
	for i := 0; i < t.NumField(); i++ {
		fmt.Printf("%+v\n", t.Field(i))
	}
	log.Println(src)
}
