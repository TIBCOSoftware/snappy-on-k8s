# Helm chart for Zeppelin notebook enabled Spark cluster
This chart is still work-in-progress

## Chart Details

* This chart launches a Zeppelin server and exposes Zeppelin as a service on port 8080
* Zeppelin notebook environment makes use of [spark-on-k8s](https://github.com/apache-spark-on-k8s/spark) project binaries to launch Spark jobs. 
* Note that the Spark cluster is only started lazily. You will see Spark driver and executor pods being launched when you run a Spark paragraph in Zeppelin. 
* In order to launch Spark jobs, this makes use of Spark's in-cluster client mode. Here in-cluster client mode means  the submission environment is within a Kubernetes cluster. For the in-cluster mode, this chart uses binaries/docker images built after applying patch for [PR#456](https://github.com/apache-spark-on-k8s/spark/pull/456) as the changes for it are not yet merged in [spark-on-k8s](https://github.com/apache-spark-on-k8s/spark) project.

## Installing the Chart
1. Make sure that Helm is setup in your Kubernetes environment. For example by following instructions at https://docs.bitnami.com/kubernetes/get-started-kubernetes/#step-4-install-helm-and-tiller
2. Get the chart and install it. 

```
  $ git clone https://github.com/SnappyDataInc/spark-on-k8s
  $ cd charts
  $ helm install --name example ./spark-k8s-zeppelin-chart/
```
The above command will deploy the helm chart and will display instructions to access Zeppelin service and Spark UI.

3. Accessing Zeppelin environment: Zeppelin service URL can be found by running following commands

*NOTE*: It may take a few minutes for the LoadBalancer IP to be available. You can watch the status of by running `kubectl get svc -w example-spark-k8s-zeppelin-chart`

```
	$ export ZEPPELIN_SERVICE_IP=$(kubectl get svc --namespace default example-spark-k8s-zeppelin-chart -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
	$ echo http://$ZEPPELIN_SERVICE_IP:8080
```
Use the URL printed by the above command to access Zeppelin.

4. Accessing Spark UI: Once Spark job is launched, Spark UI can be accessed by using following commands

NOTE: Below command will work ONLY after a Spark job is launched. Run a Spark Notebook, first. 
```
   export SPARK_UI_SERVICE_IP=$(kubectl get svc --namespace default example-spark-web-ui -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
   echo http://$SPARK_UI_SERVICE_IP:4040
```
Use the URL printed by the above command to access Spark UI.

## Configuration
TBD
