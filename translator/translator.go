package translator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/golang/protobuf/proto"
	structpb "github.com/golang/protobuf/ptypes/struct"
	//dbv2 "google.golang.org/api/datastore/v1"
)

func validate(entity datastore.Entity) {
	fmt.Println(entity)
}

func GetProperty(properties []datastore.Property, name string) interface{} {
	for _, property := range properties {
		if property.Name == name {
			return property.Value
		}
	}
	return nil
}

// ProtoMessageToDatastoreEntity will generate an Entity Protobuf that datastore understands
func ProtoMessageToDatastoreEntity(src proto.Message) datastore.Entity {
	srcValues := reflect.ValueOf(src).Elem()
	entity := datastore.Entity{}
	properties := make([]datastore.Property, 0)

	for i := 0; i < srcValues.NumField(); i++ {
		fName := srcValues.Type().Field(i).Name
		if !strings.ContainsAny(fName, "XXX_") {
			//fType := srcValues.Field(i).Type().Kind().String()
			value, err := toValue(srcValues.Field(i))
			// fmt.Printf("name:%s, type:%v, value:%v\n",name,fType,value)
			if err == nil {
				properties = append(properties, datastore.Property{
					Name:  fName,
					Value: value,
				})
			} else {
				fmt.Printf("field: %s, err: %v\n", fName, err)
			}
		}
	}
	entity.Properties = properties
	return entity
}

func DEtoPM(src datastore.Entity, dst proto.Message) {
	dstValues := reflect.ValueOf(dst).Elem()
	for i := 0; i < dstValues.NumField(); i++ {
		fName := dstValues.Type().Field(i).Name
		if !strings.Contains(fName, "XXX_") {
			fValue := GetProperty(src.Properties, fName)
			fType := dstValues.Type().Field(i).Type.Kind()
			fmt.Printf("name: %s, type: %s\n", fName, fType)
			switch fType {
			case reflect.Float64, reflect.Float32, reflect.Bool, reflect.String:
				dstValues.Field(i).Set(reflect.ValueOf(fValue))
			//since enums are inherently ints handle then specially
			case reflect.Int64, reflect.Int32:
				dstValues.Field(i).SetInt(fValue.(int64))
			case reflect.Slice:
				if dstValues.Type().Field(i).Type.Elem().Kind() == reflect.Uint8 {
					dstValues.Field(i).SetBytes([]byte(fValue.(string)))
				} else {
					// get elements from ArrayValue
					arrayValues := fValue.([]interface{})
					fmt.Println(arrayValues)
					switch dstValues.Type().Field(i).Type.Elem().Kind() {
					case reflect.String:
						array := make([]string, len(arrayValues))
						for k := 0; k < len(array); k++ {
							v := arrayValues[k]
							array[k] = v.(string)
						}
						dstValues.Field(i).Set(reflect.ValueOf(array))
					case reflect.Int32:
						array := make([]int32, len(arrayValues))
						for k := 0; k < len(array); k++ {
							v := arrayValues[k]
							array[k] = int32(v.(int64))
						}
						dstValues.Field(i).Set(reflect.ValueOf(array))
					}
				}
			case reflect.Map:
				entity := fValue.(*datastore.Entity)
				fmt.Println(entity.Properties)
				/*switch entityValue.Kind() {
				// for now only entity is present inside a map
				case reflect.TypeOf(&datastore.Entity{}).Kind():
					if !entityValue.IsNil() {
						properties := entityValue.Elem().FieldByName("Properties")
						if !properties.IsNil() {
							switch dstValues.Type().Field(i).Type.String() {
							// rudimentary impl
							case "map[string]string":
								m := make(map[string]string)
								for _, key := range properties.MapKeys() {
									v := properties.MapIndex(key)
									//fmt.Printf("key: %v, value: %v\n",key,v.FieldByName("StringValue"))
									m[key.String()] = v.FieldByName("StringValue").String()
								}
								dstValues.Field(i).Set(reflect.ValueOf(m))
							case "map[string]int32":
								m := make(map[string]int32)
								for _, key := range properties.MapKeys() {
									v := properties.MapIndex(key)
									//fmt.Printf("key: %v, value: %v\n", key, v.FieldByName("IntegerValue"))
									m[key.String()] = int32(v.FieldByName("IntegerValue").Int())
								}
								dstValues.Field(i).Set(reflect.ValueOf(m))
							}
						}
					}
				}*/
			case reflect.Ptr:
				fmt.Println("validate in struct type")
				switch dstValues.Type().Field(i).Type.String() {
				case "*structpb.Struct":
					fmt.Println("inside *structpb.Struct")
					//entityValue := fValue.(*datastore.Entity)
					//fmt.Println(entityValue)
					/*entityValue := reflect.ValueOf(fValue).FieldByName("EntityValue")
					switch entityValue.Kind() {
					// for now only entity is present inside a map
					case reflect.TypeOf(&datastore.Entity{}).Kind():
						if !entityValue.IsNil() {
							properties := entityValue.Elem().FieldByName("Properties")
							fmt.Println(properties)
							//m := make(map [string]*structpb.Value)

							for _, key := range properties.MapKeys(){
								v := reflect.ValueOf(properties.MapIndex(key))
								switch v.Kind() {
								case reflect.Struct:
									t := v.Type()
									size := v.NumField()
									fmt.Println("----------------")
									for i := 0; i < size; i++ {
										name := t.Field(i).Name
										fmt.Println(name)
									}
									fmt.Printf("type: %v, size: %v\n",t, size)
									fmt.Println("----------------")
								}
							}
						}
					}*/
				}
			default:
				fmt.Println("doesn't support yet")
			}
		}
	}
}

func toValue(fValue reflect.Value) (value interface{}, err error) {
	switch fValue.Kind() {
	case reflect.String:
		value = fValue.String()
	case reflect.Bool:
		value = fValue.Bool()
	case reflect.Int64, reflect.Int32:
		value = fValue.Int()
	case reflect.Float32:
		value = float32(fValue.Float())
	case reflect.Float64:
		value = fValue.Float()
	case reflect.Slice:
		//TODO add complex type to the slicell
		if fValue.Type().Elem().Kind() == reflect.Uint8 {
			//BlobValue is a string in the datastore entity proto
			value = string(fValue.Bytes())
		} else {
			size := fValue.Len()
			values := make([]interface{}, size)
			for i := 0; i < size; i++ {
				val, _ := toValue(fValue.Index(i))
				values[i] = val
			}
			value = values
		}
	case reflect.Map:
		mapValues := reflect.ValueOf(fValue.Interface())
		innerEntity := make([]datastore.Property, 0)
		for _, key := range mapValues.MapKeys() {
			k := fmt.Sprint(key)
			//TODO what if there is an error?
			v, _ := toValue(mapValues.MapIndex(key))
			//fmt.Printf("key; %v, value: %v\n",k,v)
			innerEntity = append(innerEntity, datastore.Property{
				Name:  k,
				Value: v,
			})
		}
		value = &datastore.Entity{
			Properties: innerEntity,
		}
	case reflect.Ptr:
		switch fValue.Type().String() {
		case "*structpb.Struct":
			fmt.Println("inside *structpb.Struct")
			if !fValue.IsNil() {
				fields := fValue.Elem().FieldByName("Fields")
				innerEntity := make([]datastore.Property, 0)
				for _, value := range fields.MapKeys() {
					v := fields.MapIndex(value).Interface().(*structpb.Value)
					//don't know if there is another way of doing this, trick here is *structpb.Value
					if x, ok := v.GetKind().(*structpb.Value_StringValue); ok {
						innerEntity = append(innerEntity, datastore.Property{
							Name:  fmt.Sprint(value),
							Value: x.StringValue,
						})
					} else if x, ok := v.GetKind().(*structpb.Value_BoolValue); ok {
						innerEntity = append(innerEntity, datastore.Property{
							Name:  fmt.Sprint(value),
							Value: x.BoolValue,
						})
					} else if x, ok := v.GetKind().(*structpb.Value_NumberValue); ok {
						//structpbStruct on supports float64
						innerEntity = append(innerEntity, datastore.Property{
							Name:  fmt.Sprint(value),
							Value: x.NumberValue,
						})
					} else if x, ok := v.GetKind().(*structpb.Value_ListValue); ok {
						err = errors.New("list is not supported yet")
						// TODO  figure out this
						innerEntity = append(innerEntity, datastore.Property{
							Name:  fmt.Sprint(value),
							Value: x.ListValue,
						})
					} else if x, ok := v.GetKind().(*structpb.Value_NullValue); ok {
						innerEntity = append(innerEntity, datastore.Property{
							Name:  fmt.Sprint(value),
							Value: x.NullValue,
						})
					}
				}
				value = &datastore.Entity{Properties: innerEntity}
			}
		case "*timestamp.Timestamp":
			fmt.Println("inside *timestamp.Timestamp")
			err = errors.New("datatype[ptr] not supported")
		case "*datastore.Entity":
			fmt.Println("inside *datastore.Entity")
			err = errors.New("datatype[ptr] not supported")
		}
		//err = errors.New("datatype[ptr] not supported")
	default:
		fmt.Println("inside default case")
	}
	return value, err
}

/*
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
			case reflect.Ptr:
				fmt.Println("inside pointer")
				switch dstValues.Type().Field(i).Type.String() {
				case "*structpb.Struct":
					//protobufStruct := structpb.Struct{}
					//fields := make(map[string]*structpb.Value)
					fmt.Println(fValue.EntityValue)
					for k, v := range fValue.EntityValue.Properties {
						fmt.Println(k)
						//figure out a way to handle this
						fmt.Println(v)
					}
				}
			}
		}
	}
}
*/
