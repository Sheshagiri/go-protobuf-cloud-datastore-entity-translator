package translator

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/struct"
	"reflect"
	//"github.com/golang/protobuf/ptypes/timestamp"
	dbv2 "google.golang.org/api/datastore/v1"
	//dbv2 "google.golang.org/appengine/datastore"
	//"github.com/Sheshagiri/protobuf-struct/models"
	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/proto"
	"strings"
)

// TranslateToDatastore will transform an interface to Entity Protobuf that Datastore understands
func TranslateToDatastore(src interface{}) {

	t := reflect.TypeOf(src)
	//v := reflect.ValueOf(src)
	for i := 0; i < t.NumField(); i++ {
		fmt.Printf("%+v\n", t.Field(i))
		if t.Field(i).Type.Kind() == reflect.TypeOf(structpb.Struct{}).Kind() {
			fmt.Println("yeyy")
		}
	}
}

func ProtoToEntity(src proto.Message) dbv2.Entity {
	srcValues := reflect.ValueOf(src).Elem()
	dst := dbv2.Entity{}
	properties := make(map[string]dbv2.Value)
	for i := 0; i < srcValues.NumField(); i++ {
		name := srcValues.Type().Field(i).Name
		if !strings.Contains(name, "XXX_") {
			value := srcValues.Field(i).Interface()
			//fmt.Printf("Type: %v, Name: %v, Value: %v\n", l.Field(i).Type(), name, value)
			fmt.Println(srcValues.Type().Field(i).Type.String())
			switch srcValues.Type().Field(i).Type.String() {
			case "bool":
				properties[name] = dbv2.Value{
					BooleanValue: value.(bool),
				}
			case "string":
				fmt.Println("string value")
				properties[name] = dbv2.Value{
					StringValue: fmt.Sprint(value),
				}
			case "*timestamp.Timestamp":
				fmt.Println("*timestamp.Timestamp")
				properties[name] = dbv2.Value{
					TimestampValue: fmt.Sprint(value),
				}
			case "*structpb.Struct":
				fmt.Println("convert to *structpb.Struct")
			}
		}
	}
	dst.Properties = properties
	return dst
}

func EntityToProto(src dbv2.Entity, dst proto.Message) {
	dstValues := reflect.ValueOf(dst).Elem()
	for i := 0; i < dstValues.NumField(); i++ {
		fieldName := dstValues.Type().Field(i).Name
		if !strings.Contains(fieldName, "XXX_") {
			fieldValue := getProperty(src.Properties, fieldName)
			fmt.Println(dstValues.Type().Field(i).Type.String())
			switch dstValues.Type().Field(i).Type.String() {
			case "string":
				fmt.Println("type: string")
				fmt.Println(fieldValue)
				dstValues.Field(i).SetString(fmt.Sprint(fieldValue))
			case "*timestamp.Timestamp":
				fmt.Println(fieldValue)
			case "*structpb.Struct":
				fmt.Println("convert to *structpb.Struct")
			}
		}
	}
}

func getProperty(properties map[string]dbv2.Value, name string) dbv2.Value {
	return properties[name]
}

func getValue(src interface{}, field string) {
	r := reflect.ValueOf(src)
	f := reflect.Indirect(r).FieldByName(field)
	fmt.Println(f.Interface())
}

// ProtoMessageToDatastoreEntity will generate an Entity Protobuf that datastore understands
func ProtoMessageToDatastoreEntity(src proto.Message) dbv2.Entity {
	srcValues := reflect.ValueOf(src).Elem()
	_, md := descriptor.ForMessage(src.(descriptor.Message))
	entity := dbv2.Entity{}
	properties := make(map[string]dbv2.Value)
	for _, field := range md.GetField() {
		fieldName := field.GetName()
		fieldType := field.GetType().String()
		fieldValue := srcValues.Field(int(field.GetNumber()) - 1)
		fmt.Printf("type: %s, name: %s, value: %v\n", fieldType, fieldName, fieldValue)

		switch fieldType {
		case "TYPE_DOUBLE":
			fmt.Println("TYPE_DOUBLE")
		case "TYPE_FLOAT":
			fmt.Println("TYPE_FLOAT")
		case "TYPE_UINT64":
			fmt.Println("TYPE_UINT64")
		case "TYPE_INT32", "TYPE_INT64":
			fmt.Println("TYPE_INT32")
			properties[fieldName] = dbv2.Value{
				IntegerValue: fieldValue.Int(),
			}
		case "TYPE_FIXED64":
			fmt.Println("TYPE_FIXED64")
		case "TYPE_FIXED32":
			fmt.Println("TYPE_FIXED32")
		case "TYPE_BOOL":
			fmt.Println("TYPE_BOOL")
			properties[fieldName] = dbv2.Value{
				BooleanValue: fieldValue.Bool(),
			}
		case "TYPE_STRING":
			fmt.Println("TYPE_STRING")
			properties[fieldName] = dbv2.Value{
				StringValue: fieldValue.String(),
			}
		case "TYPE_GROUP":
			fmt.Println("TYPE_GROUP")
		case "TYPE_MESSAGE":
			fmt.Println("TYPE_MESSAGE")
		case "TYPE_BYTES":
			fmt.Println("TYPE_BYTES")
		case "TYPE_UINT32":
			fmt.Println("TYPE_UINT32")
		case "TYPE_ENUM":
			fmt.Println("TYPE_ENUM")
		case "TYPE_SFIXED32":
			fmt.Println("TYPE_SFIXED32")
		case "TYPE_SFIXED64":
			fmt.Println("TYPE_SFIXED64")
		case "TYPE_SINT32", "TYPE_SINT64":
			fmt.Println("TYPE_SINT32")
		default:
			fmt.Println("inside default case")
		}
		fmt.Println("----")
	}
	entity.Properties = properties
	fmt.Println(entity)
	return entity
}
