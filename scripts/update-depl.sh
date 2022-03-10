#!/bin/bash

sed -i '.bak' "s/\/foodquest-api:.*$/\/foodquest-api:$1/g" ../kubernetes/foodquest-api-depl.yaml
kubectl apply -f ../kubernetes/foodquest-api-depl.yaml
