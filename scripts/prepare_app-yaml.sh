#!/bin/bash

ENVIRONMENT=$1
BUCKET_NAME=$2
IMAGE_ACCESS_ENDPOINT=$3
AUTHENTICATION_METHOD=$4
PROJECT_ID=$5

if [ "$#" -ne 5 ]; then
    echo "Incorrect number of parameters"
    exit 1
fi

echo "Preparing app.yaml .."
scripts/substitute.py app-template.yaml --output app.yaml --values environment=$ENVIRONMENT \
                                                                   bucketname=$BUCKET_NAME \
                                                                   imageaccessendpoint=$IMAGE_ACCESS_ENDPOINT \
                                                                   authenticationmethod=$AUTHENTICATION_METHOD \
                                                                   projectid=$PROJECT_ID
