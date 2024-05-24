package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// run this first to ENSURE THAT WE HAVE A SOLID COLLECTION!
type Machine_test struct {
	ID       int `bson:"machine_id"`
	Acc_data int `bson:"acc_data"`
	Vib_data int `bson:"vib_data"`
	Oth_data int `bson:"oth_data"`
}

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	//create the machine_test database with index title from struct
	machineDatabase := client.Database("machine_test")
	machineCollection := machineDatabase.Collection("Machine")

	// .Collection("Machine")

	machine1 := Machine_test{ID: 1, Acc_data: 6, Vib_data: 3, Oth_data: 4}

	//TO ENSURE THE COLLECTION IS HERE!
	_, err = machineCollection.InsertOne(context.TODO(), machine1)
	if err != nil {
		fmt.Println("insert fail")
		log.Fatal(err)
	} else {
		fmt.Println("insert OK!")
	}

	// _, err = machineCollection.DeleteMany(context.TODO(), machine1)

	// if err != nil {
	// 	fmt.Println("delete fail")
	// 	log.Fatal(err)
	// } else {
	// 	fmt.Println("delete OK")
	// }

}
