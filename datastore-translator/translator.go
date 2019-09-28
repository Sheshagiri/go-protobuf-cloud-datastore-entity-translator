package datastore_translator

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"regexp"

	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/proto"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/genproto/googleapis/datastore/v1"

	clientSDK "cloud.google.com/go/datastore"
)

// ProtoMessageToDatastoreEntity will generate an Entity Protobuf that datastore understands
func ProtoMessageToDatastoreEntity(src proto.Message, snakeCase bool, excludeFromIndexName ...string) (entity datastore.Entity, err error) {
	if len(excludeFromIndexName) > 1 {
		return entity, errors.New("exclude only works with one extension")
	}
	srcValues := reflect.ValueOf(src).Elem()
	properties := make(map[string]*datastore.Value)
	// see https://github.com/golang/protobuf/issues/372
	fd, md := descriptor.ForMessage(src.(descriptor.Message))
	// this is the default name
	excludeFromIndexExt := "[exclude_from_index]:true "
	var excludeIndex string
	// use the Extension name is passed else derive it dynamically
	if excludeFromIndexName != nil && len(excludeFromIndexName) > 0 {
		excludeFromIndexExt = fmt.Sprintf("[%s]:true ", excludeFromIndexName[0])
		excludeIndex = excludeFromIndexName[0]
	} else {
		if len(fd.GetExtension()) > 0 {
			excludeIndex = fd.GetExtension()[0].GetName()
			excludeFromIndexExt = fmt.Sprintf("[%s]:true ", excludeIndex)
		}
	}

	excludeFields := make(map[string]string, 0)
	for _, fd := range md.GetField() {
		if fd.GetOptions() != nil {
			excludeFields[fd.GetName()] = fd.GetOptions().String()
		}
	}
	for i := 0; i < srcValues.NumField(); i++ {
		fName := srcValues.Type().Field(i).Name
		if !strings.ContainsAny(fName, "XXX_") {
			var value *datastore.Value
			if value, err = toDatastoreValue(fName, srcValues.Field(i), snakeCase, excludeIndex); err != nil {
				return
			} else {
				if value != nil {
					if snakeCase {
						fName = toSnakeCase(fName)
					}
					if excludeFields[fName] == excludeFromIndexExt {
						value.ExcludeFromIndexes = true
					}
					properties[fName] = value
				}
			}
		}
	}
	entity.Properties = properties
	return
}

// DatastoreEntityToProtoMessage converts any given datastore.Entity to supplied proto.Message
func DatastoreEntityToProtoMessage(src *datastore.Entity, dst proto.Message, snakeCase bool) (err error) {
	entity, err := clientSDK.ProtoToEntity(src, snakeCase)
	if err != nil {
		return err
	}

	err = clientSDK.EntityToStruct(dst, entity)
	if err != nil {
		if strings.ContainsAny(err.Error(), "no such struct field") || strings.ContainsAny(err.Error(), "versus map[string]") {
			err = nil
			// handle google.protobuf.Struct type here
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
					// log.Printf("name: %s, type: %s\n", fName, fType)
					switch fType {
					case reflect.Map:
						if !reflect.ValueOf(fValue).IsNil() {
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
									m[key] = fromDatastoreValueToStructValue(value)
								}
								dstValues.Field(i).Set(reflect.ValueOf(m))
							}
						}
					case reflect.Ptr:
						iv := dstValues.Field(i).Interface()
						switch iv.(type) {
						case *structpb.Struct:
							if !reflect.ValueOf(fValue).IsNil() {
								switch v := reflect.ValueOf(fValue.ValueType).Interface().(type) {
								case *datastore.Value_EntityValue:
									properties := v.EntityValue.Properties
									if properties != nil {
										s := &structpb.Struct{}
										m := make(map[string]*structpb.Value)
										for key, value := range properties {
											// log.Printf("value type is: %T", value.ValueType)
											m[key] = fromDatastoreValueToStructValue(value)
										}
										s.Fields = m
										dstValues.Field(i).Set(reflect.ValueOf(s))
									}
								}
							}
						// all the other pointers like referenced protobuf's
						default:
							if !reflect.ValueOf(fValue).IsNil() {
								switch v := reflect.ValueOf(fValue.ValueType).Interface().(type) {
								case *datastore.Value_EntityValue:
									innerModel, ok := dstValues.Field(i).Interface().(proto.Message)
									if !ok {
										return errors.New("failed to translate: %s" + fName)
									}
									err = DatastoreEntityToProtoMessage(v.EntityValue, innerModel, snakeCase)
									if err != nil {
										return nil
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return err
}

func toDatastoreValue(fName string, fValue reflect.Value, snakeCase bool, excludeFromIndexName string) (*datastore.Value, error) {
	value := &datastore.Value{}
	var err error
	switch fValue.Kind() {
	case reflect.String:
		value.ValueType = &datastore.Value_StringValue{
			StringValue: fValue.String(),
		}
	case reflect.Bool:
		value.ValueType = &datastore.Value_BooleanValue{
			BooleanValue: fValue.Bool(),
		}
	case reflect.Int64, reflect.Int32:
		value.ValueType = &datastore.Value_IntegerValue{
			IntegerValue: fValue.Int(),
		}
	case reflect.Float32, reflect.Float64:
		value.ValueType = &datastore.Value_DoubleValue{
			DoubleValue: fValue.Float(),
		}
	case reflect.Slice:
		// TODO add complex type to the slice
		if fValue.Type().Elem().Kind() == reflect.Uint8 {
			// BlobValue is a string in the datastore entity proto
			value.ValueType = &datastore.Value_BlobValue{
				BlobValue: fValue.Bytes(),
			}
		} else {
			size := fValue.Len()
			values := make([]*datastore.Value, 0)
			for i := 0; i < size; i++ {
				val, err := toDatastoreValue(fName, fValue.Index(i), snakeCase, excludeFromIndexName)
				if err != nil {
					return nil, err
				}
				values = append(values, val)
			}
			value.ValueType = &datastore.Value_ArrayValue{
				ArrayValue: &datastore.ArrayValue{
					Values: values,
				},
			}
		}
	case reflect.Map:
		mapValues := reflect.ValueOf(fValue.Interface())
		entity := &datastore.Entity{}
		properties := make(map[string]*datastore.Value)
		for _, key := range mapValues.MapKeys() {
			k := fmt.Sprint(key)
			// TODO what if there is an error?
			v, _ := toDatastoreValue(fName, mapValues.MapIndex(key), snakeCase, excludeFromIndexName)
			// fmt.Printf("key; %v, value: %v\n",k,v)
			properties[k] = v
		}
		entity.Properties = properties
		value.ValueType = &datastore.Value_EntityValue{
			EntityValue: entity,
		}
	case reflect.Ptr:
		if fValue.IsNil() || !fValue.IsValid() {
			// don't return an error because we still want to retain the proto3 behaviour of having default values
			return nil, nil
		}
		iv := fValue.Interface()
		switch v := iv.(type) {
		case *structpb.Struct:
			properties := make(map[string]*datastore.Value)
			for key, value := range v.Fields {
				properties[key] = fromStructValueToDatastoreValue(value)
			}
			value.ValueType = &datastore.Value_EntityValue{
				EntityValue: &datastore.Entity{
					Properties: properties,
				},
			}
		case *timestamp.Timestamp:
			ts := fValue.Interface().(*timestamp.Timestamp)
			value.ValueType = &datastore.Value_TimestampValue{
				TimestampValue: ts,
			}
		case *structpb.Value:
			value = fromStructValueToDatastoreValue(v)
		default:
			// translate any imported protobuf's that we defined ourself
			if !fValue.IsNil() && fValue.IsValid() {
				if importedProto, ok := reflect.ValueOf(fValue.Interface()).Interface().(proto.Message); ok {
					entityOfImportedProto, err := ProtoMessageToDatastoreEntity(importedProto, snakeCase, excludeFromIndexName)
					if err != nil {
						return nil, err
					}
					value.ValueType = &datastore.Value_EntityValue{
						EntityValue: &datastore.Entity{
							Properties: entityOfImportedProto.Properties,
						},
					}
				}
			}
		}
	default:
		errString := fmt.Sprintf("[toDatastoreValue]: datatype[%s] not supported", fValue.Type().String())
		log.Println(errString)
		err = errors.New(errString)
	}

	return value, err
}

func fromStructValueToDatastoreValue(v *structpb.Value) *datastore.Value {
	pbValue := &datastore.Value{}
	switch v.GetKind().(type) {
	case *structpb.Value_StringValue:
		pbValue.ValueType = &datastore.Value_StringValue{
			StringValue: v.GetStringValue(),
		}
	case *structpb.Value_BoolValue:
		pbValue.ValueType = &datastore.Value_BooleanValue{
			BooleanValue: v.GetBoolValue(),
		}
	case *structpb.Value_NumberValue:
		pbValue.ValueType = &datastore.Value_DoubleValue{
			DoubleValue: v.GetNumberValue(),
		}
	case *structpb.Value_NullValue:
		pbValue.ValueType = &datastore.Value_NullValue{
			NullValue: v.GetNullValue(),
		}
	case *structpb.Value_ListValue:
		values := make([]*datastore.Value, 0)
		for _, val := range v.GetListValue().Values {
			values = append(values, fromStructValueToDatastoreValue(val))
		}
		pbValue.ValueType = &datastore.Value_ArrayValue{
			ArrayValue: &datastore.ArrayValue{
				Values: values,
			},
		}
	case *structpb.Value_StructValue:
		structValue := v.GetStructValue()
		properties := make(map[string]*datastore.Value, 0)
		for key, value := range structValue.GetFields() {
			properties[key] = fromStructValueToDatastoreValue(value)
		}
		pbValue.ValueType = &datastore.Value_EntityValue{
			EntityValue: &datastore.Entity{
				Properties: properties,
			},
		}
	}
	return pbValue
}

func fromDatastoreValueToStructValue(v *datastore.Value) *structpb.Value {
	pbValue := &structpb.Value{}
	iv := reflect.ValueOf(v.ValueType).Interface()
	switch v := iv.(type) {
	case *datastore.Value_BooleanValue:
		pbValue.Kind = &structpb.Value_BoolValue{BoolValue: v.BooleanValue}
	case *datastore.Value_StringValue:
		pbValue.Kind = &structpb.Value_StringValue{StringValue: v.StringValue}
	case *datastore.Value_DoubleValue:
		pbValue.Kind = &structpb.Value_NumberValue{NumberValue: v.DoubleValue}
	case *datastore.Value_NullValue:
		pbValue.Kind = &structpb.Value_NullValue{}
	case *datastore.Value_EntityValue:
		entityValue := v.EntityValue
		fields := make(map[string]*structpb.Value)
		for key, value := range entityValue.GetProperties() {
			fields[key] = fromDatastoreValueToStructValue(value)
		}
		pbValue.Kind = &structpb.Value_StructValue{
			StructValue: &structpb.Struct{
				Fields: fields,
			},
		}
	case *datastore.Value_ArrayValue:
		values := make([]*structpb.Value, 0)
		for _, val := range v.ArrayValue.Values {
			values = append(values, fromDatastoreValueToStructValue(val))
		}
		pbValue.Kind = &structpb.Value_ListValue{
			ListValue: &structpb.ListValue{
				Values: values,
			},
		}
	}
	return pbValue
}

func toSnakeCase(name string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(name, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
