#!/usr/bin/env bash

kubectl create -f db-service.yaml,db-deployment.yaml,cmpapi-service.yaml,cmpapi-claim0-persistentvolumeclaim.yaml,cmpapi-deployment.yaml