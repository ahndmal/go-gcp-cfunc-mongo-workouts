package p

type Workout struct {
	Id           string `json:"id"`
	Record       int64  `json:"record"`
	Sets         string `json:"sets"`
	Comments     string `json:"comments"`
	CreationDate string `json:"creation_date"`
	WorkoutDate  string `json:"workout_date"`
	Day          string `json:"day"`
	Week         int    `json:"week"`
	WorkoutType  string `json:"workout_type"`
	Month        string `json:"month"`
}
