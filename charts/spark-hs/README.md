# A Helm chart for Spark History Server
This chart is still work-in-progress

## Chart Details
This chart launches a Spark History Server on Kubernetes. To read Spark history events from a file system path, user needs to provide name of an exiting persistent volume claim.
A user can also set HDFS/S3 URI instead in SPARK_HISTORY_OPTS environment variable (refer to values.yaml)
 
## Installing the Chart

**Instructions to use with spark-k8s-zeppelin Helm chart**

1. First install the spark-k8s-zeppelin Helm chart by following the instructions given [here](https://github.com/SnappyDataInc/spark-on-k8s/tree/master/charts/spark-k8s-zeppelin-chart).
   For example:
 
    ```
    $ helm install --name example ./spark-k8s-zeppelin-chart/
    ```

2. After the installation of the spark-k8s-zeppelin chart, a persistent volume claim will be created for the standard PV. Get the name of the PVC

    ```
    $ kubectl get pvc
    NAME                               STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
    example-spark-k8s-zeppelin-chart   Bound     pvc-7627ffa9-11a5-11e8-8f54-42010a8001f4   8Gi        RWO            standard       3m
    ```

   Use the above persistent volume claim name to install History Server chart. This PVC will be used to mount a persistent volume from which History Server events will be read

   For example:

    ```
    $ helm install --name example-history --set eventsdir.existingClaimName=example-spark-k8s-zeppelin-chart ./spark-hs/
    ```

3.  Spark History UI URL can now be accessed as follows:
    ```
    $ export SERVICE_IP=$(kubectl get svc --namespace default example-history-spark-hs -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
    $ echo http://$SERVICE_IP:18080
    ```
    Use the URL to access History UI in a browser.
     
## Deleting the chart
Use `helm delete` command to delete the chart
   ```
   $ helm delete --purge example-history
   ```