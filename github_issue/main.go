package main

import (
	"log"
	"cloud.google.com/go/datastore"
	"context"
	//pb "google.golang.org/genproto/googleapis/datastore/v1"

	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/translator"
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/example"
)

func main() {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "st2-saas-prototype-dev")
	if err != nil {
		log.Printf("unable to create a client, %v\n", err)
	}
	// working sample don't mess around
	/*p0 := datastore.PropertyList{
		{Name: "L", Value: []interface{}{int64(12), "string", true}},
	}

	k, err := client.Put(ctx, datastore.IncompleteKey("ListValue", nil), &p0)
	if err != nil {
		log.Printf("client.Put: %v", err)
	}
	var p1 datastore.PropertyList
	if err := client.Get(ctx, k, &p1); err != nil {
		log.Printf("client.Get: %v", err)
	}
	log.Printf("property value: %v\n",p1)*/

	simpleProto := &example.ExampleNestedModel{
		Int32Key:10,
		StringKey:"a simple proto message",
	}

	translatedEntity, err := translator.ProtoMessageToDatastoreEntity(simpleProto,false)
	log.Printf("translated proto: %v\n", translatedEntity)
	key := datastore.IncompleteKey("ExampleNestedModel", nil)
	_, err = client.Put(ctx, key, &translatedEntity)
	if err != nil {
		log.Printf("client.Put: %v", err)
	}

	/*var p1 pb.Entity
	if err := client.Get(ctx, k, &p1); err != nil {
		log.Printf("client.Get: %v", err)
	}
	log.Printf("entity from datastore: %v\n",p1)*/
}
