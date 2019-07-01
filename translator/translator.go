package translator

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	dbv2 "google.golang.org/api/datastore/v1"

	//dbv1 "google.golang.org/appengine/datastore"
	"reflect"
	"strings"

	"errors"

	structpb "github.com/golang/protobuf/ptypes/struct"
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
		fName := srcValues.Type().Field(i).Name
		if !strings.ContainsAny(fName, "XXX_") {
			//fType := srcValues.Field(i).Type().Kind().String()
			value, err := toValue(srcValues.Field(i))
			// fmt.Printf("name:%s, type:%v, value:%v\n",name,fType,value)
			if err == nil {
				properties[fName] = value
			} else {
				fmt.Printf("field: %s, err: %v\n", fName, err)
			}
		}
	}
	entity.Properties = properties
	return entity
}

func DatastoreEntityToProtoMessage(src dbv2.Entity, dst proto.Message) {
	dstValues := reflect.ValueOf(dst).Elem()
	for i := 0; i < dstValues.NumField(); i++ {
		fName := dstValues.Type().Field(i).Name
		if !strings.Contains(fName, "XXX_") {
			fValue := getProperty(src.Properties, fName)
			fType := dstValues.Type().Field(i).Type.Kind()
			fmt.Printf("name: %s, type: %s\n", fName, fType)
			switch fType {
			case reflect.String:
				dstValues.Field(i).SetString(fValue.StringValue)
			case reflect.Bool:
				dstValues.Field(i).SetBool(fValue.BooleanValue)
			case reflect.Int32, reflect.Int64:
				dstValues.Field(i).SetInt(fValue.IntegerValue)
			case reflect.Float32, reflect.Float64:
				dstValues.Field(i).SetFloat(fValue.DoubleValue)
			case reflect.Map:
				entity := fValue.EntityValue
				if entity != nil {
					switch dstValues.Type().Field(i).Type.String() {
					// rudimentary impl
					case "map[string]string":
						m := make(map[string]string)
						for k, v := range entity.Properties {
							m[fmt.Sprint(k)] = v.StringValue
						}
						dstValues.Field(i).Set(reflect.ValueOf(m))
					case "map[string]int32":
						m := make(map[string]int32)
						for k, v := range entity.Properties {
							m[fmt.Sprint(k)] = int32(v.IntegerValue)
						}
						dstValues.Field(i).Set(reflect.ValueOf(m))
					}
				}
			}
		}
	}
}

func toValue(fValue reflect.Value) (value dbv2.Value, err error) {
	switch fValue.Kind() {
	case reflect.String:
		value.StringValue = fValue.String()
	case reflect.Bool:
		value.BooleanValue = fValue.Bool()
	case reflect.Int32, reflect.Int64:
		value.IntegerValue = fValue.Int()
	case reflect.Float32, reflect.Float64:
		value.DoubleValue = fValue.Float()
	case reflect.Slice:
		//TODO add complex type to the slicell
		if fValue.Type().Elem().Kind() == reflect.Uint8 {
			//BlobValue is a string in the datastore entity proto
			value.BlobValue = string(fValue.Bytes())
		} else {
			size := fValue.Len()
			values := make([]*dbv2.Value, size)
			for i := 0; i < size; i++ {
				val, _ := toValue(fValue.Index(i))
				values[i] = &val
			}
			value.ArrayValue = &dbv2.ArrayValue{
				Values: values,
			}
		}
	case reflect.Map:
		mapValues := reflect.ValueOf(fValue.Interface())
		innerEntity := make(map[string]dbv2.Value)
		for _, key := range mapValues.MapKeys() {
			k := fmt.Sprint(key)
			//TODO what if there is an error?
			v, _ := toValue(mapValues.MapIndex(key))
			//fmt.Printf("key; %v, value: %v\n",k,v)
			innerEntity[k] = v
		}
		value.EntityValue = &dbv2.Entity{
			Properties: innerEntity,
		}
	case reflect.Ptr:
		switch fValue.Type().String() {
		case "*structpb.Struct":
			fmt.Println("inside *structpb.Struct")
			fields := fValue.Elem().FieldByName("Fields")
			innerEntity := make(map[string]dbv2.Value)
			for _, value := range fields.MapKeys() {
				v := fields.MapIndex(value).Interface().(*structpb.Value)
				//don't know if there is another way of doing this, trick here is *structpb.Value
				if x, ok := v.GetKind().(*structpb.Value_StringValue); ok {
					innerEntity[fmt.Sprint(value)] = dbv2.Value{StringValue: x.StringValue}
				} else if x, ok := v.GetKind().(*structpb.Value_BoolValue); ok {
					innerEntity[fmt.Sprint(value)] = dbv2.Value{BooleanValue: x.BoolValue}
				} else if x, ok := v.GetKind().(*structpb.Value_NumberValue); ok {
					//structpbStruct on supports float64
					innerEntity[fmt.Sprint(value)] = dbv2.Value{DoubleValue: x.NumberValue}
				} else if _, ok := v.GetKind().(*structpb.Value_ListValue); ok {
					err = errors.New("list is not supported yet")
					// TODO  figure out this
					// innerEntity[fmt.Sprint(key)] = dbv2.Value{ArrayValue: x.ListValue}
				} else if x, ok := v.GetKind().(*structpb.Value_NullValue); ok {
					innerEntity[fmt.Sprint(value)] = dbv2.Value{NullValue: string(x.NullValue)}
				}
			}
			value.EntityValue = &dbv2.Entity{Properties: innerEntity}
		case "*timestamp.Timestamp":
			fmt.Println("inside *timestamp.Timestamp")
		}
		//err = errors.New("datatype[ptr] not supported")
	default:
		fmt.Println("inside default case")
	}
	return value, err
}
