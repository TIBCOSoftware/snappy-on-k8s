# Helm chart for Zeppelin notebook enabled Spark cluster
This chart is still work-in-progress

## Chart Details

* This chart launches a Zeppelin server and exposes Zeppelin as a service on port 8080
* Zeppelin notebook environment makes use of [spark-on-k8s](https://github.com/apache-spark-on-k8s/spark) project binaries to launch Spark jobs. 
* Note that the Spark cluster is only started lazily. You will see Spark driver and executor pods being launched when you run a Spark paragraph in Zeppelin. 
* In order to launch Spark jobs, this makes use of Spark's in-cluster client mode. Here in-cluster client mode means  the submission environment is within a Kubernetes cluster. For the in-cluster mode, this chart uses binaries/docker images built after applying patch for [PR#456](https://github.com/apache-spark-on-k8s/spark/pull/456) as the changes for it are not yet merged in [spark-on-k8s](https://github.com/apache-spark-on-k8s/spark) project.

## Installing the Chart
1. Make sure that Helm is setup in your Kubernetes environment. For example by following instructions at https://docs.bitnami.com/kubernetes/get-started-kubernetes/#step-4-install-helm-and-tiller

*NOTE*: In case you have RBAC (Role based access control) enabled on K8S cluster (for example k8s version 1.8.7 on GKE), some additional config will be needed as mentioned [here](https://github.com/kubeflow/tf-operator/issues/106)

[Helm](https://github.com/kubernetes/helm/blob/master/README.md) comprises of two parts: a client and a server (Tiller) inside the kube-system namespace. Tiller runs inside your Kubernetes cluster, and manages releases (installations) of your charts. To be able to do this, Tiller needs access to the Kubernetes API. By default, RBAC policies will not allow Tiller to carry out these operations, so we need to do the following:

Create a Service Account tiller for the Tiller server (in the kube-system namespace). Service Accounts are meant for intra-cluster processes running in Pods.

Bind the cluster-admin ClusterRole to this Service Account. We will use ClusterRoleBindings as we want it to be applicable in all namespaces. The reason is that we want Tiller to manage resources in all namespaces.

Update the existing Tiller deployment (tiller-deploy) to associate its pod with the Service Account tiller.

The cluster-admin ClusterRole exists by default in your Kubernetes cluster, and allows superuser operations in all of the cluster resources. The reason for binding this role is because with Helm charts, you can have deployments consisting of a wide variety of Kubernetes resources.

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
Without appropriate role the notebook may fail with a NullPointerException. You may have to grant the default 
service account appropriate permissions, for example:
 ```
 kubectl create clusterrolebinding spark-role --clusterrole=edit --serviceaccount=default:default --namespace=default
 ```
 If you want to use a different service account instead of default, you can specify it in the values.yaml file.
 
For more information on RBAC authorization and how to configure Kubernetes service accounts for pods, please refer to
[Using RBAC Authorization](https://kubernetes.io/docs/admin/authorization/rbac/) and
[Configure Service Accounts for Pods](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/).

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
   export SPARK_UI_SERVICE_IP=$(kubectl get svc --namespace default example-zeppelin-spark-web-ui -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
   echo http://$SPARK_UI_SERVICE_IP:4040
```
Use the URL printed by the above command to access Spark UI.

### Enabling Spark event logging for history server
By default, Zeppelin chart does not enable Spark event logging. Users can enable Spark event logging to view job 
details in the history server UI by configuring the chart by following the steps given below. Users can log Spark events to any HDFS compatible system (GCS/S3/HDFS etc.). 
Here we list steps to configure Spark event logging on Google Cloud Storage

1. Deploy the Spark History Server chart by following the steps given [here](https://github.com/SnappyDataInc/spark-on-k8s/tree/master/charts/spark-hs) 
   Make sure that you have created a GCS bucket and generated a json key file as given in step 1. of [history server configuration](https://github.com/SnappyDataInc/spark-on-k8s/tree/master/charts/spark-hs#steps-to-configure-and-install-chart)

2. In the values.yaml file of the Zeppelin chart, set 'sparkEventLog.enableHistoryEvents' to true and 'sparkEventLog.eventLogDir' to the URI of the GCS bucket (created in step 1 above).   
	```
	sparkEventLog:
	  enableHistoryEvents: true
	  eventLogDir: "gs://spark-history-server/"
	```
3. To allow Spark Driver to write to GCS bucket, we need to mount the json key file on the driver pod. First copy the json file into 'conf/secrets' directory of spark-k8s-zeppelin chart
	```
	$ cp sparkonk8s-test.json spark-k8s-zeppelin-chart/conf/secrets/
	```
    Also set 'mountSecrets' field of values.yaml file to true. When 'mountSecrets' 
    is set to true json key file will be mounted on path '/etc/secrets' of the pod.  
    
	```
	mountSecrets: true
	```
4.  Lastly in the values.yaml file, modify the attribute 'environment.SPARK_SUBMIT_OPTIONS' so that spark-submit uses json key file while 
    accessing the GCS bucket. In the example below, configuration option 'spark.hadoop.google.cloud.auth.service.account.json.keyfile' is 
    set to the mounted json key file  
    
	```
	environment:
	# Provide configuration parameters, use syntax as expected by spark-submit
	SPARK_SUBMIT_OPTIONS: >-
	  --kubernetes-namespace default
	  --conf spark.kubernetes.driver.docker.image=snappydatainc/spark-driver:v2.2.0-kubernetes-0.5.1
	  --conf spark.kubernetes.executor.docker.image=snappydatainc/spark-executor:v2.2.0-kubernetes-0.5.1
	  --conf spark.executor.instances=2
	  --conf spark.hadoop.google.cloud.auth.service.account.json.keyfile=/etc/secrets/sparkonk8s-test.json
	```
After following the steps as given above, Spark will log job events to the GCS bucket 'gs://spark-history-server/'. You may use the
[Spark History Server UI](https://github.com/SnappyDataInc/spark-on-k8s/tree/master/charts/spark-hs) to view details of jobs.

## Configuration
The following table lists the configuration parameters available for this chart

| Parameter               | Description                        | Default                                                    |
| ----------------------- | ---------------------------------- | ---------------------------------------------------------- |
| `replicaCount`          |  No of replicas of Zeppelin server |     `1`                                                    |
| `image.repository`      |  Docker repo for the image         |     `SnappyDataInc`                                        |
| `image.tag`             |  Tag for the Docker image          |     `zeppelin:0.7.3-spark-v2.2.0-kubernetes-0.5.0`         | 
| `image.pullPolicy`      |  Pull policy for the image         |     `IfNotPresent`                                         |
| `zeppelinService.type`  |  K8S service type for Zeppelin     |     `LoadBalancer`                                         |
| `zeppelinService.port`  |  Port for Zeppelin service         |      `8080`                                                |
| `sparkWebUI.type`       |  K8S service type for for Spark UI |     `LoadBalancer`                                         |
| `sparkWebUI.port`       |  Port for Spark service                          |     `4040`                                   |
| `serviceAccount`        |  Service account used to deploy Zeppelin and run Spark jobs |     `default`                                    |
| `environment`           |  Environment variables that need to be defined in containers. For example, those required by Spark and Zeppelin |        |
| `environment.SPARK_SUBMIT_OPTIONS` | Configuration options for spark-submit, used by Zeppelin while running Spark jobs | |
| `sparkEventLog.enableHistoryEvents` | Whether to enable Spark event logging, required by History server |  `false` |
| `sparkEventLog.eventLogDir` | URL of the GCS bucket where Spark event logs will be written | |
| `mountSecrets` | If true, files in 'secrets' directory will be mounted on path /etc/secrets | `false` |   |
| `noteBookStorage.usePVForNoteBooks` | Whether to use persistent volume to store notebooks | `true`|
| `noteBookStorage.notebookDir` | Absoluter path on the mounted persistent volume where notebooks will be stored | `/notebooks` |
| `persistence.enabled` | Whether to mount a persistent volume | `true` |
| `persistence.existingClaim` | An existing PVC to be used while mounting a PV. If this is not specified a dynamic PV will be created and a PVC will be generated for it | - |
| `persistence.storageClass`  | Storage class to be used for creating a PVC, if an `existingClaim` is not specified. If unspecified default will be chosen. Ignored if `persistence.existingClaim` is specified | |
| `persistence.accessMode`    | Access mode for the dynamically generated PV and its PVC. Ignored if `persistence.existingClaim` is specified | `ReadWriteOnce`|
| `persistence.size`    | Size of the dynamically generated PV and its PVC. Ignored if `persistence.existingClaim` is specified | 8Gi |
| `resources`           | CPU and Memory resources for the Zeppelin pod  | |
| `global.umbrellaChart` | Internal attribute. Do not modify | `false` | 

These configuration attributes can be set in the `values.yaml` file or while using the helm install command, for example, 

```
# set an attribute while using helm install command
helm install --name zeppelin --set serviceAccount=spark ./spark-k8s-zeppelin-chart
```


