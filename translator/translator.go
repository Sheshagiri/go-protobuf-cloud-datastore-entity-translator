package translator

import (
	"fmt"
	"reflect"
	"github.com/golang/protobuf/ptypes/struct"
	//"github.com/golang/protobuf/ptypes/timestamp"
	dbv2 "google.golang.org/api/datastore/v1"
	//dbv2 "google.golang.org/appengine/datastore"
	//"github.com/Sheshagiri/protobuf-struct/models"
	"github.com/golang/protobuf/proto"
	"strings"
)

func TranslateToDatastore(src interface{}) {

	t := reflect.TypeOf(src)
	//v := reflect.ValueOf(src)
	for i := 0; i < t.NumField(); i++ {
		fmt.Printf("%+v\n", t.Field(i))
		if t.Field(i).Type.Kind() == reflect.TypeOf(structpb.Struct{}).Kind()   {
			fmt.Println("yeyy")
		}
	}
}


func ProtoToEntity(src proto.Message) dbv2.Entity {
	srcValues := reflect.ValueOf(src).Elem()
	dst := dbv2.Entity{}
	properties := make(map[string]dbv2.Value)
	for i :=0; i < srcValues.NumField(); i++ {
		name := srcValues.Type().Field(i).Name
		if !strings.Contains(name, "XXX_") {
			value := srcValues.Field(i).Interface()
			//fmt.Printf("Type: %v, Name: %v, Value: %v\n", l.Field(i).Type(), name, value)
			switch srcValues.Type().Field(i).Type.String() {
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
	for i :=0; i < dstValues.NumField(); i++ {
		fieldName := dstValues.Type().Field(i).Name
		if !strings.Contains(fieldName, "XXX_") {
			fieldValue := getProperty(src.Properties,fieldName)
			fmt.Println(dstValues.Type().Field(i).Type.String())
			switch dstValues.Type().Field(i).Type.String() {
			case "string":
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