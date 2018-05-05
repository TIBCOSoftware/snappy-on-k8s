# Copyright (c) Jupyter Development Team.
# Distributed under the terms of the Modified BSD License.

# Refer to https://github.com/SnappyDataInc/spark-on-k8s/tree/master/docs/building-images.md#jupyter-image
# for instructions to build the Docker image.
# This Dockerfile should be present in the same directory where spark-on-k8s distribution
# directory (spark-2.2.0-k8s-0.5.0-bin-2.7.3) is kept.

FROM jupyter/scipy-notebook

USER root

# Copied from pyspark notebook- start
# Spark dependencies
ENV APACHE_SPARK_VERSION 2.2.0
ENV HADOOP_VERSION 2.7

RUN apt-get -y update && \
    apt-get install --no-install-recommends -y openjdk-8-jre-headless ca-certificates-java && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Copied from pyspark notebook- end

####### Begin changes for Spark-on-k8s #################

RUN mkdir -p /opt/spark && \
    mkdir -p /opt/spark/work-dir \
    touch /opt/spark/RELEASE && \
    rm -f /bin/sh && \
    ln -sv /bin/bash /bin/sh && \
    chgrp root /etc/passwd && chmod ug+rw /etc/passwd

COPY spark-2.2.0-k8s-0.5.0-bin-2.7.3/jars /opt/spark/jars
COPY spark-2.2.0-k8s-0.5.0-bin-2.7.3/bin /opt/spark/bin
COPY spark-2.2.0-k8s-0.5.0-bin-2.7.3/sbin /opt/spark/sbin
COPY spark-2.2.0-k8s-0.5.0-bin-2.7.3/conf /opt/spark/conf
COPY spark-2.2.0-k8s-0.5.0-bin-2.7.3/dockerfiles/spark-base/entrypoint.sh /opt/

ADD spark-2.2.0-k8s-0.5.0-bin-2.7.3/examples /opt/spark/examples
ADD spark-2.2.0-k8s-0.5.0-bin-2.7.3/python /opt/spark/python

# Copy aws and gcp jars
# COPY aws_gcp_jars/hadoop-aws-2.7.3.jar /opt/spark/jars
# COPY aws_gcp_jars/aws-java-sdk-1.7.4.jar /opt/spark/jars
# COPY aws_gcp_jars/gcs-connector-latest-hadoop2.jar /opt/spark/jars

ENV SPARK_HOME /opt/spark

ENV PYTHON_VERSION 2.7.13
ENV PYSPARK_PYTHON python
ENV PYSPARK_DRIVER_PYTHON python
ENV PYTHONPATH ${SPARK_HOME}/python/:${SPARK_HOME}/python/lib/py4j-0.10.4-src.zip:${PYTHONPATH}

CMD SPARK_CLASSPATH="${SPARK_HOME}/jars/*" && \
    env | grep SPARK_JAVA_OPT_ | sed 's/[^=]*=\(.*\)/\1/g' > /tmp/java_opts.txt && \
    readarray -t SPARK_DRIVER_JAVA_OPTS < /tmp/java_opts.txt && \
    if ! [ -z ${SPARK_MOUNTED_CLASSPATH+x} ]; then SPARK_CLASSPATH="$SPARK_MOUNTED_CLASSPATH:$SPARK_CLASSPATH"; fi && \
    if ! [ -z ${SPARK_SUBMIT_EXTRA_CLASSPATH+x} ]; then SPARK_CLASSPATH="$SPARK_SUBMIT_EXTRA_CLASSPATH:$SPARK_CLASSPATH"; fi && \
    if ! [ -z ${SPARK_EXTRA_CLASSPATH+x} ]; then SPARK_CLASSPATH="$SPARK_EXTRA_CLASSPATH:$SPARK_CLASSPATH"; fi && \
    if ! [ -z ${SPARK_MOUNTED_FILES_DIR+x} ]; then cp -R "$SPARK_MOUNTED_FILES_DIR/." .; fi && \
    if ! [ -z ${SPARK_MOUNTED_FILES_FROM_SECRET_DIR+x} ]; then cp -R "$SPARK_MOUNTED_FILES_FROM_SECRET_DIR/." .; fi && \
    ${JAVA_HOME}/bin/java "${SPARK_DRIVER_JAVA_OPTS[@]}" -cp "$SPARK_CLASSPATH" -Xms$SPARK_DRIVER_MEMORY -Xmx$SPARK_DRIVER_MEMORY -Dspark.driver.bindAddress=$SPARK_DRIVER_BIND_ADDRESS $SPARK_DRIVER_CLASS $PYSPARK_PRIMARY $PYSPARK_FILES $SPARK_DRIVER_ARGS


# Copied from pyspark notebook- start
# Spark config
ENV PYTHONPATH $SPARK_HOME/python:$SPARK_HOME/python/lib/py4j-0.10.4-src.zip
ENV SPARK_OPTS --driver-java-options=-Xms1024M --driver-java-options=-Xmx4096M --driver-java-options=-Dlog4j.logLevel=info
# Copied from pyspark notebook- end

RUN chown -R $NB_USER:users /opt/spark

####### End changes for Spark-on-k8s ##########################

USER $NB_USER

# Added to anble python2 notebooks
RUN conda create --quiet --yes \
     -n ipykernel_py2 python=2 ipykernel && \
     source activate ipykernel_py2 && \
     python -m ipykernel install --user


RUN  source activate ipykernel_py2 && \
     conda install --yes \
     matplotlib \
     scipy \
     numpy \
     pandas \ 
     nltk \
     tensorflow && \
     source activate ipykernel_py2 && \
     pip install \
     sklearn \
     wordcloud \
     treeinterpreter

