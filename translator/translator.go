package translator

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/timestamp"
	"regexp"
	"time"
)

func GetProperty(properties []datastore.Property, name string) interface{} {
	for _, property := range properties {
		if property.Name == name {
			return property.Value
		}
	}
	return nil
}

// ProtoMessageToDatastoreEntity will generate an Entity Protobuf that datastore understands
func ProtoMessageToDatastoreEntity(src proto.Message, snakeCase bool) (entity datastore.Entity, err error) {
	srcValues := reflect.ValueOf(src).Elem()
	properties := make([]datastore.Property, 0)

	for i := 0; i < srcValues.NumField(); i++ {
		fName := srcValues.Type().Field(i).Name
		if !strings.ContainsAny(fName, "XXX_") {
			var value interface{}
			if value, err = toValue(srcValues.Field(i)); err != nil {
				return
			} else {
				if snakeCase {
					fName = toSnakeCase(fName)
				}
				properties = append(properties, datastore.Property{
					Name:  fName,
					Value: value,
				})
			}
		}
	}
	entity.Properties = properties
	return
}

// DatastoreEntityToProtoMessage converts any given datastore.Entity to supplied proto.Message
func DatastoreEntityToProtoMessage(src datastore.Entity, dst proto.Message, snakeCase bool) (err error) {
	dstValues := reflect.ValueOf(dst).Elem()
	for i := 0; i < dstValues.NumField(); i++ {
		fName := dstValues.Type().Field(i).Name
		if !strings.Contains(fName, "XXX_") {
			keyName := fName
			if snakeCase {
				keyName = toSnakeCase(fName)
			}
			fValue := GetProperty(src.Properties, keyName)
			fType := dstValues.Type().Field(i).Type.Kind()
			log.Printf("name: %s, type: %s\n", fName, fType)
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
					// TODO use type to dynamically make an array and make use of toValue function
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
				switch dstValues.Type().Field(i).Type.String() {
				// rudimentary impl, as I can't get hold of the type with Kind() here, look at Indirect() later
				case "map[string]string":
					m := make(map[string]string)
					for _, property := range entity.Properties {
						m[property.Name] = property.Value.(string)
					}
					dstValues.Field(i).Set(reflect.ValueOf(m))
				case "map[string]int32":
					m := make(map[string]int32)
					for _, property := range entity.Properties {
						m[property.Name] = int32(property.Value.(int64))
					}
					dstValues.Field(i).Set(reflect.ValueOf(m))
				case "map[string]*structpb.Value":
					errString := "map[string]*structpb.Value is not supported yet"
					log.Println(errString)
					err = errors.New(errString)
				}
			case reflect.Ptr:
				switch dstValues.Type().Field(i).Type.String() {
				case "*timestamp.Timestamp":
					if t, ok := fValue.(time.Time); ok {
						ts, _ := ptypes.TimestampProto(t)
						dstValues.Field(i).Set(reflect.ValueOf(ts))
					}
				case "*structpb.Struct":
					entityValue := fValue.(*datastore.Entity)
					m := make(map[string]*structpb.Value)
					for _, property := range entityValue.Properties {
						v := reflect.ValueOf(property.Value).Kind()
						switch v {
						case reflect.String:
							m[property.Name] = &structpb.Value{
								Kind: &structpb.Value_StringValue{property.Value.(string)},
							}
						case reflect.Bool:
							m[property.Name] = &structpb.Value{
								Kind: &structpb.Value_BoolValue{property.Value.(bool)},
							}
						case reflect.Float64:
							m[property.Name] = &structpb.Value{
								Kind: &structpb.Value_NumberValue{property.Value.(float64)},
							}
						case reflect.Int32:
							m[property.Name] = &structpb.Value{
								Kind: &structpb.Value_NullValue{},
							}
						}
					}
					s := &structpb.Struct{
						Fields: m,
					}
					dstValues.Field(i).Set(reflect.ValueOf(s))
				}
			default:
				errString := fmt.Sprintf("datatype[%s] not supported", fType)
				err = errors.New(errString)
			}
		}
	}
	return
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
		//TODO add complex type to the slice
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
		// I can't get hold of the type with Kind() here, look at Indirect() later
		case "*structpb.Struct":
			//log.Println("inside *structpb.Struct")
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
					} else if _, ok := v.GetKind().(*structpb.Value_ListValue); ok {
						err = errors.New("list inside a google.protobuf.Struct is not supported yet")
						// TODO  figure out this
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
			if !fValue.IsNil() {
				ts := fValue.Interface().(*timestamp.Timestamp)
				value, _ = ptypes.Timestamp(ts)
			}
		default:
			errString := fmt.Sprintf("datatype[%s] not supported", fValue.Type().String())
			log.Println(errString)
			err = errors.New(errString)
		}
	default:
		errString := fmt.Sprintf("datatype[%s] not supported", fValue.Type().String())
		log.Println(errString)
		err = errors.New(errString)
	}
	return value, err
}

func toSnakeCase(name string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(name, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
