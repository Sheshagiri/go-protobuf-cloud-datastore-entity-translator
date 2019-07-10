package playground

import (
	"fmt"
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/unsupported"
	"google.golang.org/genproto/googleapis/datastore/v1"
	"testing"
	"reflect"
	"strings"
	clientSDK "cloud.google.com/go/datastore"
)

func Test1(t *testing.T) {
	input := &unsupported.Child{
		Name: "child-1",
		Parent: &unsupported.Parent{
			Name: "parent-1",
		},
	}
	fmt.Println("Input Proto: ", input)
	/*inputEntity, err := translator.ProtoMessageToDatastoreEntity(input, false)
	fmt.Println("Error: ", err)*/
	inputEntity := &datastore.Entity{
		Properties: map[string]*datastore.Value{
			"Name": {ValueType: &datastore.Value_StringValue{StringValue: "child-1"}},
			"Parent": {ValueType: &datastore.Value_EntityValue{
				EntityValue: &datastore.Entity{
					Properties: map[string]*datastore.Value{
						"name": {ValueType: &datastore.Value_StringValue{StringValue: "parent-1"}},
					},
				},
			}},
		},
	}

	fmt.Println("Input Entity: ", inputEntity)
	output := unsupported.Child{}
	/*fmt.Println("------------------------")
	detopm(inputEntity, output)
	fmt.Println("------------------------")
	fmt.Println("Output Proto: ", output)*/
	translate(output)
	fmt.Println("output: ",output.GetParent())
}

func detopm(entity *datastore.Entity, dst interface{}) {
	fmt.Println("Output kind: ",reflect.ValueOf(dst).Kind())
	dstValues := reflect.ValueOf(dst).Elem()
	//printAll(dstValues)
	for i:=0;i<dstValues.NumField();i++{
		fName := dstValues.Type().Field(i).Name
		fType := dstValues.Type().Field(i).Type
		fValue := entity.Properties[fName]
		fmt.Println("Value is: ",fValue)
		if ! strings.ContainsAny(fName, "XXX_"){
			fmt.Println(fName," : ",fType)
			switch fType.Kind() {
			case reflect.Ptr:
				iv := dstValues.Field(i).Interface()
				fmt.Println("interface: ",iv)
			}
			/*switch fType.Kind() {
			case reflect.String:
				dstValues.Field(i).SetString(fValue.GetStringValue())
			case reflect.Ptr:
				p := reflect.Zero(fType).Elem()
				v := reflect.ValueOf(p)
				if v.Kind() == reflect.Struct{
					fmt.Println(v)
					val := reflect.ValueOf(v.Interface())
					fmt.Println(reflect.Indirect(val).IsNil())
					iv := reflect.ValueOf(p).Interface()
					fmt.Println(reflect.ValueOf(iv).FieldByName("Name"))
					fmt.Println("inside a struct")
					fmt.Println("number of fields: ",v.NumField())
					vv := v.FieldByName("Name")
					if vv.IsValid() {
						fmt.Println("is valid")
					}
				}
			}*/
		}
	}
}


func printAll(v reflect.Value) {
	s := v
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())

		if f.Kind().String() == "struct" {
			x1 := reflect.ValueOf(f.Interface())
			fmt.Printf("type2: %s\n", x1)
			printAll(x1)
		}
	}
}

func translate(obj interface{}) interface{} {
	// Wrap the original in a reflect.Value
	original := reflect.ValueOf(obj)

	copy := reflect.New(original.Type()).Elem()
	translateRecursive(copy, original)

	// Remove the reflection wrapper
	return copy.Interface()
}

func translateRecursive(copy, original reflect.Value) {
	switch original.Kind() {
	// The first cases handle nested structures and translate them recursively

	// If it is a pointer we need to unwrap and call once again
	case reflect.Ptr:
		// To get the actual value of the original we have to call Elem()
		// At the same time this unwraps the pointer so we don't end up in
		// an infinite recursion
		originalValue := original.Elem()
		// Check if the pointer is nil
		if !originalValue.IsValid() {
			return
		}
		// Allocate a new object and set the pointer to it
		copy.Set(reflect.New(originalValue.Type()))
		// Unwrap the newly created pointer
		translateRecursive(copy.Elem(), originalValue)

		// If it is an interface (which is very similar to a pointer), do basically the
		// same as for the pointer. Though a pointer is not the same as an interface so
		// note that we have to call Elem() after creating a new object because otherwise
		// we would end up with an actual pointer
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()
		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()
		translateRecursive(copyValue, originalValue)
		copy.Set(copyValue)

		// If it is a struct we translate each field
	case reflect.Struct:
		for i := 0; i < original.NumField(); i += 1 {
			fmt.Println("Original --> Name: ",original.Field(i).Type().Name())
			fmt.Println("Copy --> Name: ",copy.Field(i).Type().Name())
			translateRecursive(copy.Field(i), original.Field(i))
		}

		// If it is a slice we create a new slice and translate each element
	case reflect.Slice:
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i += 1 {
			translateRecursive(copy.Index(i), original.Index(i))
		}

		// If it is a map we create a new map and translate each value
	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			// New gives us a pointer, but again we want the value
			copyValue := reflect.New(originalValue.Type()).Elem()
			translateRecursive(copyValue, originalValue)
			copy.SetMapIndex(key, copyValue)
		}

		// Otherwise we cannot traverse anywhere so this finishes the the recursion

		// If it is a string translate it (yay finally we're doing what we came for)
	case reflect.String:
		translatedString := dict[original.Interface().(string)]
		copy.SetString(translatedString)

		// And everything else will simply be taken from the original
	default:
		copy.Set(original)
	}

}

var dict = map[string]string{
	"Hello!":                 "Hallo!",
	"What's up?":             "Was geht?",
	"translate this":         "übersetze dies",
	"point here":             "zeige hier her",
	"translate this as well": "übersetze dies auch...",
	"and one more":           "und noch eins",
	"deep":                   "tief",
}

func TestCheckNestedStruct(t *testing.T) {
	child := unsupported.Child{
		Parent:&unsupported.Parent{
			Name:"this is a parent.",
		},
	}
	v := reflect.ValueOf(&child).Elem()
	parent := v.FieldByName("Parent")
	if parent.Kind()== reflect.Ptr {
		fmt.Println("parent is a pointer.")
		parent := parent.Elem()
		newParent := reflect.New(parent.Type()).Interface()
		fmt.Println(newParent)
	}
}

func TestEmptyMessage(t *testing.T){
	dst := unsupported.Child{}
	fmt.Println(dst)
	clientSDK.EntityToStruct(dst,nil)
	fmt.Println(dst)
}