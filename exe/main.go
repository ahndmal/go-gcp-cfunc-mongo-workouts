package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

func main() {
	uri := os.Getenv("DB_URL")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("workouts").Collection("workouts")
	var workout bson.M
	err = coll.FindOne(context.TODO(), bson.D{{"workout_type", "BACK"}}).Decode(&workout)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title %s\n", "title")
		return
	}
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(workout, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)
}
