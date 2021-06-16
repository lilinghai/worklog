package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
	"time"
)

// 大量的 update，insert，delete 之后，start tiflash 节点
func main() {
	var dsn string
	var clusterName string
	flag.StringVar(&clusterName, "cn", "simple", "cluster name")
	flag.StringVar(&dsn, "dsn", "root:@tcp(127.0.0.1:4000)/test", "dsn")
	flag.Parse()
	// tpcc customer, PRIMARY KEY (`c_w_id`,`c_d_id`,`c_id`)
	// 10k , 10, 3000
	warehouse := 10000
	dmlSql := []string{
		"update customer set c_payment_cnt=c_payment_cnt+1 where c_id = %d",
		"delete from customer where c_id= %d",
		"insert into customer(c_id,c_d_id,c_w_id,c_first,c_middle,c_last,c_street_1, c_street_2,c_city,c_state,c_zip,c_phone,c_since,c_credit,c_credit_lim,c_discount,c_balance,c_ytd_payment,c_payment_cnt,c_delivery_cnt,c_data)" +
			"select c_id,c_d_id,%d,c_first,c_middle,c_last,c_street_1," +
			"c_street_2,c_city,c_state,c_zip,c_phone,c_since,c_credit,c_credit_lim,c_discount," +
			"c_balance,c_ytd_payment,c_payment_cnt,c_delivery_cnt,c_data from customer where c_w_id=1",
	}
	apSql := []string{
		"select count(*) from customer",
		"select sum(c_payment_cnt) from customer",
	}
	mdb, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	for true {
		ubegin := 1
		dbegin := 5001
		log.Println("stop tiflash nodes " + clusterName)
		_, err := shellCommand(fmt.Sprintf("tiup cluster stop %s -R tiflash", clusterName))
		if err != nil {
			log.Fatalln(err)
		}
		time.Sleep(1 * time.Minute)
		_, err = mdb.Exec("set @@tidb_isolation_read_engines='tikv'")
		if err != nil {
			log.Println(err)
		}
		for _, s := range dmlSql {
			p := dbegin
			if strings.Contains(s, "update") {
				p = ubegin
			}
			if strings.Contains(s, "insert") {
				p = warehouse + ubegin
			}
			fmt.Println("execute sql " + fmt.Sprintf(s, p))
			_, err := mdb.Exec(fmt.Sprintf(s, p))
			if err != nil {
				log.Fatalln(err)
			}
			dbegin += 1
			ubegin += 1
		}
		log.Println("start tiflash nodes " + clusterName)
		_, err = shellCommand(fmt.Sprintf("tiup cluster start %s -R tiflash", clusterName))
		if err != nil {
			log.Fatalln(err)
		}
		_, err = mdb.Exec("set @@tidb_isolation_read_engines='tiflash'")
		if err != nil {
			log.Println(err)
		}
		for i := 0; i < 10; i++ {
			for _, aps := range apSql {
				log.Println("execute sql " + aps)
				_, err := selectCnt(mdb, aps)
				if err != nil {
					if strings.Contains(err.Error(), "Region epoch not match") ||
						strings.Contains(err.Error(), "Region is unavailable") ||
						strings.Contains(err.Error(), "close of nil channel") ||
						strings.Contains(err.Error(), "TiFlash server timeout") ||
						strings.Contains(err.Error(), "MPP Task canceled because it seems hangs") {
						log.Println(err)
					} else {
						log.Fatalln(err)
					}
				}
			}
		}
	}
	time.Sleep(20 * time.Hour)
}
