package translator

import (
	"testing"
	"reflect"
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/example"
	"fmt"
	"strings"
)

func TestProtoToEntity(t *testing.T) {
	src := &example.ExampleDBModel{
		ComplexArrayKey: []*example.ExampleNestedModel{
			{StringKey:"string-1"},
			{StringKey:"string-2"},
		},
	}
	fmt.Println(src)
	dst := &example.ExampleDBModel{}
	dstValues := reflect.ValueOf(dst).Elem()
	for i := 0; i < dstValues.NumField(); i++ {
		fName := dstValues.Type().Field(i).Name
		fType := dstValues.Type().Field(i).Type.Kind()
		//nm1 := &example.ExampleNestedModel{StringKey:"string-1"}
		//nm2 := &example.ExampleNestedModel{StringKey:"string-2"}
		if !strings.ContainsAny(fName, "XXX_"){
			switch fType {
			case reflect.Slice:
				v := dstValues.Type().Field(i).Type.Elem()
				switch v.Kind() {
				case reflect.String:
					fmt.Println("String: ",fName)
				case reflect.Int32:
					fmt.Println("Int32: ",fName)
				case reflect.Ptr:
					fmt.Println("Pointer: ", fName)
					sliceTyp := reflect.TypeOf(dstValues.Field(i).Type().Name())
					fmt.Println("slice type: ", sliceTyp)
					slice := reflect.MakeSlice(sliceTyp, 0, 0).Interface()
					fmt.Printf("%#v", slice)
				default:
					fmt.Println("Default: ",fName)
				}
			}
		}
	}
}
