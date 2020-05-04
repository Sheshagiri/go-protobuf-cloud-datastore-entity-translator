# Translate any proto message to Google Datastore Entity and vice-versa


[![Build Status](https://travis-ci.org/Sheshagiri/go-protobuf-cloud-datastore-entity-translator.svg?branch=master)](https://travis-ci.org/Sheshagiri/go-protobuf-cloud-datastore-entity-translator)
[![codecov](https://codecov.io/gh/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/branch/master/graph/badge.svg)](https://codecov.io/gh/Sheshagiri/go-protobuf-cloud-datastore-entity-translator)
[![DeepSource](https://static.deepsource.io/deepsource-badge-light-mini.svg)](https://deepsource.io/gh/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/?ref=repository-badge)

# Background

This is largely inspired from being able to persist protocol buffers to Google Cloud Datastore. Protobuf messages that the datstore supports are listed [here](https://github.com/googleapis/googleapis/blob/c50d9e822e19e069b7e3758736ea58cb4f35267c/google/datastore/v1/entity.proto#L188).
 
This repository acts as a translator to translate any given ``proto.Message`` to ``datastore.Entity`` that the datastore understands and 
``datastore.Entity`` to any ``proto.Message``.

This repository also addresses some of the limitations that [google-cloud-go](https://github.com/googleapis/google-cloud-go/tree/master/datastore) has.
Currently the go datastore library `google-cloud-go` doesn't support `maps`, `google.protobuf.Struct`, `google.protobuf.Value` types 
and it also doesn't expose `Put` and `Get` functions that operate on `datastore.Entity`. I added the support for being able to use `datastore.Entity` by adding `PutEntity` and `GetEntity` in a fork [here](https://github.com/Sheshagiri/google-cloud-go).

Issues that inspired this solution:

https://github.com/googleapis/google-cloud-go/issues/1474

https://github.com/googleapis/google-cloud-go/issues/1489

https://github.com/googleapis/google-cloud-go/issues/680

Following is an example of the same.

```
package main
import (
	"context"
	"cloud.google.com/go/datastore"
	"log"
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/execution"
	"github.com/golang/protobuf/ptypes"
	"github.com/pborman/uuid"
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
		Uuid:uuid.New(),
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

	// 7. convert the value fetched from datastore to protobuf
	dsExecReq, err := translator.DatastoreEntityToProtoMessage(dsEntity,&execution.ExecutionRequest{}, true)
	if err != nil {
		log.Fatalf("error while converting to proto message, %v", err)
	}

	// 8. simply log it :)
	log.Println(dsExecReq)
}
```

## Support for excluding the fields from being indexed buy Google Cloud Datastore
By default, Google Cloud Datastore automatically indexes each entity (model) property.

Indexing each field (entity property) is usually not desired nor needed. It also has some limitations (for example, 
size of a simple field which is to be indexed is limited to 1500 bytes, etc.). In addition to that, unnecessary
 indexing causes increased storage space consumption.

This library allows you to define which model fields to exclude from index on the field basis utilizing Protobuf field 
options extension.

For example:
```proto
syntax = "proto3";

import "google/protobuf/descriptor.proto";

// Custom Protobuf option which specifies which model fields should be excluded
// from index
// NOTE: Keep in mind that it's important not to change the option name
// ("exclude_from_index") since this library uses that special option name to
// determine if a field should be excluded from index.
extend google.protobuf.FieldOptions {
    bool exclude_from_index = 50000;
}

message ExampleDBModelWithOptions1 {
    string string_key_one = 1 [(exclude_from_index) = true];
    string string_key_two = 2;
    string string_key_three = 3 [(exclude_from_index) = true];
    string string_key_four = 4;
    int32 int32_field_one = 5;
    int32 int32_field_two = 6 [(exclude_from_index) = true];
}
```
In this example, fields string_key_one, string_key_three and int32_field_two won't be
 indexed (https://cloud.google.com/datastore/docs/concepts/indexes#unindexed_properties).

Note: If the extension is imported then `<package>.<extension>` should be passed to the `ProtoMessageToDatastoreEntity`.
Following is an example
```Go
datastoreEntity, err := translator.ProtoMessageToDatastoreEntity(dbModel, true, "models.exclude_from_index")
```
My colleague at work wrote an [equivalent translator in python](https://github.com/Kami/python-protobuf-cloud-datastore-entity-translator).

## Tested with go1.12.6
