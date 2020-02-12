#!/bin/bash

POD="$(kubectl get pods --selector app=bookstore-1 -nsmc --no-headers | grep 'Running' | awk '{print $1}')"

kubectl port-forward "$POD" -n smc 15000:15000
