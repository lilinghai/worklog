package main

/*
go-tpc tpcc run --warehouses 10000 -H 172.16.4.157 -P 4009 -D tpcc -T 500
tiup cluster  scale-out tiflash-test scale-out.yaml  -u llh -p
./bin/go-tpc tpcc check --warehouses 10000 -H 172.16.4.204 -P 4009 -D tpcc -T 1 （with tiflash engine）
*/
