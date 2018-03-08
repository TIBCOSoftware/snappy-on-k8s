This chart is work-in-progress

# A Helm chart for Jupyter notebook environment with Spark-on-k8s binaries

## Chart Details

* This chart launches Jupyter pod and exposes Jupyter as a service on port 8888.
  * You can provide configurations for Jupyter notebook by placing your jupiter_notebook_config.py under conf/jupyter directory of this chart.
* Also, you can run your PySpark applications from the Jupyter notebook environment by explicitly creating SparkContext.
  * This launches your Spark cluster in in-cluster client mode. The in-cluster client mode means that the submission environment is within a kubernetes cluster.
  * You can place your spark-defaults.conf under conf/spark/ directory of this chart which will be used to configure the Spark cluster.
  * A basic spark-defaults.conf is provided which specifies the docker images from SnappyData Inc.
  * This chart uses docker images built with spark-on-k8s binaries after applying patch for [PR#456](https://github.com/apache-spark-on-k8s/spark/pull/456) as the changes for it are not yet merged in [spark-on-k8s](https://github.com/apache-spark-on-k8s/spark) project.

## Installing the Chart
1. Make sure that Helm is setup in your kubernetes environment. You may refer to these [steps](https://docs.bitnami.com/kubernetes/get-started-kubernetes/#step-4-install-helm-and-tiller).
2. Clone the chart from GitHub and install it. 

```bash
  $ git clone https://github.com/SnappyDataInc/spark-on-k8s
  $ helm install --name jupyter ./spark-on-k8s/charts/jupyter-with-spark/
```

  The above command will deploy the helm chart and will display instructions to access Jupyter service and Spark UI.
  Note that the Spark UI will be available only after you start a SparkContext in Jupyter notebook.

3. Jupyter server URL can be found by running following commands:

```bash
  $ export JUPYTER_SERVICE_IP=$(kubectl get svc --namespace default jupyter-notebook-jupyter-with-spark -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
  $ echo http://$JUPYTER_SERVICE_IP:{{ .Values.jupyterService.port }}
```
  *NOTE*: It may take a few minutes before the LoadBalancer's IP is available. You can watch the status of by running `kubectl get svc -w jupyter-notebook-jupyter-with-spark`

# Running Spark application via notebook 
While creating the SparkContext, in addition to the usual conf settings required to launch Spark on kubernetes, user must
also set the "spark.master" in SparkConf as given below, to be able to launch the driver and executors on the kubernetes cluster.

```python
conf.set("spark.master", "k8s://https://{}:{}".format(
        os.getenv("KUBERNETES_SERVICE_HOST"),
        os.getenv("KUBERNETES_SERVICE_PORT")
      ))
```

If you want to change the default Spark UI port from 4040 to something else in the notebook, first update it in the [values.yml](values.yaml) under `sparkWebUI` attribute before installing the chart.
Otherwise, you may not be able access the Spark UI from outside the kubernetes cluster.

Once SparkContext is created in the notebook, Spark UI URL can be obtained by using following commands:

```bash
   $ export SPARK_UI_SERVICE_IP=$(kubectl get svc --namespace default spark-web-ui -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
   $ echo http://$SPARK_UI_SERVICE_IP:{{ .Values.sparkWebUI.port }}
```

A sample Spark job to calculate the value of Pi is given below.

```python
from __future__ import print_function

import sys
from random import random
from operator import add
import os
from pyspark.sql import SparkSession


spark = SparkSession\
      .builder\
      .appName("PythonPi")\
      .config("spark.submit.deployMode", "client")\
      .config("spark.master", "k8s://https://{}:{}".format(
        os.getenv("KUBERNETES_SERVICE_HOST"),
        os.getenv("KUBERNETES_SERVICE_PORT")
      ))\
      .config("spark.app.name", "spark-pi")\
      .config("spark.kubernetes.driver.docker.image", "snappydatainc/spark-driver-py:v2.2.0-7b8c9f5")\
      .config("spark.kubernetes.executor.docker.image", "snappydatainc/spark-executor-py:v2.2.0-7b8c9f5")\
      .config("spark.kubernetes.initcontainer.docker.image", "snappydatainc/spark-init:v2.2.0-7b8c9f5")\
      .config("spark.kubernetes.docker.image.pullPolicy", "Always")\
      .config("spark.kubernetes.driver.pod.name", os.getenv("HOSTNAME"))\
      .config("spark.kubernetes.namespace", "default")\
      .config("spark.executor.instances", "2")\
      .getOrCreate()
partitions = 2
n = 100000 * partitions

def f(_):
    x = random() * 2 - 1
    y = random() * 2 - 1
    return 1 if x ** 2 + y ** 2 <= 1 else 0

count = spark.sparkContext.parallelize(range(1, n + 1), partitions).map(f).reduce(add)
print("Pi is roughly %f" % (4.0 * count / n))
```