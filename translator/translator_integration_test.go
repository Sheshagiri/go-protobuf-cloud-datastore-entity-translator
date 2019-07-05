// +build integration

package translator

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/example"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/struct"
	"gotest.tools/assert"
	"log"
	"testing"
)

func TestIntegration(t *testing.T) {
	ctx := context.Background()
	// 1. create a new datastore client
	client, err := datastore.NewClient(ctx, "st2-saas-prototype-dev")
	assert.NilError(t, err)

	// 2. create a key that we plan to save into
	key := datastore.NameKey("Example_DB_Model", "complex_proto_2", nil)

	// 3. create protobuf
	srcProto := &example.ExampleDBModel{
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
		//TimestampKey: ptypes.TimestampNow(),
	}

	log.Printf("original proto: %v", srcProto)
	// 4. translate the source protobuf to datastore.Entity
	translatedSrcProto, err := ProtoMessageToDatastoreEntity(srcProto, true)
	assert.NilError(t, err)

	// 5. save the translated protobuf to datastore
	_, err = client.PutEntity(ctx, key, &translatedSrcProto)
	assert.NilError(t, err)

	// 6. get the saved protobuf from cloud datastore
	datastoreEntity, err := client.GetEntity(ctx, key)
	assert.NilError(t, err)
	log.Printf("entity from cloud datastore: %v", datastoreEntity)

	// 7. create a protobuf that we plan to decode into
	dstProto := &example.ExampleDBModel{}

	// 8. translate the protobuf from datastore.Entity{} to our own protobuf
	err = DatastoreEntityToProtoMessage(datastoreEntity, dstProto, true)
	assert.NilError(t, err)

	log.Printf("original proto                   : %v", srcProto)
	log.Printf("datastore entity to proto message: %v", dstProto)

	// 9. start validating srcProto and dstProto should be equal
	assert.Equal(t, srcProto.GetStringKey(), dstProto.GetStringKey())
	assert.Equal(t, srcProto.GetBoolKey(), dstProto.GetBoolKey())
	assert.Equal(t, srcProto.GetInt32Key(), dstProto.GetInt32Key())
	assert.Equal(t, srcProto.GetInt64Key(), dstProto.GetInt64Key())
	assert.Equal(t, srcProto.GetFloatKey(), dstProto.GetFloatKey())
	assert.Equal(t, srcProto.GetDoubleKey(), dstProto.GetDoubleKey())

	//assert string array
	assert.DeepEqual(t, srcProto.GetStringArrayKey(), dstProto.GetStringArrayKey())
	//assert int32 array
	assert.DeepEqual(t, srcProto.Int32ArrayKey, dstProto.Int32ArrayKey)
	// enums are converted to int's in datastore
	assert.Equal(t, srcProto.GetEnumKey(), dstProto.GetEnumKey())
	//assert map[string]string
	assert.Equal(t, srcProto.GetMapStringString()["string-key-1"], dstProto.GetMapStringString()["string-key-1"])
	assert.Equal(t, srcProto.GetMapStringInt32()["int-key-2"], dstProto.GetMapStringInt32()["int-key-2"])

	//TODO BlobValue returns a string
	assert.DeepEqual(t, srcProto.GetBytesKey(), dstProto.GetBytesKey())

	//extra check to see if they are really equal
	assert.Equal(t, srcProto.GetStructKey().Fields["struct-key-string"].GetStringValue(), dstProto.GetStructKey().Fields["struct-key-string"].GetStringValue())

	//assert google.protobuf.Struct
	assert.DeepEqual(t, srcProto.GetStructKey().Fields, dstProto.GetStructKey().Fields)

	//assert google.protobuf.timestamp
	assert.DeepEqual(t, srcProto.GetTimestampKey().Seconds, dstProto.GetTimestampKey().Seconds)

}

func TestEmptyProtoMessage(t *testing.T) {
	ctx := context.Background()
	// 1. create a new datastore client
	client, err := datastore.NewClient(ctx, "st2-saas-prototype-dev")
	assert.NilError(t, err)

	// 2. create a key that we plan to save into
	key := datastore.NameKey("Example_DB_Model", "complex_proto_empty", nil)

	srcProto := &example.ExampleDBModel{}
	translatedProto, err := ProtoMessageToDatastoreEntity(srcProto, false)
	assert.NilError(t, err)

	_, err = client.PutEntity(ctx, key, &translatedProto)
	//e expect an error when the whole proto is empty
	assert.Error(t, err, "rpc error: code = Internal desc = grpc: error while marshaling: proto: oneof field has nil value")
}

func TestProtoWithNilPointer(t *testing.T) {
	ctx := context.Background()
	// 1. create a new datastore client
	client, err := datastore.NewClient(ctx, "st2-saas-prototype-dev")
	assert.NilError(t, err)

	// 2. create a key that we plan to save into
	key := datastore.NameKey("Example_DB_Model", "complex_proto_empty", nil)

	srcProto := &example.ExampleDBModel{
		TimestampKey: ptypes.TimestampNow(),
	}
	translatedProto, err := ProtoMessageToDatastoreEntity(srcProto, false)
	assert.NilError(t, err)

	_, err = client.PutEntity(ctx, key, &translatedProto)
}
