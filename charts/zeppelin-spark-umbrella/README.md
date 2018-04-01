# Umbrella Helm chart to deploy Zeppelin
This chart deploys Zeppelin, Spark History Server, Spark Resource Staging Server, 
and Spark Shuffle Service on Kubernetes. This chart uses individual subcharts for 
the respective components. A user can configure each subchart in parent chart's 
`values.yaml` file.

# Usage
1. Get the dependencies (subcharts) required by the chart
  ```
  helm dep up zeppelin-spark-umbrella
  ```
  The above command will bring in dependent subcharts in the charts directory. For example:
  
  ```
  $ ls zeppelin-spark-umbrella/charts/
  jupyter-with-spark-0.1.0.tgz spark-hs-0.1.0.tgz  spark-k8s-zeppelin-chart-0.1.0.tgz  spark-rss-0.1.0.tgz  spark-shuffle-0.1.0.tgz

  ```
  
2. Copy you json key file to access the GCS bucket in the secrets directory of parent chart.
  For example:
  
  ```
  cp sparkonk8s-test.json zeppelin-spark-umbrella/secrets/
  ```

3. Modify the  `values.yaml` file. Specify the following required attributes for history server configuration
    1. historyServer.historyServerConf.eventsDir: specify the GCS bucket path from which the history 
    server will read logs. For example:
    ```
      historyServerConf:
        eventsDir: "gs://spark_history_server_testing/"
    ```
    2. historyServer.environment: configure SPARK_HISTORY_OPTS to specify the JSON key file. JSON key file will 
    be available in path /etc/secrets of the pod. For example:
     ```
     environment:
       SPARK_HISTORY_OPTS: -Dspark.hadoop.google.cloud.auth.service.account.json.keyfile=/etc/secrets/sparkonk8s-test.json
     ```
     3. Similar to step 2, configure Zeppelin to log events to the same GCS bucket
     ```
     zeppelin:
       environment:
         SPARK_SUBMIT_OPTIONS: >-
            --kubernetes-namespace default
            --conf spark.kubernetes.driver.docker.image=snappydatainc/spark-driver:v2.2.0-kubernetes-0.5.1
            --conf spark.kubernetes.executor.docker.image=snappydatainc/spark-executor:v2.2.0-kubernetes-0.5.1
            --conf spark.executor.instances=2
            --conf spark.hadoop.google.cloud.auth.service.account.json.keyfile=/etc/secrets/sparkonk8s-test.json
       sparkEventLog:
         enableHistoryEvents: true
         eventLogDir: "gs://spark_history_server_testing/"
     ```
     4. For Jupyter notebook server, you can enable history logging by setting sparkEventLog.enableHistoryEvents to true.
     ```python
      sparkEventLog:
        enableHistoryEvents: true
     ```
       
4.  Deploy the umbrella chart
    ```
    helm install --name zp ./zeppelin-spark-umbrella/
    ```
    
5. Jupyter, Zeppelin and History Server URL can be accessed from the K8s UI dashboard (services section)
