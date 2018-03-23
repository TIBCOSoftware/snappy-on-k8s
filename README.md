# Multi cloud Spark application service on PKS

## Introduction
[Kubernetes](kubernetes.io) is an open source project designed specifically for container orchestration. 
Kubernetes offers a number of key [features](https://kubernetes.io/docs/concepts/overview/what-is-kubernetes/), 
including multiple storage APIs, container health checks, manual or automatic scaling, rolling upgrades and 
service discovery.
[Pivotal Container Service - PKS](https://pivotal.io/platform/pivotal-container-service) is a solution to manage 
Kubernetes clusters across private and public clouds. It leverages [BOSH]() to offer a uniform way to instantiate, 
deploy, and manage highly available Kubernetes clusters on a cloud platform like GCP, VMWare vSphere or AWS. 

This project provides a streamlined way of deploying, scaling and managing Spark applications. Besides support for Spark 
batch jobs, this project leverages [Helm charts](https://helm.sh/) to allow deployment of 
common Spark application recipes - using [Apache Zeppelin](https://zeppelin.apache.org/) and/or [Jupyter](jupyter.org) for interactive, collaborative 
workloads. 
Helm is a package manager for kubernetes and the most productive way to find, install, share, upgrade and use even the 
most complex kubernetes applications. So, for instance, with a single command you can deploy additional components like 
HDFS or Elastic search for our Spark applications.  

This project is a collaborative effort between SnappyData and Pivotal. 



## Features
- Full support for Spark 2.2 applications running on PKS 1.x on both Google cloud and on-prem VMWare cloud environments.
The project leverages [spark-on-k8s]() work.
- Deploy batch spark jobs using built-in support for kubernetes
- Helm chart to deploy Zeppelin, centralized logging, monitoring across apps (using History server)
- Helm chart to deploy Jupyter,  centralized logging, monitoring across apps (using History server)
- Spark applications can be Java, Scala or Python
- Spark applications can dynamically scale

We use Helm charts to abstract the developer from having to understand kubernetes concepts and simply focus on 
configuration that matters. Think recipes that come with sensible defaults for common Spark workloads on PKS. 

<Provide a graphic that depicts our Helm charts based architecture .. high-level )

We showcase the use of cloud storage (e.g. Google Cloud Storage) to manage logs/events but show how the use persistent 
volumes within the charts make the architecture portable across clouds. 

## Pre-requisites and assumptions:
- A running kubernetes or PKS cluster. We only support Kubernetes 1.9 (or higher) and PKS 1.0.1(or higher). If you don't
have access to a Kubernetes/PKS cluster, see pointers below. 
- Install kubectl on your local development machine and configure access to the kubernetes/PKS cluster. See instructions for 
kubectl [here](https://kubernetes.io/docs/tasks/tools/install-kubectl/). If you are using Google cloud, you will find 
instuctions for setting up Google Cloud SDK ('gcloud') along with kubectl 
[here](https://kubernetes.io/docs/tasks/tools/install-kubectl/).
- You must have appropriate permissions to list, create, edit and delete pods in your cluster. You can verify that you 
can list these resources by running kubectl auth can-i <list|create|edit|delete> pods.
- The service account credentials used by the driver pods must be allowed to create pods, services and configmaps.
- You must have Kubernetes DNS configured in your cluster.


### Getting access to a PKS or Kubernetes cluster
- PKS: You will need to install the PKS CLI and setup a PKS cluster on either a vSphere cloud or Google cloud. 
The instructions for this is [here](https://docs.pivotal.io/runtimes/pks/1-0/prerequisites.html)
- Google Cloud Platform: Login to your Google account and goto the [Cloud console](console.cloud.google.com) to launch a 
GKE cluster
- Minikube: If either of the above options is difficult, you may setup a test cluster on your local machine using 
[minikube](https://kubernetes.io/docs/getting-started-guides/minikube/). We recommend using the latest release of minikube 
with the DNS addon enabled.
- If using Minikube, be aware that the default minikube configuration is not enough for running Spark applications. 
We recommend 3 CPUs and 4g of memory to be able to start a simple Spark application with a single executor.


### Setup Helm charts

[Helm](https://github.com/kubernetes/helm/blob/master/README.md) comprises of two parts: a client and a server (Tiller) inside 
the kube-system namespace. Tiller runs inside your Kubernetes cluster, and manages releases (installations) of your charts. 
To be able to do this, Tiller needs access to the Kubernetes API. By default, RBAC policies will not allow Tiller to carry 
out these operations. 

To install Helm follow the steps [here](https://docs.pivotal.io/runtimes/pks/1-0/configure-tiller-helm.html). The instructions
are applicable for any kubernetes cluster (PKS or GKE or Minikube).


## Lots more to come ....

### High level architecture

![High Level Architecture](k8s-helm-spark-architecture-draw.io.png)

   
   TODO: Can we host our charts into a repo like suggested here (either into a public GCS bucket or through our
    GitHub git-pages itself) - https://github.com/kubernetes/helm/blob/master/docs/chart_repository.md#hosting-chart-repositories. 
    Does it make sense to get listed within the list of incubator charts in k8s?
