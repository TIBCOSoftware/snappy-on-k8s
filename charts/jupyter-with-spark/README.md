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

### Enabling event logging on Google Cloud Storage (GCS)
To enable event logging for the spark application which you would launch via Jupyter notebooks, follow these steps **before** installing this chart.

   1. Set `sparkEventLog.enableHistoryEvents` to true in values.yaml.
 
   ```python
   sparkEventLog:
     enableHistoryEvents: true
   ```
 
   2. Specify the location of event log directory on your GCS bucket in values.yaml.

   ```python
   # eventsLogDir should point to a URI of GCS bucket where history events will be dumped
   eventLogDir: "gs://spark-history-server-store/"
   ```

   3. Uncomment the line in conf/spark/spark-defaults.conf to specify the keyfile for accessing your GCS bucket.

   **NOTE:** You should only update the name of your keyfile, retaining its path to be /etc/secrets/.

   ```python
   # Uncomment below line and replace sparkonk8s-test.json with the actual name of your keyfile
   # to enable access to Google Cloud Storage.
   spark.hadoop.google.cloud.auth.service.account.json.keyfile   /etc/secrets/bucket-access-key.json
   ```

   4. Place your GCS keyfile under conf/secrets directory of the chart.

Refer to the [readme doc](../spark-hs/README.md) of Spark History Server chart, to launch the history server and monitor event logs of all your Spark applications.

### Setting password for your Jupyter notebook server
This chart by default enables authentication for your Jupyter notebook server. The default password is abc123

Change this password by setting it as the value of `jupyterService.password` attribute in values.yaml.

```python
jupyterService:
  type: LoadBalancer
  port: 8888
  # Set your password to access the notebook server. A default ('abc123') has been set for you.
  # Setting the password to empty string will disable the authentication (not recommended).
  password: 'your-new-password'
```

**NOTE**: if you set a new password via conf/jupyter/jupyter_notebook_config.py, it'll still be overridden by what you specify in values.yaml.

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
      .config("spark.executor.instances", "2")\
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