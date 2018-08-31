
# Building and publishing Docker images

## Prerequisites

You should have Docker installed on your local setup from where you would want to build and publish the Docker images.

Refer to [this page](https://docs.docker.com/install) to get information about installing Docker.

## Spark Images

The binaries used to build the Spark images are based on the [spark-on-k8s](https://github.com/apache-spark-on-k8s/spark) project, with few additional changes.
These have been committed into a clone of branch-2.2-kubernetes branch in above repository, and is available as a branch in SnappyData's fork of Apache Spark.

Get the latest branch:

```bash
$ git clone https://github.com/SnappyDataInc/spark.git -b snappy/branch-2.2-kubernetes
```

Go to the checkout directory and build the project using [maven](https://maven.apache.org/install.html).
Also, package the build into a tarball, which will be needed when building the Docker images for Jupyter and Apache Zeppelin.

```bash
$ ./build/mvn -Pkubernetes -DskipTests clean package
$ ./dev/make-distribution.sh --name 2.7.3 --pip --tgz -Phadoop-2.7 -Phive -Phive-thriftserver -Pkubernetes
```

Now that the binaries are built, you also need to download and place following jars into the directories
assembly/target/scala-2.11/jars and dist/jars of your checkout.
These are needed for enabling access to Google Cloud Storage and AWS S3 buckets, which your Spark applications may need.

1. [aws-java-sdk-1.7.4.jar](http://central.maven.org/maven2/com/amazonaws/aws-java-sdk/1.7.4/aws-java-sdk-1.7.4.jar)
2. [hadoop-aws-2.7.3.jar](http://central.maven.org/maven2/org/apache/hadoop/hadoop-aws/2.7.3/hadoop-aws-2.7.3.jar)
3. [gcs-connector-latest-hadoop2.jar](https://storage.googleapis.com/hadoop-lib/gcs/gcs-connector-latest-hadoop2.jar)

Now build and publish the Docker images to your DockerHub account. It may take several minutes depending upon your network speed.

```bash
$ ./sbin/build-push-docker-images.sh -r <your-docker-repo-name> -t <image-tag> build
```

Make sure you are logged in to your Docker Hub account before publishing the images:

```bash
$ docker login
Login with your Docker ID to push and pull images from Docker Hub. If you don't have a Docker ID, head over to https://hub.docker.com to create one.
Username: <your-account-username>
Password: <password>
$ ./sbin/build-push-docker-images.sh -r <your-docker-repo-name> -t <image-tag> push
```

## Jupyter Image

This image will contain the Spark binaries you built above apart from the dependencies needed for Jupyter Notebook server.

Extract the Spark tarball generated above into a directory where you have copied the [Dockerfile for Jupyter image](../dockerfiles/jupyter/Dockerfile) to.

Make sure that the third party jars needed to access GCS and AWS S3 are copied to the jars directory of the extracted tarball.

Build and publish the Jupyter image:

```bash
$ docker build -t <your-docker-repo-name>/jupyter-notebook:<image-tag> -f Dockerfile .
$ docker push <your-docker-repo-name>/jupyter-notebook:<image-tag>
```

For example:
```bash
$ docker build -t snappydatainc/jupyter-notebook:5.2.2-spark-v2.2.0-kubernetes-0.5.1 -f Dockerfile .
$ docker push snappydatainc/jupyter-notebook:5.2.2-spark-v2.2.0-kubernetes-0.5.1
```

## Zeppelin Image

This image will contain the Spark binaries built earlier apart from the dependencies needed for launching Apache Zeppelin server.

Extract the Spark tarball generated above into a directory where you have copied the [Dockerfile for Zeppelin image](../dockerfiles/zeppelin/Dockerfile) to.
Also, copy the script [setSparkEnvVars.sh](../dockerfiles/zeppelin/setSparkEnvVars.sh) to the same location.

Make sure that the third party jars needed to access GCS and AWS S3 are copied to the jars directory of the extracted tarball.

Build and publish the Zeppelin image.

```bash
$ docker build -t <your-docker-repo-name>/zeppelin:<image-tag> -f Dockerfile .
$ docker push <your-docker-repo-name>/zeppelin:<image-tag>
```

For example:
```bash
$ docker build -t snappydatainc/zeppelin:0.7.3-spark-v2.2.0-kubernetes-0.5.1 -f Dockerfile .
$ docker push snappydatainc/zeppelin:0.7.3-spark-v2.2.0-kubernetes-0.5.1
```

## SnappyData Image

The SnappyData Docker image available on DockerHub is built using the OSS version of the product. Docker image with
SnappyData Enterprise bits will be available soon.

Currently, some manual steps are needed to build this image which will be automated later.

- Download the Snappydata OSS tarball of the version you need and available on
[GitHub releases page](https://github.com/snappydatainc/snappydata/releases) and extract its content into a directory.

- Copy the Dockerfile and start script required for SnappyData image
[from this branch](https://github.com/SnappyDataInc/snappy-cloud-tools/blob/SNAP-2280/docker) into the extracted
SnappyData directory.

- Copy the [SnappyData interpreter jar](https://github.com/SnappyDataInc/zeppelin-interpreter/releases) for
Apache Zeppelin into the jars directory.

- Optionally, one can also add the third party jar needed to access GCS to the jars directory. The libraries to access
AWS S3 and HDFS are already included.

- Switch to the extracted directory to build and publish the SnappyData image using following commands.

    ```bash
    $ cd <extracted-snappydata-directory>
    $ docker build -t <your-docker-repo-name>/snappydata:<image-tag> -f Dockerfile .
    $ docker push <your-docker-repo-name>/snappydata:<image-tag>
    ```

    For example:
    ```bash
    $ docker build -t snappydatainc/snappydata:1.0.1 -f Dockerfile .
    $ docker push snappydatainc/snappydata:1.0.1
    ```
