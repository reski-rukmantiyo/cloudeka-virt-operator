#!/bin/bash

old_file="config/webhook/manifests.old.yaml"
new_file="config/webhook/manifests.yaml"

cp $old_file $new_file
kubectl apply -f config/webhook/manifests.yaml 