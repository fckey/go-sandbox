# Docker based operation samples

Start to use GCP with gcloud, running Python, or the other operation often requires installation of the underlying toolsets. 

Docker enable you to reproduce any operations with environment agnostic but you should just need to have Docker runtime.  

# Prerequisite 

1. Install Docker

    Mac:  https://www.docker.com/products/docker-desktop

2. Authenticate gcloud 

    ```
       docker-compose run --rm gcloud gcloud auth login
    ``` 


# Scenarios
## Run Cloud Build 
Runs Cloud Build from a container.

How to run:
```shell script
docker-compose run --rm gcloud /bin/sh cloudbuild.sh ${GCP_PROJECT_ID}
```

This scenario uses:

- "gcloud" in docker-compose file
- cloudbuild.sh - Shell script to start build
- cloudbuild.yaml  - Config of Cloud build

Below steps are done in the container:

1. Create file and directories
2. Run build based on config file
3. Local and remote shows the output 


## Run Go
This runs Go script in a container.

How to run: 

```shell script
docker-compose run --rm go_sample
```

This scenario uses:

- "go_sample" in docker-compose file
- sample.go - Go file ran in the container

## Create Cloud Composer via terraform
Required Setup:
- Create service account which has necessary permission. For this example which create Composer environment follow below
  1. Create service account like "terraform-dev@" 
  1. Give "Composer Administrator" and "Editor"
  1. Download key and save as terra/terraform-key.json

How to Run:
```shell script
# Just plan
docker-compose run --rm create_comp  plan

# Create 
docker-compose run --rm create_comp  apply
```