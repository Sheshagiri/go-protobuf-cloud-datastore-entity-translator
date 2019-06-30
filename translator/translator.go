package translator

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/struct"
	"reflect"
	//"github.com/golang/protobuf/ptypes/timestamp"
	dbv2 "google.golang.org/api/datastore/v1"
	//dbv2 "google.golang.org/appengine/datastore"
	//"github.com/Sheshagiri/protobuf-struct/models"
	"errors"
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
	fmt.Printf("Entity: %v\n",entity)
	fmt.Println("@@@@@@@")
	return entity
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
		//TODO add complex type to the slice
		if fValue.Type().Elem().Kind() == reflect.Uint8 {
			//BlobValue is a string in the datastore entity proto
			value.BlobValue = string(fValue.Bytes())
		} else {
			size := fValue.Len()
			values := make([]*dbv2.Value,size)
			for i := 0; i < size ; i++ {
				val, _ := toValue(fValue.Type().Elem().Kind().String(),fValue.Index(i))
				values[i] = &val
			}
			value.ArrayValue = &dbv2.ArrayValue{
				Values:values,
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
