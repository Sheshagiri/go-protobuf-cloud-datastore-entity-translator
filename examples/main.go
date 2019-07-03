package main

import (
	"context"

	"log"

	"cloud.google.com/go/datastore"

	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/example"
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/translator"
)

func main() {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "st2-saas-prototype-dev")
	if err != nil {
		log.Println(err)
	}
	simpleProto := &example.ExampleNestedModel{
		StringKey:"validate",
		Int32Key:204,
	}
	entity, _ := translator.ProtoMessageToDatastoreEntity(simpleProto,true)
	log.Println(entity)
	key := datastore.NameKey("ExampleNestedModel","key-1",nil)
	//entity.Key = key
	_, err = client.Put(ctx,key,&entity)
	log.Println(err)
}
