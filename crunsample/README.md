# Overview
This is sample usage of Cloud Run in GCP. This module contain usecase of ,Twitter API, Cloud PubSub, and Slack API via managers in the repository.

# Getting stargted

## Run go process
1. Add 'go' directory to your GOAPTH by following command  ```export GOPATH="${GOPATH}:${PATH_TO_PROJECT}/go"```

1. Install necessary dependency by ```dep ensure --vendor-only``` If dep is not installed, use ```go get -u github.com/golang/dep/cmd/dep```

1.  You are ready to run main.go

1. Each libray has own README so please go through them 

#

# Use in GCP
## Initialize to use with GCP
Set environment variables below before starting the app. They are included in scripts/env.sh
```
export GCP_PROJECT_ID=""
export SERVICE_ACCOUNT=""
export JOB_NAME=""
export SERVICE_ACCOUNT_EMAIL=$SERVICE_ACCOUNT@$GCP_PROJECT_ID.iam.gserviceaccount.com

```

## How to build
https://cloud.google.com/run/docs/building/containers
```
gcloud builds submit --project=$GCP_PROJECT_ID --tag gcr.io/$GCP_PROJECT_ID/$JOB_NAME
```
## How to deploy
https://cloud.google.com/run/docs/deploying
```
gcloud beta run deploy $JOB_NAME --project=$GCP_PROJECT_ID --image=gcr.io/$GCP_PROJECT_ID/$JOB_NAME:latest
```

With environment variables
https://cloud.google.com/run/docs/configuring/environment-variables
```
gcloud beta run deploy $JOB_NAME --project=$GCP_PROJECT_ID --image=gcr.io/$GCP_PROJECT_ID/$JOB_NAME:latest \
--update-env-vars="TWITTER_CONSUMER_KEY=$TWITTER_CONSUMER_KEY,TWITTER_CONSUMER_SECRET=$TWITTER_CONSUMER_SECRET,TWITTER_ACCESS_TOKEN=$TWITTER_ACCESS_TOKEN,TWITTER_ACCESS_TOKEN_SECRET=$TWITTER_ACCESS_TOKEN_SECRET"
```
## How to call
```
curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" URL
```
## How to test the docker file in the local
https://cloud.google.com/run/docs/testing/local
```
PORT=8080 && docker run -p 8080:${PORT} -e PORT=${PORT} gcr.io/$GCP_PROJECT_ID/$JOB_NAME:latest
```

## Set up Cloud Scheduler
Set up Cloud Run end point can be periodically called by Cloud Scheduler as batch job.

### Set up Service Account
To call non-public Cloud Run end point, ["run.invoker"](https://cloud.google.com/run/docs/securing/managing-access) need to be set. 
This section describes how to create service account which can invoke the instance.  

[IAM Docs](https://cloud.google.com/iam/docs/creating-managing-service-accounts)


```
SERVICE_ACCOUNT="twitter-analysis"
JOB_NAME="crunsample"
SERVICE_ACCOUNT_EMAIL=$SERVICE_ACCOUNT@$GCP_PROJECT_ID.iam.gserviceaccount.com
gcloud iam service-accounts create $SERVICE_ACCOUNT --project=$GCP_PROJECT_ID --display-name="Twitter Analysis"
gcloud beta run services add-iam-policy-binding $JOB_NAME \
    --member=serviceAccount:$SERVICE_ACCOUNT_EMAIL \
    --role="roles/run.invoker" --project=$GCP_PROJECT_ID
```

### Cloud Scheduler

``` 
gcloud scheduler jobs create http ${JOB_ID} --schedule="every 10 mins" --uri=${URI} --oidc-service-account-email=${SERVICE_ACCOUNT_EMAIL} --project=$GCP_PROJECT_ID
```