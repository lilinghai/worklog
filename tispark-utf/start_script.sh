#!/bin/bash

export SPARK_MASTER_HOST=`hostname`
#printf "spark.driver.host ${SPARK_MASTER_HOST} \nspark.driver.port 8089 \nspark.driver.blockManager.port 8090 \nspark.tispark.write.allow_spark_sql   True \nspark.sql.extensions   org.apache.spark.sql.TiExtensions \nspark.tispark.pd.addresses ${PD_ADDR}" > /spark/conf/spark-defaults.conf
printf "spark.tispark.write.allow_spark_sql   True \nspark.sql.crossJoin.enabled   True \nspark.sql.extensions   org.apache.spark.sql.TiExtensions \nspark.tispark.pd.addresses ${PD_ADDR}" > /spark/conf/spark-defaults.conf
chown 1000:1000 conf/spark-defaults.conf

. "/spark/sbin/spark-config.sh"

. "/spark/bin/load-spark-env.sh"

export SPARK_HOME=/spark
export SPARK_CONF_DIR=/spark/conf

if [ "$DEPLOY" = "thrift" ]
then
  #ln -sf /dev/stdout $SPARK_MASTER_LOG/spark-thrift.out
  #bin/spark-submit --class org.apache.spark.sql.hive.thriftserver.HiveThriftServer2 --name Thrift JDBC/ODBC Server --master spark://${SPARK_MASTER_HOST}:7077 >> $SPARK_MASTER_LOG/spark-thrift.out
  sbin/start-thriftserver.sh --master spark://${SPARK_MASTER_HOST}:7077
  tail -f /spark/logs/*
elif [ "$DEPLOY" = "worker" ]
then
  #ln -sf /dev/stdout $SPARK_MASTER_LOG/spark-thrift.out
  #bin/spark-submit --class org.apache.spark.sql.hive.thriftserver.HiveThriftServer2 --name Thrift JDBC/ODBC Server --master spark://${SPARK_MASTER_HOST}:7077 >> $SPARK_MASTER_LOG/spark-thrift.out
  sbin/start-slave.sh spark://${SPARK_MASTER_HOST}:7077 --webui-port $SPARK_WORKER_WEBUI_PORT
  tail -f /spark/logs/*
else
  #ln -sf /dev/stdout $SPARK_MASTER_LOG/spark-master.out
  #bin/spark-class org.apache.spark.deploy.master.Master \
  #  --ip $SPARK_MASTER_HOST --port $SPARK_MASTER_PORT --webui-port $SPARK_MASTER_WEBUI_PORT >> $SPARK_MASTER_LOG/spark-master.out
  sbin/start-master.sh --host $SPARK_MASTER_HOST --port $SPARK_MASTER_PORT --webui-port $SPARK_MASTER_WEBUI_PORT
  tail -f /spark/logs/*
fi
