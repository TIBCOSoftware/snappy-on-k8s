#!/bin/bash

export MASTER=k8s://$KUBERNETES_SERVICE_HOST:$KUBERNETES_SERVICE_PORT
export SPARK_SUBMIT_OPTIONS="--conf spark.kubernetes.namespace=default --conf spark.kubernetes.driver.pod.name=$HOSTNAME --conf spark.kubernetes.container.image=snappydatainc/spark:v2.4-dist"
