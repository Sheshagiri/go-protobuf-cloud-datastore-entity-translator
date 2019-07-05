package translator

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/timestamp"
	pb "google.golang.org/genproto/googleapis/datastore/v1"
	"regexp"
)

// ProtoMessageToDatastoreEntity will generate an Entity Protobuf that datastore understands
func ProtoMessageToDatastoreEntity(src proto.Message, snakeCase bool) (entity pb.Entity, err error) {
	srcValues := reflect.ValueOf(src).Elem()
	properties := make(map[string]*pb.Value)

	for i := 0; i < srcValues.NumField(); i++ {
		fName := srcValues.Type().Field(i).Name
		if !strings.ContainsAny(fName, "XXX_") {
			var value *pb.Value
			if value, err = toValue(srcValues.Field(i)); err != nil {
				return
			} else {
				if snakeCase {
					fName = toSnakeCase(fName)
				}
				properties[fName] = value
			}
		}
	}
	entity.Properties = properties
	return
}

// DatastoreEntityToProtoMessage converts any given datastore.Entity to supplied proto.Message
func DatastoreEntityToProtoMessage(src *pb.Entity, dst proto.Message, snakeCase bool) (err error) {
	dstValues := reflect.ValueOf(dst).Elem()
	for i := 0; i < dstValues.NumField(); i++ {
		fName := dstValues.Type().Field(i).Name
		if !strings.Contains(fName, "XXX_") {
			keyName := fName
			if snakeCase {
				keyName = toSnakeCase(fName)
			}
			fValue := src.Properties[keyName]
			fType := dstValues.Type().Field(i).Type.Kind()
			log.Printf("name: %s, type: %s\n", fName, fType)
			switch fType {
			case reflect.Bool:
				dstValues.Field(i).SetBool(fValue.GetBooleanValue())
			case reflect.String:
				dstValues.Field(i).SetString(fValue.GetStringValue())
			case reflect.Float64, reflect.Float32:
				dstValues.Field(i).SetFloat(fValue.GetDoubleValue())
				//since enums are inherently ints handle then specially
			case reflect.Int64, reflect.Int32:
				dstValues.Field(i).SetInt(fValue.GetIntegerValue())
			case reflect.Slice:
				if dstValues.Type().Field(i).Type.Elem().Kind() == reflect.Uint8 {
					dstValues.Field(i).SetBytes(fValue.GetBlobValue())
				} else {
					// get elements from ArrayValue
					arrayValues := fValue.GetArrayValue().Values
					// TODO use type to dynamically make an array and make use of toValue function
					switch dstValues.Type().Field(i).Type.Elem().Kind() {
					case reflect.String:
						array := make([]string, len(arrayValues))
						for k := 0; k < len(array); k++ {
							v := arrayValues[k]
							array[k] = v.GetStringValue()
						}
						dstValues.Field(i).Set(reflect.ValueOf(array))
					case reflect.Int32:
						array := make([]int32, len(arrayValues))
						for k := 0; k < len(array); k++ {
							v := arrayValues[k]
							array[k] = int32(v.GetIntegerValue())
						}
						dstValues.Field(i).Set(reflect.ValueOf(array))
					}
				}
			case reflect.Map:
				entity := fValue.GetEntityValue()
				switch dstValues.Type().Field(i).Type.String() {
				// rudimentary impl, as I can't get hold of the type with Kind() here, look at Indirect() later
				case "map[string]string":
					m := make(map[string]string)
					for key, value := range entity.Properties {
						m[key] = value.GetStringValue()
					}
					dstValues.Field(i).Set(reflect.ValueOf(m))
				case "map[string]int32":
					m := make(map[string]int32)
					for key, value := range entity.Properties {
						m[key] = int32(value.GetIntegerValue())
					}
					dstValues.Field(i).Set(reflect.ValueOf(m))
				case "map[string]*structpb.Value":
					m := make(map[string]*structpb.Value)
					for key, value := range entity.Properties {
						if x, ok := reflect.ValueOf(value).Interface().(*pb.Value); ok {
							m[key] = &structpb.Value{
								Kind: &structpb.Value_StringValue{
									StringValue: x.GetStringValue(),
								},
							}
							continue
						}
						if x, ok := reflect.ValueOf(value).Interface().(*pb.Value); ok {
							m[key] = &structpb.Value{
								Kind: &structpb.Value_BoolValue{
									BoolValue: x.GetBooleanValue(),
								},
							}
							continue
						}
						if x, ok := reflect.ValueOf(value).Interface().(*pb.Value); ok {
							m[key] = &structpb.Value{
								Kind: &structpb.Value_NumberValue{
									NumberValue: x.GetDoubleValue(),
								},
							}
							continue
						}
						if x, ok := reflect.ValueOf(value).Interface().(*pb.Value); ok {
							m[key] = &structpb.Value{
								Kind: &structpb.Value_NullValue{
									NullValue: x.GetNullValue(),
								},
							}
							continue
						}
						//TODO handle list and struct specially
						/*if x, ok := reflect.ValueOf(value).Interface().(*pb.Value); ok {
							m[key] = &structpb.Value{
								Kind:&structpb.Value_StructValue{
									StructValue:x.GetEntityValue(),
								},
							}
						}
						if x, ok := reflect.ValueOf(value).Interface().(*pb.Value); ok {
							m[key] = &structpb.Value{
								Kind:&structpb.Value_ListValue{
									ListValue:x.GetArrayValue(),
								},
							}
						}*/
					}
					dstValues.Field(i).Set(reflect.ValueOf(m))
				}
			case reflect.Ptr:
				if !reflect.ValueOf(fValue).IsNil() {
					switch dstValues.Type().Field(i).Type.String() {
					case "*timestamp.Timestamp":
						if fValue.GetTimestampValue() != nil {
							dstValues.Field(i).Set(reflect.ValueOf(fValue.GetTimestampValue()))
						}
					case "*structpb.Struct":
						entityValue := fValue.GetEntityValue()
						if entityValue != nil {
							s := &structpb.Struct{}
							m := make(map[string]*structpb.Value)
							for key, value := range entityValue.Properties {
								log.Printf("value type is: %T", value.ValueType)
								if val, ok := reflect.ValueOf(value.ValueType).Interface().(*pb.Value_DoubleValue); ok {
									m[key] = &structpb.Value{
										Kind:&structpb.Value_NumberValue{
											val.DoubleValue,
										},
									}
								} else if val, ok := reflect.ValueOf(value.ValueType).Interface().(*pb.Value_StringValue); ok {
									m[key] = &structpb.Value{
										Kind:&structpb.Value_StringValue{
											val.StringValue,
										},
									}
								} else if val, ok := reflect.ValueOf(value.ValueType).Interface().(*pb.Value_BooleanValue); ok {
									m[key] = &structpb.Value{
										Kind:&structpb.Value_BoolValue{
											val.BooleanValue,
										},
									}
								} else if _, ok := reflect.ValueOf(value.ValueType).Interface().(*pb.Value_NullValue); ok {
									m[key] = &structpb.Value{
										Kind:&structpb.Value_NullValue{},
									}
								} /*else if _, ok := reflect.ValueOf(value.ValueType).Interface().(*pb.Value_ArrayValue); ok {
									m[key] = &structpb.Value{
										Kind:&structpb.Value_ListValue{

										},
									}
								} else if _, ok := reflect.ValueOf(value.ValueType).Interface().(*pb.Value_EntityValue); ok {
									m[key] = &structpb.Value{
										Kind:&structpb.Value_StructValue{

										},
									}
								}*/
							}
							s.Fields = m
							dstValues.Field(i).Set(reflect.ValueOf(s))
						}
					}
				}
			default:
				errString := fmt.Sprintf("datatype[%s] not supported", fType)
				err = errors.New(errString)
			}
		}
	}
	return
}

func toValue(fValue reflect.Value) (*pb.Value, error) {
	value := &pb.Value{}
	var err error
	switch fValue.Kind() {
	case reflect.String:
		value.ValueType = &pb.Value_StringValue{
			StringValue: fValue.String(),
		}
	case reflect.Bool:
		value.ValueType = &pb.Value_BooleanValue{
			BooleanValue: fValue.Bool(),
		}
	case reflect.Int64, reflect.Int32:
		value.ValueType = &pb.Value_IntegerValue{
			IntegerValue: fValue.Int(),
		}
	case reflect.Float32, reflect.Float64:
		value.ValueType = &pb.Value_DoubleValue{
			DoubleValue: fValue.Float(),
		}
	case reflect.Slice:
		//TODO add complex type to the slice
		if fValue.Type().Elem().Kind() == reflect.Uint8 {
			//BlobValue is a string in the datastore entity proto
			value.ValueType = &pb.Value_BlobValue{
				BlobValue: fValue.Bytes(),
			}
		} else {
			size := fValue.Len()
			values := make([]*pb.Value, 0)
			for i := 0; i < size; i++ {
				val, _ := toValue(fValue.Index(i))
				values = append(values, val)
			}
			value.ValueType = &pb.Value_ArrayValue{
				ArrayValue: &pb.ArrayValue{
					Values: values,
				},
			}
		}
	case reflect.Map:
		mapValues := reflect.ValueOf(fValue.Interface())
		entity := &pb.Entity{}
		properties := make(map[string]*pb.Value)
		for _, key := range mapValues.MapKeys() {
			k := fmt.Sprint(key)
			//TODO what if there is an error?
			v, _ := toValue(mapValues.MapIndex(key))
			//fmt.Printf("key; %v, value: %v\n",k,v)
			properties[k] = v
		}
		entity.Properties = properties
		value.ValueType = &pb.Value_EntityValue{
			EntityValue: entity,
		}
	case reflect.Ptr:
		switch fValue.Type().String() {
		// I can't get hold of the type with Kind() here, look at Indirect() later
		case "*structpb.Struct":
			log.Println("inside *structpb.Struct")
			if !fValue.IsNil() {
				fields := fValue.Elem().FieldByName("Fields")
				properties := make(map[string]*pb.Value)
				for _, value := range fields.MapKeys() {
					v := fields.MapIndex(value).Interface().(*structpb.Value)
					//don't know if there is another way of doing this, trick here is *structpb.Value
					if x, ok := v.GetKind().(*structpb.Value_StringValue); ok {
						properties[fmt.Sprint(value)] = &pb.Value{
							ValueType: &pb.Value_StringValue{
								StringValue: x.StringValue,
							},
						}
					} else if x, ok := v.GetKind().(*structpb.Value_BoolValue); ok {
						properties[fmt.Sprint(value)] = &pb.Value{
							ValueType: &pb.Value_BooleanValue{
								BooleanValue: x.BoolValue,
							},
						}
					} else if x, ok := v.GetKind().(*structpb.Value_NumberValue); ok {
						//structpbStruct on supports float64
						properties[fmt.Sprint(value)] = &pb.Value{
							ValueType: &pb.Value_DoubleValue{
								DoubleValue: x.NumberValue,
							},
						}
					} else if _, ok := v.GetKind().(*structpb.Value_ListValue); ok {
						err = errors.New("list inside a google.protobuf.Struct is not supported yet")
						// TODO  figure out this
					} else if x, ok := v.GetKind().(*structpb.Value_NullValue); ok {
						properties[fmt.Sprint(value)] = &pb.Value{
							ValueType: &pb.Value_NullValue{
								NullValue: x.NullValue,
							},
						}
					}
				}
				value.ValueType = &pb.Value_EntityValue{
					EntityValue: &pb.Entity{
						Properties: properties,
					},
				}
			}
		case "*timestamp.Timestamp":
			if !fValue.IsNil() {
				ts := fValue.Interface().(*timestamp.Timestamp)
				value.ValueType = &pb.Value_TimestampValue{
					TimestampValue: ts,
				}
			}
		case "*structpb.Value":
		if !fValue.IsNil() {
			v := fValue.Interface().(*structpb.Value)
			//don't know if there is another way of doing this, trick here is *structpb.Value
			if x, ok := v.GetKind().(*structpb.Value_StringValue); ok {
				value.ValueType = &pb.Value_StringValue{
					StringValue: x.StringValue,
				}
			} else if x, ok := v.GetKind().(*structpb.Value_BoolValue); ok {
				value.ValueType = &pb.Value_BooleanValue{
					BooleanValue: x.BoolValue,
				}
			} else if x, ok := v.GetKind().(*structpb.Value_NumberValue); ok {
				//structpbStruct on supports float64
				value.ValueType = &pb.Value_DoubleValue{
					DoubleValue: x.NumberValue,
				}
			} else if _, ok := v.GetKind().(*structpb.Value_ListValue); ok {
				err = errors.New("list inside a google.protobuf.Struct is not supported yet")
				// TODO  figure out this
			} else if x, ok := v.GetKind().(*structpb.Value_NullValue); ok {
				value.ValueType = &pb.Value_NullValue{
					NullValue: x.NullValue,
				}
			}
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
