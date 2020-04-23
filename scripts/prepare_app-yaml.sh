#!/bin/bash

ENVIRONMENT=$1

if [ "$#" -ne 8 ]; then
    echo "Incorrect number of parameters"
    exit 1
fi

echo "Preparing app.yaml .."
scripts/substitute.py app-template.yaml --output app.yaml --values environment=$ENVIRONMENT
