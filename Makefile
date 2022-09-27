deploy:
	~/google-cloud-sdk/bin/gcloud functions deploy go_gcp_cfunc_mongo_workouts --trigger-http --runtime=go116 --entry-point=GetWorkouts --allow-unauthenticated --memory=256MB