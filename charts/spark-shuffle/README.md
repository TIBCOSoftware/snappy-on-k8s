# A Helm chart to launch Spark shuffle service on k8s

This chart launches Spark shuffle service as a daemonset on Kubernetes. Spark shuffle service is 
 required for dynamic executor scaling in Spark


## Installing the Chart

To install the chart use following command:

  ```
  $ helm install --name example-shuffle ./spark-shuffle/
  ```

Above command will deploy the chart and display labels attached to the Shuffle pods in 'NOTES'.

Example output:

  ```
  NOTES:
  Created a Spark shuffle daemonset with following Pod labels:
   app: spark-shuffle-service
   spark-version: 2.2.0

  ```
Above mentioned labels can be used in spark-submit configuration options to enable dynamic executor scaling
  
For example in values.yaml of [spark-k8s-zeppelin chart](https://github.com/SnappyDataInc/spark-on-k8s/tree/master/charts/spark-k8s-zeppelin-chart), modify the SPARK_SUBMIT_OPTIONS
 as given below (note the options for dynamicAllocation and shuffle). This will enable dynamic executor scaling and use the 
 shuffle service installed above.

```
  SPARK_SUBMIT_OPTIONS: >-
     --kubernetes-namespace default
     --conf spark.kubernetes.driver.docker.image=snappydatainc/spark-driver:v2.2.0-kubernetes-0.5.1
     --conf spark.kubernetes.executor.docker.image=snappydatainc/spark-executor:v2.2.0-kubernetes-0.5.1
     --conf spark.driver.cores="300m"
     --conf spark.dynamicAllocation.enabled=true
     --conf spark.shuffle.service.enabled=true
     --conf spark.kubernetes.shuffle.namespace=default
     --conf spark.kubernetes.shuffle.labels="app=spark-shuffle-service,spark-version=2.2.0"
     --conf spark.dynamicAllocation.initialExecutors=0
     --conf spark.dynamicAllocation.minExecutors=1
     --conf spark.dynamicAllocation.maxExecutors=5
```
  