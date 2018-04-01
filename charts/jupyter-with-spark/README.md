This chart is work-in-progress

# A Helm chart for Jupyter notebook environment with Spark-on-k8s binaries

## Chart Details

* This chart launches Jupyter pod and exposes Jupyter as a service on port 8888.
  * You can provide configurations for Jupyter notebook by placing your jupiter_notebook_config.py under conf/jupyter directory of this chart.
* Also, you can run your PySpark applications from the Jupyter notebook environment by explicitly creating SparkContext.
  * This launches your Spark cluster in in-cluster client mode. The in-cluster client mode means that the submission environment is within a kubernetes cluster.
  * You can place your spark-defaults.conf under conf/spark/ directory of this chart which will be used to configure the Spark cluster.
  * A basic spark-defaults.conf is provided which specifies the docker images from SnappyData, Inc and few other defaults.
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

3. Jupyter notebook server URL can be found by running the commands displayed after running `helm install ...` above. 

```bash
  $ export JUPYTER_SERVICE_IP=$(kubectl get svc --namespace default jupyter-jupyter-notebook-ui -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
  $ echo http://$JUPYTER_SERVICE_IP:8888
```
  *NOTE*: It may take a few minutes before the LoadBalancer's IP is available. You can watch the status of by running `kubectl get svc -w upyter-jupyter-notebook-ui`

# Running Spark application via notebook 
Create a new notebook based on Python 2 and launch a Spark cluster by creating SparkContext.

```python
spark = SparkSession.builder.config("spark.app.name", "spark-pi")\
      .config("spark.executor.instances", "2")\
      .getOrCreate()
```
Note that, other spark configurations like `spark.master`, `spark.submit.deployMode`, docker image names for driver and executor pods, kubernetes namespace are already set for you.
See spark-defaults.conf under conf/spark/ directory of the chart for more details. You can also override the conf properties or set new ones in the notebook.

```python
spark = SparkSession.builder.config("spark.app.name", "spark-pi")\
      .config("spark.kubernetes.namespace", "dev-namespace")\
      .config("spark.executor.instances", "5")\
      .getOrCreate()
```

You can change the default Spark UI port from 4040 to something else in the notebook, either by updating it in the [values.yml](values.yaml) under `sparkWebUI` attribute 
or using `--set sparkWebUI.port=<new-port>` when invoking `helm install ...`.

Once SparkContext is created in the notebook, Spark UI URL can be obtained by using commands displayed by `helm install`:

```bash
   $ export SPARK_UI_SERVICE_IP=$(kubectl get svc --namespace default jupyter-jupyter-spark-web-ui -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
   $ echo http://$SPARK_UI_SERVICE_IP:4040
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
      .config("spark.app.name", "spark-pi")\
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