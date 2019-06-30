package translator

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	dbv2 "google.golang.org/api/datastore/v1"
	"reflect"
	"strings"
)

func getProperty(properties map[string]dbv2.Value, name string) dbv2.Value {
	return properties[name]
}

// ProtoMessageToDatastoreEntity will generate an Entity Protobuf that datastore understands
func ProtoMessageToDatastoreEntity(src proto.Message) dbv2.Entity {
	srcValues := reflect.ValueOf(src).Elem()
	entity := dbv2.Entity{}
	properties := make(map[string]dbv2.Value)

	for i := 0; i < srcValues.NumField(); i++ {
		name := srcValues.Type().Field(i).Name
		if !strings.ContainsAny(name, "XXX_") {
			fType := srcValues.Field(i).Type().Kind().String()
			value, err := toValue(fType, srcValues.Field(i))
			// fmt.Printf("name:%s, type:%v, value:%v\n",name,fType,value)
			if err == nil {
				properties[name] = value
			} else {
				fmt.Printf("err: %v\n", err)
			}
		}
	}
	entity.Properties = properties
	fmt.Println("@@@@@@@")
	fmt.Printf("Entity: %v\n", entity)
	fmt.Println("@@@@@@@")
	return entity
}

func DatastoreEntityToProtoMessage(src dbv2.Entity, dst proto.Message) {
	dstValues := reflect.ValueOf(dst).Elem()
	for i := 0; i < dstValues.NumField(); i++ {
		fName := dstValues.Type().Field(i).Name
		if !strings.Contains(fName, "XXX_") {
			fValue := getProperty(src.Properties, fName)
			fType := dstValues.Type().Field(i).Type.String()
			fmt.Printf("name: %s, type: %s\n", fName, fType)
			switch fType {
			case "string":
				dstValues.Field(i).SetString(fValue.StringValue)
			case "bool":
				dstValues.Field(i).SetBool(fValue.BooleanValue)
			case "int32", "int64":
				dstValues.Field(i).SetInt(fValue.IntegerValue)
			case "float32", "float64":
				dstValues.Field(i).SetFloat(fValue.DoubleValue)
			}
		}
	}
}

func toValue(fType string, fValue reflect.Value) (value dbv2.Value, err error) {
	switch fType {
	case "string":
		value.StringValue = fValue.String()
	case "bool":
		value.BooleanValue = fValue.Bool()
	case "int32", "int64":
		value.IntegerValue = fValue.Int()
	case "float32", "float64":
		value.DoubleValue = fValue.Float()
	case "slice":
		//TODO add complex type to the slicell
		if fValue.Type().Elem().Kind() == reflect.Uint8 {
			//BlobValue is a string in the datastore entity proto
			value.BlobValue = string(fValue.Bytes())
		} else {
			size := fValue.Len()
			values := make([]*dbv2.Value, size)
			for i := 0; i < size; i++ {
				val, _ := toValue(fValue.Type().Elem().Kind().String(), fValue.Index(i))
				values[i] = &val
			}
			value.ArrayValue = &dbv2.ArrayValue{
				Values: values,
			}
		}
	case "map":
		err = errors.New("datatype[map] not supported")
	case "ptr":
		err = errors.New("datatype[ptr] not supported")
	default:
		fmt.Println("inside default case")
	}
	return value, err
}
