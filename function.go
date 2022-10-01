package p

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
)

func GetWorkouts(writer http.ResponseWriter, req *http.Request) {
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
	//workouts := make([]bson.M, 0)
	var workouts []bson.D
	//var workouts []Workout
	params := req.URL.Query()
	if len(params.Get("type")) == 0 { // get by TYPE if requested
		cursor, err := coll.Find(context.TODO(), bson.D{{"type", params.Get("type")}})
		if err == mongo.ErrNoDocuments {
			fmt.Printf("No document was found with the type %s\n", "type")
			return
		}
		err2 := cursor.All(context.TODO(), &workouts)
		if err2 != nil {
			log.Panicln(err2)
		}
		jsonData, err := json.MarshalIndent(workouts, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Fprint(writer, string(jsonData))
	}
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title %s\n", "title")
		return
	}
	err2 := cursor.All(context.TODO(), &workouts)
	if err2 != nil {
		log.Panicln(err2)
	}
	log.Printf("Type of workouts is %T", workouts) // reflect.TypeOf(workouts)
	log.Println(workouts[0])

	jsonData, err3 := json.Marshal(workouts)
	if err3 != nil {
		log.Panicln(err)
	}
	req.Header.Add("Content-Type", "application/json")

	fmt.Fprint(writer, string(jsonData))
}
