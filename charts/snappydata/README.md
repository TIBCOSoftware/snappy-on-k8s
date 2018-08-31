
# SnappyData on Kubernetes

[SnappyData](http://www.snappydata.io) fuses Apache Spark with an in-memory database to deliver a compute+data engine capable of 
stream processing, transactions, interactive analytics and prediction in a single cluster.

## Key SnappyData Features on top of Spark

- High concurrency
- Low latency, true interactive speeds
- Native, in-memory indexes
- Data integrity through keys and constraints
- Mutability with support for transactions and snapshot isolation
- Shared in-memory state across Spark applications
- Engine for serving predictions
- Extended support for Kubernetes and fully certified for PKS
- HA through distributed consensus based replication
- Support for semi-structured data : Nested data objects, Graph and JSON data
- Framework for Change-data-capture and streaming ETL
- Built-in support for Notebook based development (both Zeppelin and Jupyter)
- Data prep web tool for self-service data cataloging and prep. 
- Enterprise grade security

## Pre-requisites and assumptions:
- We need a running kubernetes or PKS cluster. We only support Kubernetes 1.9 (or higher) and PKS 1.0.0(or higher).
- You must have Helm deployed. Follow the steps given below if Helm is not deployed in your Kubernetes environment

### Setup Helm charts

[Helm](https://github.com/kubernetes/helm/blob/master/README.md) comprises of two parts: a client and a server (Tiller) inside 
the kube-system namespace. Tiller runs inside your Kubernetes cluster, and manages releases (installations) of your charts. 
To install Helm follow the steps [here](https://docs.pivotal.io/runtimes/pks/1-0/configure-tiller-helm.html). The instructions
are applicable for any kubernetes cluster (PKS or GKE or Minikube).

*Note that the Docker image of SnappyData is built using its Open Source (OSS) version. The image of its Enterprise version will be made available soon.*

## Quickstart

We use SnappyData Helm chart to deploy SnappyData on Kubernetes. Use `helm install` command to deploy the SnappyData chart

```
# fetch the chart repo ....
git clone https://github.com/SnappyDataInc/spark-on-k8s
# go to charts directory
cd charts
# install the chart in namespace 'snappy'
helm install --name snappydata --namespace snappy ./snappydata/
```

The above command will deploy the chart in a namespace called 'snappy' and will display the Kubernetes objects created 
by the chart on the console(service, statefulsets etc.).

SnappyData Helm chart uses Kubernetes [statefulsets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/) 
to launch locator, lead and servers. By default SnappyData helm chart will deploy SnappyData cluster consisting of 
1 locator, 1 lead and 2 servers and services to access SnappyData endpoints.

You may monitor the Kubernetes UI dashboard to check the status of the components as it takes a few minutes for all 
servers to be online. SnappyData chart provisions volumes dynamically for servers, locator and lead. These volumes and 
the data on it will be retained even when the chart deployment is deleted.

#### Troubleshooting tip

In case the cluster or some of the member containers fail to launch successfully, users can check SnappyData cluster logs to get more details.

Even if the cluster or a particular member is not running, users can access the SnappyData logs by mounting the persistent volume claims of respective member containers (locator, servers or lead) on a dumb container. We have provided an experimental scripts to do just that.

The script is available [here](../../utils/snappy-debug-pod.sh). It simply takes the name(s) of the persistent volume claim(s) to be mounted and mounts it on the container as /data0 (and /data1, so on).

You can find out the names of persistent volume claims on Kubernetes Dashboard or by running `kubectl get persistentvolumeclaims`.
 
For more details, simply execute the script without any arguments and it'll display its usage.

### Accessing SnappyData UI

SnappyData UI dashboard can be access using `snappydata-leader-public` service. Use URL of the form `externalIp:5050` to 
view SnappyData UI dashboard in a browser. 

To view SnappyData services running in K8S cluster use `kubectl get svc --namespace=snappy` command

```
# check the SnappyData services running in K8S cluster
kubectl get svc --namespace=snappy
# This will show output like following

NAME                        TYPE           CLUSTER-IP      EXTERNAL-IP      PORT(S)                                        AGE
snappydata-leader           ClusterIP      None            <none>           5050/TCP                                       5m
snappydata-leader-public    LoadBalancer   10.51.255.175   35.232.102.51    5050:31964/TCP,8090:32700/TCP,3768:32575/TCP   5m
snappydata-locator          ClusterIP      None            <none>           10334/TCP,1527/TCP                             5m
snappydata-locator-public   LoadBalancer   10.51.241.224   104.198.47.162   1527:31957/TCP                                 5m
snappydata-server           ClusterIP      None            <none>           1527/TCP                                       5m
snappydata-server-public    LoadBalancer   10.51.248.27    35.232.16.4      1527:31853/TCP                                 5m

```

From the above output, the URL  `35.232.102.51:5050` can be used to view dashboard in browser.


### Executing queries using Snappy shell

You may use Snappy shell to connect to SnappyData and execute queries. The Snappy shell need not run withing the K8S 
cluster. 

To run Snappy shell from your laptop/workstation, download the SnappyData distribution from 
[SnappyData github releases page](https://github.com/SnappyDataInc/snappydata/releases)

```
# check the SnappyData services running in K8S cluster
kubectl get svc --namespace=snappy
# This will show output like following

NAME                        TYPE           CLUSTER-IP      EXTERNAL-IP      PORT(S)                                        AGE
snappydata-leader           ClusterIP      None            <none>           5050/TCP                                       5m
snappydata-leader-public    LoadBalancer   10.51.255.175   35.232.102.51    5050:31964/TCP,8090:32700/TCP,3768:32575/TCP   5m
snappydata-locator          ClusterIP      None            <none>           10334/TCP,1527/TCP                             5m
snappydata-locator-public   LoadBalancer   10.51.241.224   104.198.47.162   1527:31957/TCP                                 5m
snappydata-server           ClusterIP      None            <none>           1527/TCP                                       5m
snappydata-server-public    LoadBalancer   10.51.248.27    35.232.16.4      1527:31853/TCP                                 5m

# The above output displays services created for SnappyData along with external IP (if applicable) and ports available 
# for external connections.
# Launch snappy shell using command 'bin/snappy' and connect to locator-public endpoint using URL fo the form 'externalURL:1527'
# For example 
snappy> connect client '104.198.47.162:1527';

# you may now create tables and fire queries
snappy> create table t1(col1 int, col2 int) using column;
snappy> insert into t1 values(1, 1);
1 row inserted/updated/deleted

```

### Connecting using JDBC driver

To connect to the SnappyData cluster, using JDBC, use URL of the form `jdbc:snappydata://<locatorIP>:<locatorClientPort>/`
For Kubernetes deployments JDBC clients can connect to SnappyData system using the locator-public service.

```
# check the SnappyData services running in K8S cluster
kubectl get svc --namespace=snappy
# This will show output like following

NAME                        TYPE           CLUSTER-IP      EXTERNAL-IP      PORT(S)                                        AGE
snappydata-leader           ClusterIP      None            <none>           5050/TCP                                       5m
snappydata-leader-public    LoadBalancer   10.51.255.175   35.232.102.51    5050:31964/TCP,8090:32700/TCP,3768:32575/TCP   5m
snappydata-locator          ClusterIP      None            <none>           10334/TCP,1527/TCP                             5m
snappydata-locator-public   LoadBalancer   10.51.241.224   104.198.47.162   1527:31957/TCP                                 5m
snappydata-server           ClusterIP      None            <none>           1527/TCP                                       5m
snappydata-server-public    LoadBalancer   10.51.248.27    35.232.16.4      1527:31853/TCP                                 5m

```

Use the locator-public service's external IP and port to connect to SnappyData system using JDBC 
connections. For example from the above output the JDBC URL to be used will be  `jdbc:snappydata://104.198.47.162:1527/`

>Note: Refer to [SnappyData documentation](http://snappydatainc.github.io/snappydata/howto/connect_using_jdbc_driver/) for 
example JDBC program and how to obtain JDBC driver using Maven/SBT co-ordinates 

<!--- TODO: Point SnappyData docs?
### Using SnappyData with Zeppelin
--->

<!--- TODO 
### Using Spark shell
--->

### Submitting a SnappyData job

>Note: To understand how to write a SnappyData job, refer to [how-to section in documentation for SnappyData jobs](http://snappydatainc.github.io/snappydata/howto/run_spark_job_inside_cluster/).
To run a job in Kubernetes deployment, follow the the steps given below

To run Snappy job from your laptop/workstation, download the SnappyData distribution from 
[SnappyData github releases page](https://github.com/SnappyDataInc/snappydata/releases). Use the `bin/snappy-job.sh` script 
to submit and run the job in the SnappyData cluster.  For submitting the job use leader-public service that exposes port 
8090 to run jobs.

```
# check the SnappyData services running in K8S cluster
kubectl get svc --namespace=snappy
# This will show output like following

NAME                        TYPE           CLUSTER-IP      EXTERNAL-IP      PORT(S)                                        AGE
snappydata-leader           ClusterIP      None            <none>           5050/TCP                                       5m
snappydata-leader-public    LoadBalancer   10.51.255.175   35.232.102.51    5050:31964/TCP,8090:32700/TCP,3768:32575/TCP   5m
snappydata-locator          ClusterIP      None            <none>           10334/TCP,1527/TCP                             5m
snappydata-locator-public   LoadBalancer   10.51.241.224   104.198.47.162   1527:31957/TCP                                 5m
snappydata-server           ClusterIP      None            <none>           1527/TCP                                       5m
snappydata-server-public    LoadBalancer   10.51.248.27    35.232.16.4      1527:31853/TCP                                 5m

```

Use the `snappydata-leader-public` services's external IP and port 8090 to submit a job. For example:

```
# first cd to your SnappyData product dir
cd $SNAPPY_HOME
# Submit the job using snappydata-leader-public service etxernalIp and port 8090 in the -lead option
# The command below submits a job called CreatePartitionedRowTable as given the SnappyData how-to docs
bin/snappy-job.sh submit
  --app-name CreatePartitionedRowTable
  --class org.apache.spark.examples.snappydata.CreatePartitionedRowTable
  --app-jar examples/jars/quickstart.jar
  --lead 35.232.102.51:8090
```

### Deleting the SnappyData helm chart

Delete the Helm chart using `helm delete` command. Dynamically provisioned volumes and the data on it will be 
retained even the chart deployment is deleted.

```
$ helm delete --purge snappydata
```

## Detailed configuration for SnappyData chart

Users can modify `values.yaml` file to configure the SnappyData chart. Configuration options can be specified in the 
respective attributes of locators, leaders and servers in `values.yaml`

The following table lists the configuration parameters available for this chart


| Parameter| Description | Default |
| ---------| ------------| --------|
| `image` |  Docker repo from which the SnappyData Docker image will be pulled    |  `snappydatainc/snappydata`   |
| `imageTag` |  Tag of the SnappyData Docker image to be pulled |   |
| `imagePullPolicy` | Pull policy for the image.  | `IfNotPresent` |
| `locators.conf` | List of configuration options to be passed to locators | |
| `locators.resources` | Resource configuration for the locator Pods. User can configure CPU/memory requests and limits using this | `locators.requests.memory` is set to 1024Mi |
| `locators.persistence.storageClass` | Storage class to be used while dynamically provisioning a volume | Default value is not defined so `default` storage class for the cluster will be chosen  |
| `locators.persistence.accessMode` | [Access mode](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes) for the dynamically provisioned volume | `ReadWriteOnce` |
| `locators.persistence.size` | Size of the dynamically provisioned volume | `10Gi` |
| `servers.replicaCount` | Number of servers to be started in SnappyData system | `2` |
| `servers.conf` | List of configuration options to be passed to servers | |
| `servers.resources` | Resource configuration for the server Pods. User can configure CPU/memory requests and limits using this | `servers.requests.memory` is set to 4096Mi |
| `servers.persistence.storageClass` | Storage class to be used while dynamically provisioning a volume | Default value is not defined so `default` storage class for the cluster will be chosen  |
| `servers.persistence.accessMode` | [Access mode](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes) for the dynamically provisioned volume | `ReadWriteOnce` |
| `servers.persistence.size` | Size of the dynamically provisioned volume | `10Gi` |
| `leaders.conf` | List of configuration options to be passed to leaders | |
| `leaders.resources` | Resource configuration for the server Pods. User can configure CPU/memory requests and limits using this | `leaders.requests.memory` is set to 4096Mi |
| `leaders.persistence.storageClass` | Storage class to be used while dynamically provisioning a volume | Default value is not defined so `default` storage class for the cluster will be chosen  |
| `leaders.persistence.accessMode` | [Access mode](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes) for the dynamically provisioned volume | `ReadWriteOnce` |
| `leaders.persistence.size` | Size of the dynamically provisioned volume | `10Gi` |


**Example**: start 4 servers each with heap size of 2048MB
```
servers:
  replicaCount: 4
  ## config options for servers
  conf: "-heap-size=2048m"
```

Users may specify SnappyData [configuration parameters](https://snappydatainc.github.io/snappydata/configuring_cluster/configuring_cluster/#configuration-files) 
in the `servers.conf`, `locators.conf` and `leaders.conf` attributes for servers, locators and leaders respectively.


## Description of various Kubernetes objects deployed for SnappyData chart

### Statefulsets for Servers, Leaders and Locators

[Kubernetes statefulsets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/) are used to 
manage stateful applications. Statefulsets provide many benefits such as stable, unique network identifiers, stable 
persistent storage, ordered deployment and scaling, graceful deletion and rolling updates. SnappyData Helm charts 
deploys statefulsets for servers, leaders and locators. The chart by default will create 2 data servers (pods) for 
SnappyData servers and 1 pod each for locator and leader. Upon deletion of the Helm deployment, each pod will gracefully
terminate the SnappyData process running on it.  

### Description of services that expose external endpoints

SnappyData Helm chart creates services in order to allow user to make JDBC connections, execute Spark jobs and access
SnappyData UI etc.  Services of type LoadBalancer have external IP address assigned and can be used to connect from outside 
of Kubernetes cluster.  

To check the service created for SnappyData deployment, use command `kubectl get svc --namespace=snappy`. For example if a SnappyData 
release is deployed using command `helm install --name snappydata ./snappydata` as given in the [Quickstart](#Quickstart)

```
# check the SnappyData services running in K8S cluster
kubectl get svc --namespace=snappy
# This will show output like following

NAME                        TYPE           CLUSTER-IP      EXTERNAL-IP      PORT(S)                                        AGE
snappydata-leader           ClusterIP      None            <none>           5050/TCP                                       5m
snappydata-leader-public    LoadBalancer   10.51.255.175   35.232.102.51    5050:31964/TCP,8090:32700/TCP,3768:32575/TCP   5m
snappydata-locator          ClusterIP      None            <none>           10334/TCP,1527/TCP                             5m
snappydata-locator-public   LoadBalancer   10.51.241.224   104.198.47.162   1527:31957/TCP                                 5m
snappydata-server           ClusterIP      None            <none>           1527/TCP                                       5m
snappydata-server-public    LoadBalancer   10.51.248.27    35.232.16.4      1527:31853/TCP                                 5m

```

In the above out output, three services viz. `snappydata-leader-public`, `snappydata-locator-public` and 
`snappydata-server-public`  of type LoadBalancer have been created. These services have external IPs assigned and therefore 
can be accessed from outside of Kubernetes. Remaining services that have no external ip are created for internal use.
 
`snappydata-leader-public` service exposes port 5050 for SnappyData UI and port 8090 for to accept SnappyData jobs.
`snappydata-locator-public` service exposes port 1527 to accept JDBC connections.
 
### Persistent Volumes 

A Pod in a SnappyData deployment has a persistent volume mounted it. This volume is dynamically provisioned and is used
to store data directory for SnappyData. The persistent volume (PV) is mounted on path `/opt/snappydata/work` in each Pod. 
These volumes and the data on it will be retained even if the chart deployment is deleted. If the chart is deployed again 
with same chart name and the volume exists then the existing volume is used instead of provisioning a new volume.   


## Future work

In future we plan to enable 
 - Backup and Restore in K8s environment
 - Snappy commands such as validate-disk-store, compact-disk-store, compact-all-disk-stores, 
   revoke-missing-disk-store, unblock-disk-store, modify-disk-store, show-disk-store-metadata etc.
 
<!---
## Troubleshooting
--->

