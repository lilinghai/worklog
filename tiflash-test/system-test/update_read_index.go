package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
	"time"
)

// 大量的 update，insert，delete 之后，start tiflash 节点
func main() {
	// tpcc customer, PRIMARY KEY (`c_w_id`,`c_d_id`,`c_id`)
	// 10k , 10, 3000
	warehouse := 10000
	stopTiflash := "tiup cluster stop tiflash-test -R tiflash"
	startTiflash := "tiup cluster stop tiflash-test -R tiflash"

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
	mdb, err := sql.Open("mysql", "root:@tcp(172.16.4.204:4009)/tpcc")
	if err != nil {
		log.Fatalln(err)
	}
	for true {
		ubegin := 1
		dbegin := 5001
		shellCommand(stopTiflash)
		time.Sleep(2 * time.Minute)
		for _, s := range dmlSql {
			p := dbegin
			if strings.Contains(s, "update") {
				p = ubegin
			}
			if strings.Contains(s, "insert") {
				p = warehouse + ubegin
			}
			_, err := mdb.Exec(fmt.Sprintf(s), p)
			if err != nil {
				log.Fatalln(err)
			}
			dbegin += 1
			ubegin += 1
		}
		shellCommand(startTiflash)
		available := selectCnt(mdb, `select count(*) from information_schema.tiflash_replica where TABLE_SCHEMa="tpcc" and TABLE_NAME="customer" and AVAILABLE=1`)
		for available == 0 {
			time.Sleep(30 * time.Second)
			available = selectCnt(mdb, `select count(*) from information_schema.tiflash_replica where TABLE_SCHEMa="tpcc" and TABLE_NAME="customer" and AVAILABLE=1`)
		}
		_, err = mdb.Exec("set @@tidb_allow_fallback_to_tikv='tiflash'")
		if err != nil {
			log.Println(err)
		}
		_, err = mdb.Exec("set @@tidb_isolation_read_engines='tiflash'")
		if err != nil {
			log.Println(err)
		}
		for i := 0; i < 10; i++ {
			for _, aps := range apSql {
				selectCnt(mdb, aps)
			}
		}
	}
	time.Sleep(20 * time.Hour)
}
