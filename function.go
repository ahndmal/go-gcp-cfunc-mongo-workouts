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
	queryParams := req.URL.Query()

	wType := queryParams.Get("wType")
	wDate := queryParams.Get("wDate")
	month := queryParams.Get("month")
	comments := queryParams.Get("comments")

	var paramName string
	var param string

	if len(wType) > 0 {
		param = wType
		paramName = "workout_type"
	} else if len(wDate) > 0 {
		param = wDate
		paramName = "workout_date"
	} else if len(month) > 0 {
		param = month
		paramName = "month"
	} else if len(comments) > 0 {
		param = comments
		paramName = "comments"
	}

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

	//var workouts []bson.D
	//var workouts []Workout
	workouts := make([]Workout, 0)

	// filter / query
	filter := bson.D{{paramName, param}}
	findOptions := options.Find().SetSort(bson.D{{"record", -1}}) // find using filter and sort
	//cursor, err := coll.Find(context.TODO(), bson.D{{"workout_type", wType}})	// find by type
	cursor, err := coll.Find(context.TODO(), filter, findOptions)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the %s %s \n", paramName, param)
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

	jsonData, err3 := json.Marshal(workouts)
	if err3 != nil {
		log.Panicln(err)
	}
	req.Header.Add("Content-Type", "application/json")
	//req.Header.Add("Accept", "application/json")
	//req.Response.Header.Add("Content-Type", "application/json")

	fmt.Fprint(writer, string(jsonData))
}
