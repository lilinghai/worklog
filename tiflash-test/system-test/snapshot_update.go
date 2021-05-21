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
// snapshot apply 的时候 update 操作
func main() {
	// tpcc customer, PRIMARY KEY (`c_w_id`,`c_d_id`,`c_id`)
	// 10k , 10, 3000
	warehouse := 10000
	scaleOut := "tiup cluster scale-out simple topology/scale-out.yaml -u llh -p"
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
	shellCommand(scaleOut)
	go func(){
		_, err = mdb.Exec("set @@tidb_isolation_read_engines='tiflash'")
		if err != nil {
			log.Println(err)
		}
		for i := 0; i < 10; i++ {
			for _, aps := range apSql {
				selectCnt(mdb, aps)
			}
		}
	}()
	go func() {
		for true {
			ubegin := 1
			dbegin := 5001
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
		}
	}()
	time.Sleep(20 * time.Hour)
}
