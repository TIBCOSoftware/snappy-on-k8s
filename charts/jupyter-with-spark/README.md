
# A Helm chart for Jupyter notebook enabled Spark cluster

## Chart Details

* This chart launches Jupyter notebook server pod and exposes Jupyter notebook as a service on port 8888.
  * You can provide configurations for Jupyter notebook by placing your jupiter_notebook_config.py under [conf/jupyter](conf/jupyter) directory of this chart.
* Also, you can run your PySpark applications from the Jupyter notebook environment by explicitly creating SparkContext.
  * This launches your Spark cluster in in-cluster client mode. The in-cluster client mode means that the submission environment is within a kubernetes cluster.
  * You can configure your Spark cluster by providing your custom spark-defaults.conf under [conf/spark/](conf/spark) directory of this chart.
  * A basic [spark-defaults.conf](conf/spark/spark-defaults.conf) is provided which specifies the docker images from SnappyData, Inc. and few other defaults.
  * This chart uses docker images built with spark-on-k8s binaries after applying patch for [PR#456](https://github.com/apache-spark-on-k8s/spark/pull/456) as the changes for it are not yet merged in [spark-on-k8s](https://github.com/apache-spark-on-k8s/spark) project.

## Installing the Chart
1. Make sure that Helm is setup in your kubernetes environment. You may refer to these [steps](https://docs.bitnami.com/kubernetes/get-started-kubernetes/#step-4-install-helm-and-tiller).
2. Clone the chart from GitHub and install it. 

```bash
  $ git clone https://github.com/SnappyDataInc/spark-on-k8s
  $ helm install --name jupyter --namespace spark ./spark-on-k8s/charts/jupyter-with-spark/
```

  The above command will deploy the helm chart and will display instructions to access Jupyter service and Spark UI.
  Note that the Spark UI will be available only after you start a SparkContext in Jupyter notebook.

3. Jupyter notebook server URL can be found by running the commands displayed after running `helm install ...` above. 

```bash
  $ export JUPYTER_SERVICE_IP=$(kubectl get svc --namespace default jupyter-jupyter-spark -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
  $ echo http://$JUPYTER_SERVICE_IP:8888
```
  *NOTE*: It may take a few minutes before the LoadBalancer's IP is available. You can watch the status of by running `kubectl get svc -w jupyter-jupyter-spark`

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

   3. Specify name of your keyfile for accessing your GCS bucket in conf/spark/spark-defaults.conf

   **NOTE:** You should only update the name of your keyfile, retaining its path to be /etc/secrets/.

   ```python
   # Replace bucket-access-key.json with the actual name of your keyfile
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

## Running Spark application via notebook 
Create a new notebook based on Python 2 and launch a Spark cluster by creating SparkContext.

```python
spark = SparkSession.builder.config("spark.app.name", "spark-pi")\
      .config("spark.executor.instances", "2")\
      .getOrCreate()
```
Note that, other spark configurations like `spark.master`, `spark.submit.deployMode`, docker image names for driver and executor pods, kubernetes namespace are already set for you.
See [spark-defaults.conf](conf/spark/spark-defaults.conf) under conf/spark/ directory of the chart for more details. You can also override the conf properties or set new ones in the notebook.

```python
spark = SparkSession.builder.config("spark.app.name", "spark-pi")\
      .config("spark.kubernetes.namespace", "dev-namespace")\
      .config("spark.executor.instances", "2")\
      .getOrCreate()
```

You can change the default Spark UI port from 4040 to something else in the notebook, either by updating it in the [values.yml](values.yaml) under `jupyterService` attribute 
or while installing the chart as `helm install --name jupyter --set jupyterService.sparkUIPort=5050 ./spark-on-k8s/charts/jupyter-with-spark/`.

Once you create SparkContext in your Jupyter notebook, you can access the Spark UI via the URL obtained by running below commands (displayed by `helm install`):

```bash
   $ export SPARK_UI_SERVICE_IP=$(kubectl get svc --namespace default jupyter-jupyter-spark -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
   $ echo http://$SPARK_UI_SERVICE_IP:4040
```

A sample Spark job to calculate the value of Pi is given below which you can try running in your Jupyter notebook.

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

## Configuration Properties List
   The following table lists the configuration parameters available for this chart. You can modify these in [values.yaml](values.yaml).
   
   | Parameter               | Description                           | Default                                     |
   | ----------------------- | ------------------------------------- | ------------------------------------------- |
   | `image.repository`      |  Docker repo/name for the image       |     `SnappyDataInc/jupyter-notebook`        |
   | `image.tag`             |  Tag for the Docker image             |     `5.2.2-spark-v2.2.0-kubernetes-0.5.1`   | 
   | `image.pullPolicy`      |  Pull policy for the image            |     `IfNotPresent`                          |
   | `jupyterService.type`   |  K8S service type for Jupyter server  |     `LoadBalancer`                          |
   | `jupyterService.jupyterPort`   |  Port for Jupyter notebook service    |     `8080`                           |
   | `jupyterService.sparkUIPort`   |  Port for Spark service               |     `4040`                           |
   | `serviceAccount`        |  Service account used to deploy Jupyter and run Spark jobs |     `default`          |
   | `sparkEventLog.enableHistoryEvents` | Whether to enable Spark event logging, required by History server |  `false` |
   | `sparkEventLog.eventLogDir` | URL of the GCS bucket where Spark event logs will be written | |
   | `mountSecrets`          | If true, files in 'secrets' directory will be mounted on path /etc/secrets | `false` |   |
   | `persistence.enabled`   | Whether to mount a persistent volume  | `true` |
   | `persistence.existingClaim` | An existing PVC to be used while mounting a PV. If this is not specified a dynamic PV will be created and a PVC will be generated for it |  |
   | `persistence.storageClass`  | Storage class to be used for creating a PVC, if an `existingClaim` is not specified. If unspecified default will be chosen. Ignored if `persistence.existingClaim` is specified | |
   | `persistence.accessMode`    | Access mode for the dynamically generated PV and its PVC. Ignored if `persistence.existingClaim` is specified | `ReadWriteOnce`|
   | `persistence.size`      | Size of the dynamically generated PV and its PVC. Ignored if `persistence.existingClaim` is specified | 8Gi |
   | `resources`             | CPU and Memory resources for the Jupyter pod  | |
   | `global.umbrellaChart`  | Internal attribute. Do not modify.    | `false` | 
   

   You can also set these attributes with the helm install command.
   
   For example: 
   ```
   # set an attribute while using helm install command
   helm install --name jupyter --namespace spark --set serviceAccount=spark ./jupyter-with-spark
   ```
