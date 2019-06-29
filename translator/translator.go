package translator

import (
	"fmt"
	"reflect"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/timestamp"
	// dbv1 "google.golang.org/api/datastore/v1"
	dbv2 "google.golang.org/appengine/datastore"
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
	l := reflect.ValueOf(src).Elem()
	dst := dbv2.Entity{}
	properties := make([]dbv2.Property, 0)
	for i :=0; i < l.NumField(); i++ {
		name := l.Type().Field(i).Name
		if !strings.Contains(name, "XXX_") {
			value := l.Field(i).Interface()
			//fmt.Printf("Type: %v, Name: %v, Value: %v\n", l.Field(i).Type(), name, value)
			properties = append(properties, dbv2.Property{
				Name:name,
				Value: value,
			})
		}
	}
	dst.Properties = properties
	//fmt.Println(dst)
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
				timestamp.Timestamp{

				}
			case "*structpb.Struct":
				fmt.Println("convert to *structpb.Struct")
			}
		}
	}
}

func getProperty(properties []dbv2.Property, name string) interface{} {
	for _, property := range properties {
		if property.Name == name {
			return property.Value
		}
	}
	return nil
}

func getValue(src interface{}, field string) {
	r := reflect.ValueOf(src)
	f := reflect.Indirect(r).FieldByName(field)
	fmt.Println(f.Interface())
}

func decodeValue(v dbv2.Property) interface{}{
	return nil
}