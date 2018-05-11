
# SnappyData on Kubernetes

[SnappyData](www.snappydata.io) fuses Apache Spark with an in-memory database to deliver a compute+data engine capable of 
stream processing, transactions, interactive analytics and prediction in a single cluster.

## SnappyData Features

- High concurrency
- Low latency, true interactive speeds
- Native, in-memory indexes
- Data integrity through keys and constraints
- Mutability with support for transactions and snapshot isolation
- Shared in-memory state across Spark applications
- Engine for serving predictions
- Extended support for kubernetes and fully certified for PKS
- HA through distributed consensus based replication
- Support for semi-structured data : Nested data objects, Graph and JSON data
- Framework for Change-data-capture and streaming ETL
- Built-in support for Notebook based development (both Zeppelin and Jupyter)
- Data prep web tool for self-service data cataloging and prep. 
- Enterprise grade security

## Pre-requisites and assumptions:
- We need a running kubernetes or PKS cluster. We only support Kubernetes 1.9 (or higher) and PKS 1.0.0(or higher).
- You must have appropriate role for the K8S service account you are using *details to be added here*
- You must have Helm deployed. Follow the steps given below if Helm is not deployed in your Kubernetes environment

### Setup Helm charts

[Helm](https://github.com/kubernetes/helm/blob/master/README.md) comprises of two parts: a client and a server (Tiller) inside 
the kube-system namespace. Tiller runs inside your Kubernetes cluster, and manages releases (installations) of your charts. 
To install Helm follow the steps [here](https://docs.pivotal.io/runtimes/pks/1-0/configure-tiller-helm.html). The instructions
are applicable for any kubernetes cluster (PKS or GKE or Minikube).


## Quickstart

We use SnappyData Helm chart to deploy SnappyData on Kubernetes. Use `helm install` command to deploy the SnappyData chart

```
# go to charts directory
$ cd charts
# install the chart
$ helm install --name snappydata ./snappydata/
```

The above command will display the Kubernetes objects created by the chart on the console(service, stateful sets etc.).

You may monitor the Kubernetes UI dashboard to check the status of the components as the it takes a few minutes for all 
servers to be online. SnappyData chart provisions volumes dynamically for servers, locator and lead with default 10Gi size. 
These volumes and the data on it will be retained even the chart deployment is deleted.

SnappyData helm chart uses Kubernetes [statefulsets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/) 
to launch locator, lead and servers.

By default SnappyData helm chart will deploy SnappyData cluster consisting of 1 locator, 1 lead and 2 servers and service 
to access SnappyData endpoints.

### Connecting using Snappy shell
 
### Connecting using JDBC driver 

### Using SnappyData with Zeppelin

### Using Spark shell

### Deleting the SnappyData helm chart

```
$ helm delete --purge snappydata
```

## Detailed configuration for SnappyData chart

Users can modify `values.yaml` file to configure the SnappyData chart. Configuration options can be specified in the 
respective attributes of locators, leaders and servers in `values.yaml`

Example: start 4 servers each with heap size of 2048MB
```
servers:
  replicaCount: 2
  ## config options for servers
  conf: "-heap-size=2048m"
```

Users may specify SnappyData [configuration parameters](https://snappydatainc.github.io/snappydata/configuring_cluster/configuring_cluster/#configuration-files) 
in the `servers.conf` attribute. Similalrly users may specify configuration options in the respective section of locators and leaders.

**TODO: Add a table for attributes that can be provided in values.yaml**

## Description of various Kubernetes objects deployed for SnappyData chart

### Statefulsets for Servers, Leaders and Locators

### Description services that expose external endpoints
 
### Persistent Volumes 


## Managing SnappyData in Kubernetes

### Backup and Restore

### Using other Snappy commands 
Document those that are applicable and tested.
 
stats, merge-logs, validate-disk-store, compact-disk-store, compact-all-disk-stores, revoke-missing-disk-store, unblock-disk-store, 
modify-disk-store, show-disk-store-metadata, shut-down-all, print-stacks, run, install-jar, replace-jar, remove-jar

## Future work
Any features that not currently implemented but are planned

## Troubleshooting

