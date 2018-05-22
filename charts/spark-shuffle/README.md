# A Helm chart to launch Spark shuffle service on k8s

This chart launches Spark shuffle service as a daemonset on Kubernetes. Spark shuffle service is 
 required for dynamic executor scaling in Spark


## Installing the Chart

To install the chart use following command:

  ```
  $ helm install --name example-shuffle --namespace spark ./spark-shuffle/
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
  
For example in values.yaml of [spark-k8s-zeppelin chart](https://github.com/SnappyDataInc/spark-on-k8s/tree/master/charts/zeppelin-with-spark), modify the SPARK_SUBMIT_OPTIONS
 as given below (note the options for dynamicAllocation and shuffle). This will enable dynamic executor scaling and use the 
 shuffle service installed above.

```
  SPARK_SUBMIT_OPTIONS: >-
     --conf spark.kubernetes.driver.docker.image=snappydatainc/spark-driver:v2.2.0-kubernetes-0.5.1
     --conf spark.kubernetes.executor.docker.image=snappydatainc/spark-executor:v2.2.0-kubernetes-0.5.1
     --conf spark.driver.cores="300m"
     --conf spark.local.dir=/tmp/spark-local
     --conf spark.dynamicAllocation.enabled=true
     --conf spark.shuffle.service.enabled=true
     --conf spark.kubernetes.shuffle.namespace=spark
     --conf spark.kubernetes.shuffle.labels="app=spark-shuffle-service,spark-version=2.2.0"
     --conf spark.dynamicAllocation.initialExecutors=0
     --conf spark.dynamicAllocation.minExecutors=1
     --conf spark.dynamicAllocation.maxExecutors=5
```

## Configuration
The following table lists the configuration parameters available for this chart

| Parameter               | Description                        | Default                                                    |
| ----------------------- | ---------------------------------- | ---------------------------------------------------------- |
| `image.repository`      |  Docker repo for the shuffle service image         |     `SnappyDataInc`                        |
| `image.tag`             |  Tag for the Docker image          |     `spark-shuffle:v2.2.0-kubernetes-0.5.1`        | 
| `image.pullPolicy`      |  Pull policy for the image         |     `IfNotPresent`                                 |
| `serviceAccount`        |  Service account used to deploy shuffle service daemonset |     `default`               |
| `shufflePodLabels` | Labels assigned to pods of the shuffle service. These can be used to target a particular service while running jobs. By default two labels are created `app: spark-shuffle-service` and `spark-version: 2.2.0`| |
| `resources`           | CPU and Memory resources for the RSS pod  | |
| `global.umbrellaChart` | Internal attribute. Do not modify | `false` | 

These configuration attributes can be set in the `values.yaml` file or while using the helm install command, for example, 

```
# set an attribute while using helm install command
helm install --name shuffle --namespace spark --set serviceAccount=spark ./spark-shuffle
```
  