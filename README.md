# Multi cloud Spark application service on PKS

## Introduction
[Kubernetes](kubernetes.io) is an open source project designed specifically for container orchestration. 
Kubernetes offers a number of key [features](https://kubernetes.io/docs/concepts/overview/what-is-kubernetes/), 
including multiple storage APIs, container health checks, manual or automatic scaling, rolling upgrades and 
service discovery.
[Pivotal Container Service - PKS](https://pivotal.io/platform/pivotal-container-service) is a solution to manage 
Kubernetes clusters across private and public clouds. It leverages [BOSH]() to offer a uniform way to instantiate, 
deploy, and manage highly available Kubernetes clusters on a cloud platform like GCP, VMWare vSphere or AWS. 

This project provides a streamlined way of deploying, scaling and managing Spark applications. Spark 2.3 added support
for Kubernetes as a cluster manager. This project leverages [Helm charts](https://helm.sh/) to allow deployment of 
common Spark application recipes - using [Apache Zeppelin](https://zeppelin.apache.org/) and/or [Jupyter](jupyter.org) 
for interactive, collaborative workloads. It also automates logging of all events across batch jobs and Notebook driven
applications to log events to shared storage for offline analysis.   


Helm is a package manager for kubernetes and the most productive way to find, install, share, upgrade and use even the 
most complex kubernetes applications. So, for instance, with a single command you can deploy additional components like 
HDFS or Elastic search for our Spark applications.  

This project is a collaborative effort between SnappyData and Pivotal. 



## Features
- Full support for Spark 2.2 applications running on PKS 1.x on both Google cloud and on-prem VMWare cloud environments.
The project leverages [spark-on-k8s]() work.
- Deploy batch spark jobs using kubernetes master as the cluster/resource manager
- Helm chart to deploy Zeppelin, centralized logging, monitoring across apps (using History server)
- Helm chart to deploy Jupyter,  centralized logging, monitoring across apps (using History server)
- Use kubernetes persistent volumes for notebooks and event logging for collaboration and historical analysis
- Spark applications can be Java, Scala or Python
- Spark applications can dynamically scale

We use Helm charts to abstract the developer from having to understand kubernetes concepts and simply focus on 
configuration that matters. Think recipes that come with sensible defaults for common Spark workloads on PKS. 

We showcase the use of cloud storage (e.g. Google Cloud Storage) to manage logs/events but show how the use persistent 
volumes within the charts make the architecture portable across clouds. 

## Pre-requisites and assumptions:
- We need a running kubernetes or PKS cluster. We only support Kubernetes 1.9 (or higher) and PKS 1.0.1(or higher). If 
you already have access to a Kubernetes cluster skip the next section. 


### Getting access to a PKS or Kubernetes cluster
If you would like to deploy on-prem you can either use Minikube (local developer machine) or get PKS environment setup 
using vSphere. 

#### Option (1) - PKS
- PKS on vSphere: Follow these [instructions](https://docs.pivotal.io/runtimes/pks/1-0/vsphere.html) 
- PKS on GCP: Follow these [instructions](https://docs.pivotal.io/runtimes/pks/1-0/gcp.html)
- Create a Kubernetes cluster using PKS CLI : Once PKS is setup you will need to create a k8s cluster as described 
[here](https://docs.pivotal.io/runtimes/pks/1-0/using.html)

#### Option (2) - Kubernetes on Google Cloud Platform (GCP)
- Login to your Google account and goto the [Cloud console](console.cloud.google.com) to launch a GKE cluster

#### Option (3) - Minikube on your local machine
- If either of the above options is difficult, you may setup a test cluster on your local machine using 
[minikube](https://kubernetes.io/docs/getting-started-guides/minikube/). We recommend using the latest release of minikube 
with the DNS addon enabled.
- If using Minikube, be aware that the default minikube configuration is not enough for running Spark applications. 
We recommend 3 CPUs and 4g of memory to be able to start a simple Spark application with a single executor.

### Steps if a Kubernetes cluster is available 
- If using PKS, you will need to install the PKS command line tool. See instructions 
[here](https://docs.pivotal.io/runtimes/pks/1-0/installing-pks-cli.html)
- Install kubectl on your local development machine and configure access to the kubernetes/PKS cluster. See instructions for 
kubectl [here](https://kubernetes.io/docs/tasks/tools/install-kubectl/). If you are using Google cloud, you will find 
instuctions for setting up Google Cloud SDK ('gcloud') along with kubectl 
[here](https://kubernetes.io/docs/tasks/tools/install-kubectl/).
- You must have appropriate permissions to list, create, edit and delete pods in your cluster. You can verify that you 
can list these resources by running kubectl auth can-i <list|create|edit|delete> pods.
- The service account credentials used by the driver pods must be allowed to create pods, services and configmaps.
- You must have Kubernetes DNS configured in your cluster.


### Setup Helm charts

[Helm](https://github.com/kubernetes/helm/blob/master/README.md) comprises of two parts: a client and a server (Tiller) inside 
the kube-system namespace. Tiller runs inside your Kubernetes cluster, and manages releases (installations) of your charts. 
To install Helm follow the steps [here](https://docs.pivotal.io/runtimes/pks/1-0/configure-tiller-helm.html). The instructions
are applicable for any kubernetes cluster (PKS or GKE or Minikube).


### Quickstart

#### Launch Spark and Notebook servers

We use the zeppelin-spark-umbrella chart to deploy Zeppelin, Spark Resource Staging Server, and Spark Shuffle Service 
on Kubernetes. This chart is composed from individual sub-charts for each of the components. 
You can read more about Helm umbrella charts 
[here](https://github.com/kubernetes/helm/blob/master/docs/charts_tips_and_tricks.md#complex-charts-with-many-dependencies)

You can configure the components in the umbrella chart's 'values.yaml' (see zeppelin-spark-umbrella/values.yaml) or in 
each of the individual subchart's 'values.yaml' file. The umbrella chart's 'values.yaml' will override the component one.

```text
# fetch the chart repo ....
git clone https://github.com/SnappyDataInc/spark-on-k8s

# Now, install the chart
cd charts
helm install --name spark-all ./zeppelin-spark-umbrella/
```
The above command will deploy the helm chart and will display instructions to access Zeppelin service and Spark UI.

> Note that this command will return quickly and kubernetes controllers will work in the background to achieve the state
specified in the chart. The command below can be used to access the notebook environment from any browser. 

```text
kubectl get services -w
# Note: this could take a while to complete. Use '-w' option to wait for state changes. 
```

Once everything is up and running you will see something like this:
```text
NAME                       TYPE           CLUSTER-IP      EXTERNAL-IP      PORT(S)                      AGE
kubernetes                 ClusterIP      10.11.240.1     <none>           443/TCP                      32d
spark-all-rss              LoadBalancer   10.11.249.80    35.194.62.140    10000:31000/TCP              2h
spark-all-spark-web-ui     LoadBalancer   10.11.241.218   35.188.200.10    4040:32033/TCP               2h
spark-all-zeppelin         LoadBalancer   10.11.249.81    35.225.23.173    8080:30927/TCP               2h
```
> Access the zeppelin notebook environment using URL <external-ip>:8080 from any browser.
> Spark UI is accessible using URL <external-ip>:4040. 
NOTE that this will work only after you have run at least one Spark job. Spark Driver (and hence UI) is lazily started. 
Simply navigate to <Zeppelin Home>; Click 'Zeppelin Tutorial' and then 'Basic Features(Spark)'. Run the 'Load
data' paragraph followed by one or more SQL paragraphs. 


#### Launch a Spark batch job

The spark distribution with support for kubernetes can be downloaded 
[here](https://github.com/apache-spark-on-k8s/spark/releases/tag/v2.2.0-kubernetes-0.5.0)

We will use spark-submit from this distribution to deploy a batch job. Example below runs the built in SparkPi job. 
The 'local://' URL will result in looking for the JAR in the launched container. 

```text
# Find your Kubernetes Master server IP using 'kubectl config view' 
 bin/spark-submit --master k8s://https://<K8S-API-SERVER-IP>:443 --deploy-mode cluster --name spark-pi \
 --class org.apache.spark.examples.SparkPi --conf spark.executor.instances=1 --conf spark.app.name=spark-pi \
 --conf spark.kubernetes.driver.docker.image=snappydatainc/spark-driver:v2.2.0-kubernetes-0.5.0 \
 --conf spark.kubernetes.executor.docker.image=snappydatainc/spark-executor:v2.2.0-kubernetes-0.5.0 \
  local:///opt/spark/examples/jars/spark-examples_2.11-2.2.0-k8s-0.5.0.jar 
```

#### Launch the kubernetes dashboard
> You can launch the Kubernetes dashboard (If using GCP you can get to the dashboard from the GCP console) to inspect the 
various deployed objects, associated pods and login the containers to get a better sense for what is going on
```text
# To launch the dashboard, do this ... We use a proxy to access the dashboard locally ...
kubectl proxy

# Goto URL localhost:8001/ui. The page will request a token .... 
# Get the token using ....
kubectl config view | grep -A10 "name: $(kubectl config current-context)" | awk '$1=="access-token:"{print $2}' 
```

#### Stop/delete everything
You can delete everything using 'helm delete'. Note that any changes to notebooks, data, etc will be gone too. 
```text
helm delete --purge spark-all
```

### Quickstart along with History server

> NOTE: Unlike, when using any other cluster manager for Spark (Standalone, Mesos, Yarn) there is no UI so one can monitor 
across all Spark applications or analyze the metrics after a Job completes. We use the Spark History server for this 
purpose. 

Spark applications (drivers, executors) use a shared folder to log events. And, the Spark history server reads from 
this shared folder for its UI. Spark can log using NFS, HDFS or GS (google storage).  


### Setup shared storage 
We need a shared persistent volume to manage state: Spark events from distributed applications(pods) and Notebooks (developed 
using Zeppelin or Jupyter) that you want to preserve/share. Note containers only manage ephemeral state. You need to 
have external persistence so your data survives pod failures.  

#### Steps to setup storage using Google Cloud Storage (GCS)
In this example, we use Google Cloud Storage(GCS) to persist the events generated by Spark applications. You can opt out
of the steps below if you decide to use other schemes like 'hdfs' or 's3' storage.

Using Google cloud utilities (gsutil and gcloud ; should already be setup on your local laptop), we create a GCS bucket 
and associate them with your GCP project.  

    ```
    # Create a bucket using gsutil
    # NOTE: Bucket names have to be globally unique. Pick a unique name if spark-history-server bucket exists.
    gsutil mb -c nearline gs://spark-history-server
    # Specify a account name
    export ACCOUNT_NAME=sparkonk8s-test
    # Change below to specify your Google cloud project name. Use 'gcloud config list' if you don't know. 
    export GCP_PROJECT_ID= your-gcp-project-id
    # Create a service account and generate credentials
    gcloud iam service-accounts create ${ACCOUNT_NAME} --display-name "${ACCOUNT_NAME}"
    gcloud iam service-accounts keys create "${ACCOUNT_NAME}.json" --iam-account "${ACCOUNT_NAME}@${GCP_PROJECT_ID}.iam.gserviceaccount.com"
    # Grant admin rights to the bucket
    gcloud projects add-iam-policy-binding ${GCP_PROJECT_ID} --member "serviceAccount:${ACCOUNT_NAME}@${GCP_PROJECT_ID}.iam.gserviceaccount.com" --role roles/storage.admin
    gsutil iam ch serviceAccount:${ACCOUNT_NAME}@${GCP_PROJECT_ID}.iam.gserviceaccount.com:objectAdmin gs://spark-history-server
    ```

## Lots more to come ....

### High level architecture

![High Level Architecture](k8s-helm-spark-architecture-draw.io.png)

   
   TODO: Can we host our charts into a repo like suggested here (either into a public GCS bucket or through our
    GitHub git-pages itself) - https://github.com/kubernetes/helm/blob/master/docs/chart_repository.md#hosting-chart-repositories. 
    Does it make sense to get listed within the list of incubator charts in k8s?
