deploy:
	gcloud functions deploy go_gcp_cfunc_mongo_workouts \
		--trigger-http \
		--runtime=go121 \
		--entry-point=GetWorkouts \
		--allow-unauthenticated\
		--memory=128MB