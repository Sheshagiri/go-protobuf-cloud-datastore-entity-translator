package main

import (
	"cloud.google.com/go/datastore"
	"context"
	"log"
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/example"
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/translator"
)

func main() {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "st2-saas-prototype-dev")
	if err != nil {
		log.Printf("error: %v", err)
	}

	key := datastore.NameKey("ExampleNestedModel", "proto-2", nil)
	simpleProto := &example.ExampleNestedModel{
		StringKey:"some random string key in simple proto",
		Int32Key:11,
	}

	translatedSimpleProto, _ := translator.ProtoMessageToDatastoreEntity(simpleProto, true)

	_, err = client.PutEntity(ctx, key,&translatedSimpleProto)

	if err != nil {
		log.Printf("error while saving: %v", err)
	}

	key2 := datastore.NameKey("Example_DB_Model", "complex_proto_1", nil)
	/*complexProto := &example.ExampleDBModel{
		StringKey: "some random string key for testing",
		BoolKey:   true,
		Int32Key:  int32(32),
		Int64Key:  64,
		FloatKey:  float32(10.14),
		DoubleKey: float64(10.2121),
		BytesKey:  []byte("this is a byte array"),
		StringArrayKey: []string{
			"element-1",
			"element-2",
		},
		Int32ArrayKey: []int32{
			1, 2, 3, 4, 5, 6,
		},
		EnumKey: example.ExampleEnumModel_ENUM1,
		MapStringString: map[string]string{
			"string-key-1": "k1",
			"string-key-2": "k2",
		},
		MapStringInt32: map[string]int32{
			"int-key-1": 1,
			"int-key-2": 2,
		},
		StructKey: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"struct-key-string": {Kind: &structpb.Value_StringValue{"some random string in proto.Struct"}},
				"struct-key-bool":   {Kind: &structpb.Value_BoolValue{true}},
				"struct-key-number": {Kind: &structpb.Value_NumberValue{float64(123456.12)}},
				"struct-key-null":   {Kind: &structpb.Value_NullValue{}},
			},
		},
		TimestampKey: ptypes.TimestampNow(),
	}

	translatedComplexProto, _ := translator.ProtoMessageToDatastoreEntity(complexProto, true)
	_, err = client.PutEntity(ctx, key2,&translatedComplexProto)

	if err != nil {
		log.Printf("error while saving complex proto: %v", err)
	}*/
	datastoreEntity, err := client.GetEntity(ctx, key2)
	if err!= nil {
		log.Printf("error while fetching entity: %v",err)
	}
	log.Printf("entity from cloud datastore: %v",datastoreEntity)
	reverseComplexProto := &example.ExampleDBModel{}
	err = translator.DatastoreEntityToProtoMessage(datastoreEntity, reverseComplexProto, false)
	log.Printf("datastore entity to proto message: %v",reverseComplexProto)
}
