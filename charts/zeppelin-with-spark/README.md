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
  $ helm install --name example --namespace spark ./zeppelin-with-spark/
```
The above command will deploy the helm chart and will display instructions to access Zeppelin service and Spark UI.

*NOTE*: The Spark driver pod uses a Kubernetes default service account to access the Kubernetes API 
server to create and watch executor pods. Depending on the version and setup of Kubernetes deployed, this default 
service account may or may not have the role that allows driver pods to create pods and services. 
Without appropriate role the notebook may fail with a NullPointerException. You may have to grant the default 
service account appropriate permissions, for example:
 ```
 kubectl create clusterrolebinding spark-role --clusterrole=edit --serviceaccount=spark:default --namespace=spark
 ```
 If you want to use a different service account instead of default, you can specify it in the values.yaml file.
 
For more information on RBAC authorization and how to configure Kubernetes service accounts for pods, please refer to
[Using RBAC Authorization](https://kubernetes.io/docs/admin/authorization/rbac/) and
[Configure Service Accounts for Pods](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/).

3. Accessing Zeppelin environment: Zeppelin service URLs can be found by running following commands.

*NOTE*: It may take a few minutes for the LoadBalancer IP to be available. You can watch the status of by running `kubectl get svc -w example-zeppelin`

```
    $ export ZEPPELIN_SERVICE_IP=$(kubectl get svc --namespace default example-zeppelin -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
    $ echo "Access Zeppelin at http://$ZEPPELIN_SERVICE_IP:8080"
    $ echo "Access Spark at http://$ZEPPELIN_SERVICE_IP:4040 after a Spark job is run."
```
*NOTE*: Spark UI will only be available after a Spark job is launched. Run your Spark Notebook first to do that. 


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
	$ cp sparkonk8s-test.json zeppelin-with-spark/conf/secrets/
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
	  --conf spark.kubernetes.driver.docker.image=snappydatainc/spark-driver:v2.2.0-kubernetes-0.5.1
	  --conf spark.kubernetes.executor.docker.image=snappydatainc/spark-executor:v2.2.0-kubernetes-0.5.1
	  --conf spark.executor.instances=2
	  --conf spark.hadoop.google.cloud.auth.service.account.json.keyfile=/etc/secrets/sparkonk8s-test.json
	```
After following the steps as given above, Spark will log job events to the GCS bucket 'gs://spark-history-server/'. You may use the
[Spark History Server UI](https://github.com/SnappyDataInc/spark-on-k8s/tree/master/charts/spark-hs) to view details of jobs.

## Chart Configuration

This section describes various configuration options for the Spark Zeppelin chart. 

### Configuring Persistent Volumes

[Persistent volume](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) can be used to store stateful 
data of an application. For example, users can configure directory for notebook storage, directory for Zeppelin logs 
etc. on persistent volume. By default, this chart [dynamically provisions persistent volume](https://kubernetes.io/docs/concepts/storage/dynamic-provisioning/)
using default storage provisioner. Users can provide `storageClass`, `accessMode` and `size` attributes in `values.yaml` 
file
>Note: This dynamically provisioned volume is deleted by the helm chart when the chart is deleted

Instead of using a dynamically provisioned persistent volume, users can specify persistent volume claim(PVC) for 
an already provisioned volume such as NFS based persistent volume. This can be done by modifying the 
`persistent.existingClaim` attribute.

```
persistence:
  enabled: true
  ## If 'existingClaim' is defined, PVC must be created manually before
  ## volume will be bound
  existingClaim: nfsPVC
```

In the above example, an existing PVC named 'nfsPVC' has been used. This PVC should be created manually before the 
chart is deployed. A persistent volume corresponding to this persistent volume claim will be bound to the pod.
 
>Note: A persistent volume is mounted on the `/data` folder of the Pod. Also setting 'persistence.enabled' to false will
disable provisioning of persistent volume

### Configuring Notebook Storage

Notebook storage directory can be configured by modifying `values.yaml` file. If `noteBookStorage.usePVForNoteBooks`
is set to true, notebooks will be stored on a persistent volume in the path specified by `noteBookStorage.notebookDir` 
attribute

```
noteBookStorage:
  usePVForNoteBooks: true
  # If using PV for notebook storage, 'notebookDir' will be an
  # absolute path in the mounted persistent volume
  notebookDir: "/notebooks"
```


### Configuring Zeppelin

Users can modify the interpreter options by using Zeppelin UI once the Zeppelin service is started. In addition, this 
chart provides two ways to specify Zeppelin properties. Full list of Zeppelin properties can be found in Zeppelin 
[documentation](https://zeppelin.apache.org/docs/0.7.0/install/configuration.html#zeppelin-properties).

#### Using environment variables

Users may specify environment variables for Zeppelin configuration in `values.yaml` file by providing those 
variables in the `environment` section.  For example, 


```
# Any environment variables that need to be made available to the container are defined here
# This may include environment variables used by Spark, Zeppelin
environment:
  # Provide configuration parameters, use syntax as expected by spark-submit
  SPARK_SUBMIT_OPTIONS: >-
     --conf spark.kubernetes.driver.docker.image=snappydatainc/spark-driver:v2.2.0-kubernetes-0.5.1
     --conf spark.kubernetes.executor.docker.image=snappydatainc/spark-executor:v2.2.0-kubernetes-0.5.1
     --conf spark.executor.instances=2
  ZEPPELIN_LOG_DIR: /data/logs
```

In the above example, two environment variables SPARK_SUBMIT_OPTIONS and ZEPPELIN_LOG_DIR used by Zeppelin have been 
configured. The chart will define environment variables in the Zeppelin pod using Kubernetes config maps. 

#### Using configuration files

Another way to configure Zeppelin properties is by specifying Zeppelin configuration files. These files
can be saved in `conf/zeppelin` directory of the chart. Files saved in this directory are mounted on the pod in Zeppelin's 
`conf` folder. 

For example, if a user wants to configure log4j properties, user can modify the log4.properties file in the `conf/zeppelin`
folder. This file will be mounted on the pod and Zeppelin will use the modified properties. Other configuration 
property files can also be specified in the same folder (such as shiro.ini, zeppelin-env.sh). For convenience, template files 
have been provided in the `conf/zeppelin` directory. Users can copy and use these files for configuration.
```
# copy the template file
cp conf/zeppelin/zeppelin-env.sh.template conf/zeppelin/zeppelin-env.sh
# now modify the zeppelin-env.sh
```

### Configuring Spark Jobs
Users can modify the interpreter options by using Zeppelin UI once the Zeppelin service is started. This chart provides 
two ways to specify Spark properties. Full list of Spark properties can be found in Spark documentation

#### Using environment variables

Users may specify environment variables for Spark configuration in `values.yaml` file by providing those variables in 
the `environment` section.  For example,

```
# Any environment variables that need to be made available to the container are defined here
# This may include environment variables used by Spark, Zeppelin
environment:
  # Provide configuration parameters, use syntax as expected by spark-submit
  SPARK_SUBMIT_OPTIONS: >-
     --conf spark.kubernetes.driver.docker.image=snappydatainc/spark-driver:v2.2.0-kubernetes-0.5.1
     --conf spark.kubernetes.executor.docker.image=snappydatainc/spark-executor:v2.2.0-kubernetes-0.5.1
     --conf spark.executor.instances=2
```
In the above example, environment variable SPARK_SUBMIT_OPTIONS used by Zeppelin to configure Spark jobs has been 
provided. The chart will define environment variables in the Zeppelin pod using Kubernetes config maps. 

#### Using configuration files

Another way to configure Spark properties is by specifying Spark configuration files. These files
can be saved in `conf/spark` directory of the chart. Files saved in this directory are mounted on the pod in Spark's 
conf folder. 

For example, if a user wants to configure log4j properties, user can provide the `log4.properties` file in the `conf/spark`
folder. This file will be mounted on the pod and Zeppelin will use the modified properties. Other configuration 
property files can also be specified in the same folder (such as spark-defaults.conf, spark-env.sh). For convenience, 
template files have been provided in the `conf/spark` directory. Users can copy and use these files for configuration.

```
# copy the template file
cp conf/spark/log4j.properties.template conf/spark/log4j.properties
# now modify the zeppelin-env.sh
```

## Configuration Properties List
The following table lists the configuration parameters available for this chart

| Parameter               | Description                        | Default                                                    |
| ----------------------- | ---------------------------------- | ---------------------------------------------------------- |
| `image.repository`      |  Docker repo for the image         |     `SnappyDataInc`                                        |
| `image.tag`             |  Tag for the Docker image          |     `zeppelin:0.7.3-spark-v2.2.0-kubernetes-0.5.0`         | 
| `image.pullPolicy`      |  Pull policy for the image         |     `IfNotPresent`                                         |
| `zeppelinService.type`  |  K8S service type for Zeppelin     |     `LoadBalancer`                                         |
| `zeppelinService.zeppelinPort`  |  Port for Zeppelin service         |      `8080`                                   |
| `zeppelinService.sparkUIPort`   |  Port for Spark service            |      `4040`                                   |
| `serviceAccount`        |  Service account used to deploy Zeppelin and run Spark jobs |     `default`                                    |
| `environment`           |  Environment variables that need to be defined in containers. For example, those required by Spark and Zeppelin |        |
| `environment.SPARK_SUBMIT_OPTIONS` | Configuration options for spark-submit, used by Zeppelin while running Spark jobs | |
| `sparkEventLog.enableHistoryEvents` | Whether to enable Spark event logging, required by History server |  `false` |
| `sparkEventLog.eventLogDir` | URL of the GCS bucket where Spark event logs will be written | |
| `mountSecrets` | If true, files in 'secrets' directory will be mounted on path /etc/secrets | `false` |   |
| `noteBookStorage.usePVForNoteBooks` | Whether to use persistent volume to store notebooks | `true`|
| `noteBookStorage.notebookDir` | Absoluter path on the mounted persistent volume where notebooks will be stored | `/notebooks` |
| `persistence.enabled` | Whether to mount a persistent volume | `true` |
| `persistence.existingClaim` | An existing PVC to be used while mounting a PV. If this is not specified a dynamic PV will be created and a PVC will be generated for it |  |
| `persistence.storageClass`  | Storage class to be used for creating a PVC, if an `existingClaim` is not specified. If unspecified default will be chosen. Ignored if `persistence.existingClaim` is specified | |
| `persistence.accessMode`    | Access mode for the dynamically generated PV and its PVC. Ignored if `persistence.existingClaim` is specified | `ReadWriteOnce`|
| `persistence.size`    | Size of the dynamically generated PV and its PVC. Ignored if `persistence.existingClaim` is specified | 8Gi |
| `resources`           | CPU and Memory resources for the Zeppelin pod  | |
| `global.umbrellaChart` | Internal attribute. Do not modify | `false` | 

These configuration attributes can be set in the `values.yaml` file or while using the helm install command, for example, 

```
# set an attribute while using helm install command
helm install --name zeppelin --namespace spark --set serviceAccount=spark ./zeppelin-with-spark
```
 
