#!/bin/sh
# This demonstrates `gcloud builds sumit .` symlinks.

if ! gcloud auth print-access-token >/dev/null 2>&1; then
    echo 'You have to run the following command first:' >&2
    echo '  docker-compose run --rm gcloud gcloud auth login' >&2
    exit 1
fi

if [ -z "${1}" ]; then
    echo "Specify GCP project id as an argument."
    exit 1
fi

mkdir mydir1
mkdir mydir2
echo "File in mydir2" >mydir2/file.txt
ln -s ../mydir2 mydir1/link_to_mydir2

echo
echo "Contents:"
ls -lR .

echo "Contents of mydir1/:"
ls -lR mydir1/

echo
echo "Contents of mydir1/link_to_mydir2:"
ls -lR mydir1/link_to_mydir2/

echo
echo "Cloud Build:"
gcloud builds submit --config=cloudbuild.yaml --project=${1} .
