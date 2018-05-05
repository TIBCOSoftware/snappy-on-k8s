# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# Refer to https://github.com/SnappyDataInc/spark-on-k8s/tree/master/docs/building-images.md#zeppelin-image
# for instructions to build the Docker image.
# This Dockerfile expects spark-on-k8s distribution directory (spark-2.2.0-k8s-0.5.0-bin-2.7.3) and
# the script 'setSparkEnvVars.sh' to be in the same directory where this Dockerfile is kept.

FROM ubuntu:16.04
MAINTAINER Apache Software Foundation <dev@zeppelin.apache.org>

# `Z_VERSION` will be updated by `dev/change_zeppelin_version.sh`
ENV Z_VERSION="0.7.3"
ENV LOG_TAG="[ZEPPELIN_${Z_VERSION}]:" \
    Z_HOME="/zeppelin" \
    LANG=en_US.UTF-8 \
    LC_ALL=en_US.UTF-8

RUN echo "$LOG_TAG update and install basic packages" && \
    apt-get -y update && \
    apt-get install -y locales && \
    locale-gen $LANG && \
    apt-get install -y software-properties-common && \
    apt -y autoclean && \
    apt -y dist-upgrade && \
    apt-get install -y build-essential

RUN echo "$LOG_TAG install tini related packages" && \
    apt-get install -y wget curl grep sed dpkg && \
    TINI_VERSION=`curl https://github.com/krallin/tini/releases/latest | grep -o "/v.*\"" | sed 's:^..\(.*\).$:\1:'` && \
    curl -L "https://github.com/krallin/tini/releases/download/v${TINI_VERSION}/tini_${TINI_VERSION}.deb" > tini.deb && \
    dpkg -i tini.deb && \
    rm tini.deb

ENV JAVA_HOME=/usr/lib/jvm/java-8-openjdk-amd64
RUN echo "$LOG_TAG Install java8" && \
    apt-get -y update && \
    apt-get install -y openjdk-8-jdk && \
    rm -rf /var/lib/apt/lists/*

# Should install conda first before numpy, matploylib since pip and python will be installed by conda
RUN echo "$LOG_TAG Install miniconda2 related packages" && \
    apt-get -y update && \
    apt-get install -y bzip2 ca-certificates \
    libglib2.0-0 libxext6 libsm6 libxrender1 \
    git mercurial subversion && \
    echo 'export PATH=/opt/conda/bin:$PATH' > /etc/profile.d/conda.sh && \
    wget --quiet https://repo.continuum.io/miniconda/Miniconda2-4.3.11-Linux-x86_64.sh -O ~/miniconda.sh && \
    /bin/bash ~/miniconda.sh -b -p /opt/conda && \
    rm ~/miniconda.sh
ENV PATH /opt/conda/bin:$PATH

RUN echo "$LOG_TAG Install python related packages" && \
    apt-get -y update && \
    apt-get install -y python-dev python-pip && \
    apt-get install -y gfortran && \
    # numerical/algebra packages
    apt-get install -y libblas-dev libatlas-dev liblapack-dev && \
    # font, image for matplotlib
    apt-get install -y libpng-dev libfreetype6-dev libxft-dev && \
    # for tkinter
    apt-get install -y python-tk libxml2-dev libxslt-dev zlib1g-dev && \
    pip install numpy && \
    pip install matplotlib

RUN echo "$LOG_TAG Install R related packages" && \
    echo "deb http://cran.rstudio.com/bin/linux/ubuntu xenial/" | tee -a /etc/apt/sources.list && \
    gpg --keyserver keyserver.ubuntu.com --recv-key E084DAB9 && \
    gpg -a --export E084DAB9 | apt-key add - && \
    apt-get -y update && \
    apt-get -y install r-base r-base-dev && \
    R -e "install.packages('knitr', repos='http://cran.us.r-project.org')" && \
    R -e "install.packages('ggplot2', repos='http://cran.us.r-project.org')" && \
    R -e "install.packages('googleVis', repos='http://cran.us.r-project.org')" && \
    R -e "install.packages('data.table', repos='http://cran.us.r-project.org')" && \
    # for devtools, Rcpp
    apt-get -y install libcurl4-gnutls-dev libssl-dev && \
    R -e "install.packages('devtools', repos='http://cran.us.r-project.org')" && \
    R -e "install.packages('Rcpp', repos='http://cran.us.r-project.org')" && \
    Rscript -e "library('devtools'); library('Rcpp'); install_github('ramnathv/rCharts')"

ENV SEARCH_STRING="<name>zeppelin.interpreters<\/name>"
ENV INSERT_STRING="org.apache.zeppelin.interpreter.SnappyDataZeppelinInterpreter,org.apache.zeppelin.interpreter.SnappyDataSqlZeppelinInterpreter,"
ENV LEAD_HOST="localhost"
ENV LEAD_PORT="3768"

RUN echo "$LOG_TAG Download Zeppelin binary and install interpreter for snappydata" && \
    wget -O /tmp/zeppelin-${Z_VERSION}-bin-all.tgz http://archive.apache.org/dist/zeppelin/zeppelin-${Z_VERSION}/zeppelin-${Z_VERSION}-bin-all.tgz && \
    tar -zxvf /tmp/zeppelin-${Z_VERSION}-bin-all.tgz && \
    rm -rf /tmp/zeppelin-${Z_VERSION}-bin-all.tgz && \
    mv /zeppelin-${Z_VERSION}-bin-all ${Z_HOME} && \
    cp ${Z_HOME}/conf/zeppelin-site.xml.template ${Z_HOME}/conf/zeppelin-site.xml && \
    sed -i "/${SEARCH_STRING}/{n;s/<value>/<value>${INSERT_STRING}/}" ${Z_HOME}/conf/zeppelin-site.xml && \
    ${Z_HOME}/bin/install-interpreter.sh --name snappydata --artifact io.snappydata:snappydata-zeppelin:0.7.3 && \
    ${Z_HOME}/bin/zeppelin-daemon.sh start && \
    while ! test -f  ${Z_HOME}/conf/interpreter.json; do sleep 3s; done && \
    ${Z_HOME}/bin/zeppelin-daemon.sh stop && \
    sed -i "/group\": \"snappydata\"/,/isExistingProcess\": false/{s/port\": -1/port\": ${LEAD_PORT}/}" ${Z_HOME}/conf/interpreter.json && \
    sed -i "/group\": \"snappydata\"/,/isExistingProcess\": false/{s/isExistingProcess\": false/isExistingProcess\": snappydatainc_marker/}" ${Z_HOME}/conf/interpreter.json && \
    sed -i "/snappydatainc_marker/a \"host\": \"${LEAD_HOST}\"," ${Z_HOME}/conf/interpreter.json && \
    sed -i "s/snappydatainc_marker/true/" ${Z_HOME}/conf/interpreter.json


RUN echo "$LOG_TAG Cleanup" && \
    apt-get autoclean && \
    apt-get clean

EXPOSE 8080

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

# Copy aws and gcp jars
# COPY aws_gcp_jars/hadoop-aws-2.6.0.jar /opt/spark/jars
# COPY aws_gcp_jars/aws-java-sdk-1.7.4.jar /opt/spark/jars
# COPY aws_gcp_jars/gcs-connector-latest-hadoop2.jar /opt/spark/jars

COPY setSparkEnvVars.sh /opt/

ENV SPARK_HOME /opt/spark

COPY spark-2.2.0-k8s-0.5.0-bin-2.7.3/examples /opt/spark/examples

CMD SPARK_CLASSPATH="${SPARK_HOME}/jars/*" && \
    env | grep SPARK_JAVA_OPT_ | sed 's/[^=]*=\(.*\)/\1/g' > /tmp/java_opts.txt && \
    readarray -t SPARK_DRIVER_JAVA_OPTS < /tmp/java_opts.txt && \
    if ! [ -z ${SPARK_MOUNTED_CLASSPATH+x} ]; then SPARK_CLASSPATH="$SPARK_MOUNTED_CLASSPATH:$SPARK_CLASSPATH"; fi && \
    if ! [ -z ${SPARK_SUBMIT_EXTRA_CLASSPATH+x} ]; then SPARK_CLASSPATH="$SPARK_SUBMIT_EXTRA_CLASSPATH:$SPARK_CLASSPATH"; fi && \
    if ! [ -z ${SPARK_EXTRA_CLASSPATH+x} ]; then SPARK_CLASSPATH="$SPARK_EXTRA_CLASSPATH:$SPARK_CLASSPATH"; fi && \
    if ! [ -z ${SPARK_MOUNTED_FILES_DIR+x} ]; then cp -R "$SPARK_MOUNTED_FILES_DIR/." .; fi && \
    if ! [ -z ${SPARK_MOUNTED_FILES_FROM_SECRET_DIR} ]; then cp -R "$SPARK_MOUNTED_FILES_FROM_SECRET_DIR/." .; fi && \
    ${JAVA_HOME}/bin/java "${SPARK_DRIVER_JAVA_OPTS[@]}" -cp "$SPARK_CLASSPATH" -Xms$SPARK_DRIVER_MEMORY -Xmx$SPARK_DRIVER_MEMORY -Dspark.driver.bindAddress=$SPARK_DRIVER_BIND_ADDRESS $SPARK_DRIVER_CLASS $SPARK_DRIVER_ARGS

#ENV MASTER k8s://$KUBERNETES_SERVICE_HOST:$KUBERNETES_SERVICE_PORT
#ENV SPARK_SUBMIT_OPTIONS "--kubernetes-namespace default --conf spark.kubernetes.driver.pod.name=$HOSTNAME --conf spark.kubernetes.driver.docker.image=shirishd/spark-driver:v2.2.0 --conf spark.kubernetes.executor.docker.image=shirishd/spark-executor:v2.2.0"

#CMD ["/opt/setSparkEnvVars.sh"]
CMD ["bin/bash", "-c", "source", "/opt/setSparkEnvVars.sh"]
#RUN /bin/bash -c "source /opt/setSparkEnvVars.sh"
####### End changes for Spark-on-k8s ##########################

ENTRYPOINT [ "/usr/bin/tini", "--" ]
WORKDIR ${Z_HOME}
CMD ["bin/zeppelin.sh"]
