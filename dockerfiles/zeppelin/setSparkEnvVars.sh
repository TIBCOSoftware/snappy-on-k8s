#!/bin/bash

export MASTER=k8s://$KUBERNETES_SERVICE_HOST:$KUBERNETES_SERVICE_PORT
export SPARK_SUBMIT_OPTIONS="--kubernetes-namespace default --conf spark.kubernetes.driver.pod.name=$HOSTNAME --conf spark.kubernetes.driver.docker.image=snappydatainc/spark-driver:v2.2.0-kubernetes-0.5.1 --conf spark.kubernetes.executor.docker.image=snappydatainc/spark-executor:v2.2.0-kubernetes-0.5.1"
