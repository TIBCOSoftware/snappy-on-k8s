# Helm chart for Zeppelin notebook enabled Spark cluster
This chart is still work-in-progress

## Chart Details

* This chart launches a Zeppelin server and exposes Zeppelin as a service on port 8080
* Zeppelin notebook environment makes use of [spark-on-k8s](https://github.com/apache-spark-on-k8s/spark) project binaries to launch Spark jobs. 
* Note that the Spark cluster is only started lazily. You will see Spark driver and executor pods being launched when you run a Spark paragraph in Zeppelin. 
* In order to launch Spark jobs, this makes use of Spark's in-cluster client mode. Here in-cluster client mode means  the submission environment is within a Kubernetes cluster. For the in-cluster mode, this chart uses binaries/docker images built after applying patch for [PR#456](https://github.com/apache-spark-on-k8s/spark/pull/456) as the changes for it are not yet merged in [spark-on-k8s](https://github.com/apache-spark-on-k8s/spark) project.

## Installing the Chart
1. Make sure that Helm is setup in your Kubernetes environment. For example by following instructions at https://docs.bitnami.com/kubernetes/get-started-kubernetes/#step-4-install-helm-and-tiller

*NOTE*: In case you have RBAC enabled on K8S cluster (for example k8s version 1.8.7 on GKE), some additional config will be needed as mentioned [here](https://github.com/kubeflow/tf-operator/issues/106)
For example on RBAC enabled cluster run following commands to create a service account for Tiller
```
kubectl create serviceaccount --namespace kube-system tiller
kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
kubectl patch deploy --namespace kube-system tiller-deploy -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'      
helm init --service-account tiller --upgrade
```

2. Get the chart and install it. 

```
  $ git clone https://github.com/SnappyDataInc/spark-on-k8s
  $ cd charts
  $ helm install --name example ./spark-k8s-zeppelin-chart/
```
The above command will deploy the helm chart and will display instructions to access Zeppelin service and Spark UI.

*NOTE*: The Spark driver pod uses a Kubernetes default service account to access the Kubernetes API 
server to create and watch executor pods. Depending on the version and setup of Kubernetes deployed, this default 
service account may or may not have the role that allows driver pods to create pods and services. 
Without appropriate role the notebook may faile with NullPointerException. You may grant the default 
service account appropriate permissions, for example:
 ```
 kubectl create clusterrolebinding default-cluster-rule --clusterrole=cluster-admin --serviceaccount=default:default
 ```
 If you want to use a different service account instead of default, you can specify it in the values.yaml file.

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
