package p

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	workouts := make([]bson.M, 0)
	//err = coll.FindOne(context.TODO(), bson.D{{"workout_type", "BACK"}}).Decode(&workout)
	//err = coll.FindOne(context.TODO(), bson.D{}).Decode(&workouts)
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		return
	}
	err2 := cursor.All(context.TODO(), &workouts)
	if err2 != nil {
		return
	}
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title %s\n", "title")
		return
	}
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(workouts, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Fprint(writer, string(jsonData))
}
