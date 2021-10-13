ARG SPARK_VERSION=2.4.3-hadoop2.7
FROM bde2020/spark-base:${SPARK_VERSION}

LABEL maintainer="lilinghai"
ARG TISPARK_VERSION=2.3.3
RUN echo TISPARK_VERSION: ${TISPARK_VERSION}
ENV SPARK_MASTER_PORT 7077
ENV SPARK_MASTER_WEBUI_PORT 8080
ENV SPARK_WORKER_WEBUI_PORT 8081

RUN wget https://github.com/pingcap/tispark/releases/download/v${TISPARK_VERSION}/tispark-assembly-${TISPARK_VERSION}.jar \
      && wget https://github.com/pingcap/tispark/releases/download/v${TISPARK_VERSION}/tispark-assembly-${TISPARK_VERSION}.jar.md5 \
	  && echo "  tispark-assembly-${TISPARK_VERSION}.jar" >> tispark-assembly-${TISPARK_VERSION}.jar.md5 \
	  && md5sum -c tispark-assembly-${TISPARK_VERSION}.jar.md5 \
	  && mv tispark-assembly-${TISPARK_VERSION}.jar /spark/jars \
	  && rm tispark-assembly-${TISPARK_VERSION}.jar.md5
RUN apk --update add coreutils procps rsync openssh

COPY start_script.sh /spark/
WORKDIR /spark
EXPOSE 8080 7077 10000 4040 8081

CMD ["/bin/bash", "start_script.sh"]
