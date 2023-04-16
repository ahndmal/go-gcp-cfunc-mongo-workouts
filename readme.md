## Workout Android App

GCP serverless function to get all workouts by:
- type
- date

### Query params
- wType
- wDate
- month
- year
- comments

### Deploy
```
~/google-cloud-sdk/bin/gcloud functions deploy go_gcp_cfunc_mongo_workouts --trigger-http --runtime=go116 --entry-point=GetWorkouts --allow-unauthenticated --memory=256MB
```
### LOGS
```google cloud
gcloud functions logs read
gcloud functions logs read FUNCTION_NAME --execution-id EXECUTION_ID
```
### Mongo

Install Mongo Tools

```bash
mongodb+srv://<credentials>@cluster0.t1yi6.mongodb.net/myFirstDatabase?appName=mongosh+1.6.0
# import
mongoimport --db=workouts --collection=workouts --type=csv --headerline --file=/home/malandr/Downloads/'workouts-Grid view.csv'
mongosh "mongodb+srv://cluster0.t1yi6.mongodb.net/myFirstDatabase" --apiVersion 1 --username <usern
```
