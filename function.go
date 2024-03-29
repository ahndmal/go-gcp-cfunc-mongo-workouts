package p

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"time"
)

func GetWorkouts(writer http.ResponseWriter, req *http.Request) {
	queryParams := req.URL.Query()

	writer.Header().Set("Access-Control-Allow-Origin", "https://workouts-web-static.vercel.app")

	wType := queryParams.Get("wType")
	wDate := queryParams.Get("wDate")
	month := queryParams.Get("month")
	year := queryParams.Get("year")
	comments := queryParams.Get("comments")

	var paramName string
	var paramName2 string
	var param string
	var param2 string

	if len(wType) > 0 {
		paramName = "workout_type"
		param = wType
	} else if len(wDate) > 0 {
		paramName = "workout_date"
		param = wDate
	} else if len(month) > 0 && len(year) > 0 {

		// todo: debug
		log.Println(">>>>>> len(month) > 0 && len(year) > 0 ")

		paramName = "month"
		paramName2 = "year"
		param = month
		param2 = year
	} else if len(comments) > 0 {
		paramName = "comments"
		param = comments
	}
	log.Printf(">>>>> param is { %s: %s }", paramName, param)
	log.Printf(">>>>> param2 is { %s: %s } ", paramName2, param2)

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
	var workouts []Workout
	workoutDtos := make([]WorkoutDto, 0)

	// by month and year
	if paramName == "month" && paramName2 == "year" {
		//filter := bson.D{{paramName, param}, bson.M{"createdAt": bson.M{
		//	"$gte": primitive.NewDateTimeFromTime(time.Now())
		//}}}

		//"01 Oct 22 15:04 MST"
		//workDate, err := time.Parse(time.RFC822, fmt.Sprintf("01 %s %s 00:00 MST", month[0:3], year[2:4]))

		workDate, err := time.Parse(time.RFC3339Nano, fmt.Sprintf("%s-%s-01T00:00:00.000+00:00", param2, GetMonthNumByName(param)))

		// debug
		log.Printf(">>>>> parsed DATE from 'RFC3339Nano' is : %s", workDate)

		// 2023-08-25T12:25:22.102+00:00

		cursor, err4 := coll.Find(context.TODO(), bson.M{"created_at": bson.M{
			"$gte": primitive.NewDateTimeFromTime(workDate),
		}})

		//findOptions := options.Find().SetSort(bson.D{{"record", -1}}) // find using filter and sort
		//cursor, err4 := coll.Find(context.TODO(), filter, findOptions)

		if errors.Is(err4, mongo.ErrNoDocuments) {
			fmt.Printf("No document was found with the %s %s \n", paramName, param)
			return
		}

		err2 := cursor.All(context.TODO(), &workouts)
		if err2 != nil {
			log.Panicln(err2)
		}

		for _, workout := range workouts {
			dto := convertWorkout(workout)
			workoutDtos = append(workoutDtos, dto)
		}

		jsonData, err3 := json.MarshalIndent(workoutDtos, "", "    ")
		if err3 != nil {
			panic(err)
		}

		req.Header.Add("Content-Type", "application/json")

		_, err5 := fmt.Fprint(writer, string(jsonData))
		if err5 != nil {
			log.Fatalln(err5)
		}

		// by 1 parameter
	} else if len(paramName) > 0 {
		// filter / query
		filter := bson.D{{paramName, param}}
		findOptions := options.Find().
			SetSort(bson.D{{"record", -1}}) // find using filter and sort

		//cursor, err := coll.Find(context.TODO(), bson.D{{"workout_type", wType}})	// find by type
		cursor, err4 := coll.Find(context.TODO(), filter, findOptions)
		if errors.Is(err4, mongo.ErrNoDocuments) {
			fmt.Printf("No document was found with the %s %s \n", paramName, param)
			return
		}
		err2 := cursor.All(context.TODO(), &workouts)
		if err2 != nil {
			log.Panicln(err2)
		}

		for _, workout := range workouts {
			dto := convertWorkout(workout)
			workoutDtos = append(workoutDtos, dto)
		}

		jsonData, err3 := json.MarshalIndent(workoutDtos, "", "    ")
		if err3 != nil {
			panic(err)
		}

		req.Header.Add("Content-Type", "application/json")
		//req.Header.Add("Accept", "application/json")
		//req.Response.Header.Add("Content-Type", "application/json")

		_, err5 := fmt.Fprint(writer, string(jsonData))
		if err5 != nil {
			log.Fatalln(err5)
		}
	} else {
		filter := bson.D{}
		findOptions := options.Find().SetSort(bson.D{{"record", -1}}) // find using filter and sort
		cursor, err := coll.Find(context.TODO(), filter, findOptions)

		if err != nil {
			log.Println(err)
		}

		err2 := cursor.All(context.TODO(), &workouts)
		if err2 != nil {
			log.Printf(">> Error when writing data from ALL on mongo cursor to workouts: %v", err2)
		}

		for _, workout := range workouts {
			dto := convertWorkout(workout)
			workoutDtos = append(workoutDtos, dto)
		}

		jsonData, err3 := json.MarshalIndent(workoutDtos, "", "    ")
		if err3 != nil {
			panic(err)
		}

		req.Header.Add("Content-Type", "application/json")
		//req.Header.Add("Accept", "application/json")
		//req.Response.Header.Add("Content-Type", "application/json")

		_, err5 := fmt.Fprint(writer, string(jsonData))
		if err5 != nil {
			log.Fatalln(err5)
		}
	}
}

func convertWorkout(workout Workout) WorkoutDto {
	//cDateStr, err := time.Parse(time.RFC822, workout.CreationDate)
	var wDateStr string
	//wDateStr = workout.CreationDate.String()
	wDateStr = workout.CreationDate.Format(time.RFC1123)

	dto := WorkoutDto{
		Id:           workout.Id,
		Record:       workout.Record,
		Sets:         workout.Sets,
		Comments:     workout.Comments,
		CreationDate: wDateStr,
		WorkoutDate:  workout.WorkoutDate,
		Day:          workout.Day,
		Week:         workout.Week,
		WorkoutType:  workout.WorkoutType,
		Month:        workout.Month,
		Year:         workout.Year,
	}
	return dto
}
