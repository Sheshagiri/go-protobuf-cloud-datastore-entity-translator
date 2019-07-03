package translator

import (
	"testing"

	"cloud.google.com/go/datastore"
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/example"
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/unsupported"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/stretchr/testify/assert"
	"log"
)

func TestGetProperty(t *testing.T) {
	properties := []datastore.Property{
		{
			Name:  "google",
			Value: "search engine",
		},
	}
	assert.Equal(t, "search engine", GetProperty(properties, "google").(string))
	assert.Nil(t, GetProperty(properties, "apple"))
}

func TestNestedModel(t *testing.T) {
	srcProto := &example.ExampleNestedModel{
		StringKey: "some random string",
		Int32Key:  22,
	}
	entity, err := ProtoMessageToDatastoreEntity(srcProto, true)
	// make sure there is no error
	assert.NoError(t, err)
	dstProto := &example.ExampleNestedModel{}
	err = DatastoreEntityToProtoMessage(entity, dstProto, true)
	// make sure there is no error
	assert.NoError(t, err)

	assert.Equal(t, srcProto.GetStringKey(), dstProto.GetStringKey())
}

func TestFullyPopulatedModel(t *testing.T) {
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
		TimestampKey: ptypes.TimestampNow(),
	}
	entity, err := ProtoMessageToDatastoreEntity(srcProto, true)

	// make sure there is no error
	assert.NoError(t, err)
	log.Println(entity)
	dstProto := &example.ExampleDBModel{}

	err = DatastoreEntityToProtoMessage(entity, dstProto, true)
	// make sure there is no error
	assert.NoError(t, err)

	assert.Equal(t, srcProto.GetStringKey(), dstProto.GetStringKey())
	assert.Equal(t, srcProto.GetBoolKey(), dstProto.GetBoolKey())
	assert.Equal(t, srcProto.GetInt32Key(), dstProto.GetInt32Key())
	assert.Equal(t, srcProto.GetInt64Key(), dstProto.GetInt64Key())
	assert.Equal(t, srcProto.GetFloatKey(), dstProto.GetFloatKey())
	assert.Equal(t, srcProto.GetDoubleKey(), dstProto.GetDoubleKey())
	//TODO BlobValue returns a string
	assert.Equal(t, srcProto.GetBytesKey(), dstProto.GetBytesKey())
	//assert string array
	assert.Equal(t, srcProto.GetStringArrayKey(), dstProto.GetStringArrayKey())
	//assert int32 array
	assert.Equal(t, srcProto.Int32ArrayKey, dstProto.Int32ArrayKey)
	// enums are converted to int's in datastore
	assert.Equal(t, srcProto.GetEnumKey(), dstProto.GetEnumKey())
	//assert map[string]string
	assert.Equal(t, srcProto.GetMapStringString(), dstProto.GetMapStringString())
	assert.Equal(t, srcProto.GetMapStringInt32(), dstProto.GetMapStringInt32())

	//assert google.protobuf.Struct
	assert.Equal(t, srcProto.GetStructKey(), dstProto.GetStructKey())
	//extra check to see if they are really equal
	assert.Equal(t, srcProto.GetStructKey().Fields["struct-key-string"].GetStringValue(), dstProto.GetStructKey().Fields["struct-key-string"].GetStringValue())

	//assert google.protobuf.timestamp
	assert.Equal(t, srcProto.GetTimestampKey().Seconds, dstProto.GetTimestampKey().Seconds)
}

func TestPartialModel(t *testing.T) {
	partialProto := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"struct-key-string": {Kind: &structpb.Value_StringValue{"some random string in proto.Struct"}},
			"struct-key-list":   {Kind: &structpb.Value_ListValue{}},
		},
	}
	entity, err := ProtoMessageToDatastoreEntity(partialProto, true)
	assert.NoError(t, err, err)
	log.Println(entity)
	dstProto := &structpb.Struct{}
	err = DatastoreEntityToProtoMessage(entity, dstProto, true)
	assert.Error(t, err)
}

func TestUnSupportedTypes(t *testing.T) {
	srcProto := &unsupported.Model{
		Uint32Key: uint32(10),
	}
	_, err := ProtoMessageToDatastoreEntity(srcProto, false)
	assert.EqualError(t, err, "datatype[uint32] not supported")

	entity := datastore.Entity{
		Properties: []datastore.Property{
			{
				Name:  "uint32",
				Value: uint32(10),
			},
		},
	}
	err = DatastoreEntityToProtoMessage(entity, srcProto, false)
	assert.EqualError(t, err, "datatype[uint32] not supported")
}
