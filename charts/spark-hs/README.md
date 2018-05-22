# A Helm chart for Spark History Server
[Spark History Server](https://spark.apache.org/docs/latest/monitoring.html#viewing-after-the-fact) Web UI 
allows users to view job execution details even after the application has finished execution. To use History Server, 
Spark applications should be configured to log events to a directory from which Spark History Server will read events
to construct the job execution visualization. The events directory can be a local file path, an HDFS path, or any alternative 
file system supported by Hadoop APIs. 

## Chart Details
This chart launches Spark History Server on Kubernetes. History server can read events from any 
HDFS compatible system (GCS/S3/HDFS) or a file system path mounted on the pod. A user can set GCS bucket 
URI in 'historyServerConf.eventsDir' attribute. 

*Note:* This README file describes history server configuration to read history events from GCS bucket
 
## Steps to configure and install chart

1. Setup gsutil and gcloud on your local laptop and associate them with your GCP project, create a bucket, 
create an IAM service account sparkonk8s-test, generate a json key file sparkonk8s-test.json, to grant sparkonk8s-test 
admin permission to bucket gs://spark-history-server.

    ```
    $ gsutil mb -c nearline gs://spark-history-server
    $ export ACCOUNT_NAME=sparkonk8s-test
    $ export GCP_PROJECT_ID=project-id
    $ gcloud iam service-accounts create ${ACCOUNT_NAME} --display-name "${ACCOUNT_NAME}"
    $ gcloud iam service-accounts keys create "${ACCOUNT_NAME}.json" --iam-account "${ACCOUNT_NAME}@${GCP_PROJECT_ID}.iam.gserviceaccount.com"
    $ gcloud projects add-iam-policy-binding ${GCP_PROJECT_ID} --member "serviceAccount:${ACCOUNT_NAME}@${GCP_PROJECT_ID}.iam.gserviceaccount.com" --role roles/storage.admin
    $ gsutil iam ch serviceAccount:${ACCOUNT_NAME}@${GCP_PROJECT_ID}.iam.gserviceaccount.com:objectAdmin gs://spark-history-server
    ```
    
2.  In order for history server to be able read from the GCS bucket, we need 
    to mount the json key file on the history server pod. First copy the json file into 'conf/secrets' 
    directory for spark history server chart
    
    ```
    $ cp sparkonk8s-test.json spark-hs/conf/secrets/
    ```
    
3.  Modify the 'values.yaml' file and specify the GCS bucket path created above. History server 
    will read spark events from this path 
    
    ```
    historyServerConf:
      eventsDir: "gs://spark-history-server/"
    ```
        
    Also set 'mountSecrets' field of values.yaml file to true. When 'mountSecrets' 
    is set to true json key file will be mounted on path '/etc/secrets' of the pod.  
    
    ```
    mountSecrets: true
    ```
    
    Lastly set the SPARK_HISTORY_OPTS so that history server uses json key file while 
    accessing the GCS bucket  
    
    ```
    environment:
      SPARK_HISTORY_OPTS: -Dspark.hadoop.google.cloud.auth.service.account.json.keyfile=/etc/secrets/sparkonk8s-test.json
    ```

4.  Install the chart
    ```
    helm install --name history --namespace spark ./spark-hs/
    ```
    
    Spark History UI URL can now be accessed as follows:
    ```
    $ export SERVICE_IP=$(kubectl get svc --namespace default example-history-spark-hs -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
    $ echo http://$SERVICE_IP:18080
    ```
    Use the URL to access History UI in a browser.
        
## Enable spark-submit to log spark history events
The spark-submit example below shows Spark job that logs historical events 
to the GCS bucket created in above steps. Once job finishes, use the 
Spark history server UI to view the job execution details.

  ```
  bin/spark-submit \
      --master k8s://https://<k8s-master-IP> \
      --deploy-mode cluster \
      --name spark-pi \
      --conf spark.kubernetes.namespace=spark \
      --class org.apache.spark.examples.SparkPi \
      --conf spark.eventLog.enabled=true \
      --conf spark.eventLog.dir=gs://spark-history-server/ \
      --conf spark.executor.instances=2 \
      --conf spark.hadoop.google.cloud.auth.service.account.json.keyfile=/etc/secrets/sparkonk8s-test.json \
      --conf spark.kubernetes.driver.secrets.history-secrets=/etc/secrets \
      --conf spark.kubernetes.executor.secrets.history-secrets=/etc/secrets \
      --conf spark.kubernetes.driver.docker.image=snappydatainc/spark-driver:v2.2.0-kubernetes-0.5.1 \
      --conf spark.kubernetes.executor.docker.image=snappydatainc/spark-executor:v2.2.0-kubernetes-0.5.1 \
      local:///opt/spark/examples/jars/spark-examples_2.11-2.2.0-k8s-0.5.0.jar  
  ```

     
## Deleting the chart
Use `helm delete` command to delete the chart
   ```
   $ helm delete --purge example-history
   ```
