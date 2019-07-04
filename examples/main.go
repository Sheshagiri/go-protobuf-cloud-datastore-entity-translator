package main

import (
	"context"

	"fmt"
	"time"
	"cloud.google.com/go/datastore"
	//pb "google.golang.org/genproto/googleapis/datastore/v1"
	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/models/example"

	"github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/translator"
)

func main() {
	type Employee struct {
		FirstName          string
		LastName           string
		HireDate           time.Time
		AttendedHRTraining bool
	}

	ctx := context.Background()
	simpleProto := &example.ExampleNestedModel{
		StringKey:"validate",
		Int32Key:204,
	}
	employee := &Employee{
		FirstName: "Antonio",
		LastName:  "Salieri",
		HireDate:  time.Now(),
	}
	employee.AttendedHRTraining = true
	fmt.Println(simpleProto)
	fmt.Println(employee)
	client,_ := datastore.NewClient(ctx,"st2-saas-prototype-dev")
	k1 := datastore.NameKey("Employee", employee.FirstName,nil)
	k2 := datastore.NameKey("SimpleProto", simpleProto.GetStringKey(),nil)
	_, err := client.Put(ctx,k1,employee)
	fmt.Println(err)
	_, err = client.Put(ctx,k2,simpleProto)
	k4 := datastore.NameKey("GoTranslator", "entity",nil)
	translatedEntity,err := translator.ProtoMessageToDatastoreEntity(simpleProto,false)
	_, err = client.Put(ctx,k4,&translatedEntity)
	fmt.Println(err)
}
