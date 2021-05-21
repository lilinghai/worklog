package main

import "time"

/*
go-tpc tpcc run --warehouses 10000 -H 172.16.4.157 -P 4009 -D tpcc -T 500
tiup cluster  scale-out tiflash-test scale-out.yaml  -u llh -p
./bin/go-tpc tpcc check --warehouses 10000 -H 172.16.4.204 -P 4009 -D tpcc -T 1 （with tiflash engine）
*/
func main() {
	scaleOut := "tiup cluster scale-out simple topology/scale-out.yaml -u llh -y"
	shellCommand(scaleOut)
	go func() {
		shellCommand("/home/llh/upgrade/go-tpc/bin/go-tpc tpcc check --warehouses 10000 -H 172.16.4.204 -P 4009 -D tpcc -T 10")
	}()
	go func() {
		shellCommand("go-tpc tpcc run --warehouses 10000 -H 172.16.4.157 -P 4009 -D tpcc -T 500")
	}()
	time.Sleep(20 * time.Hour)
}
