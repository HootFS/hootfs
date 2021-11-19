package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hootfs/hootfs/src/core/metastore"
	"go.mongodb.org/mongo-driver/bson"
)

type Person struct {
	Name     string
	Siblings []string
}

func main() {
	ms, err := metastore.NewMetaStore()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := ms.Disconnect(); err != nil {
			log.Fatal(err)
		}
	}()

	collection := ms.TestCollection("people")

	// res, err := collection.InsertOne(context.TODO(),
	// 	bson.M{"Name": "Sal", "Siblings": bson.A{"Joel", "Sam"}},
	// )

	// if err != nil {
	// 	log.Fatal(res)
	// }

	var result Person

	filter := bson.D{
		bson.E{Key: "Name", Value: "Sal"},
	}

	err = collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		log.Fatal(err)
	}

	// This worked BOIIIII !!!!!!!
	fmt.Println("Found : ", result.Name, " ", result.Siblings)

	// bob := Person{
	// 	Name: "Bob",
	// 	Age:  22,
	// }

	// fmt.Println("Result : ", res)
}
