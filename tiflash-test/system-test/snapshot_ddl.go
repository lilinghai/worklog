package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
	"time"
)
// snapshot apply 的时候做 ddl 操作
func main() {
	// tpcc customer, PRIMARY KEY (`c_w_id`,`c_d_id`,`c_id`)
	// 10k , 10, 3000
	scaleOut := "tiup cluster scale-out simple topology/scale-out.yaml -u llh -y"
	ddlSql := []string{
		"alter table customer",
		"alter table customer add column extra int default 10",
		"alter table customer drop primary key",
		"alter table customer add primary key(c_w_id,c_d_id,c_id)",
		"alter table customer add index(extra)",
	}
	ddlDropSql := []string{
		"alter table customer rename column extra to extra2",
		"alter table customer rename column extra2 to extra",
		"alter table customer rename index extra to extra2",
		"alter table customer rename index extra2 to extra",
		"alter table customer drop index extra",
		"alter table customer drop column extra",
		"rename table customer to customer2",
		//"rename table customer2 to customer",

		// recovery table
		//alter table partition is unsupported
		//"truncate table customer",
		//"drop table customer",

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
	go func(){
		shellCommand("/home/llh/upgrade/go-tpc/bin/go-tpc tpcc check --warehouses 10000 -H 172.16.4.204 -P 4009 -D tpcc -T 10")
	}()
	go func() {
		for true {
			ddlSql2 := append(ddlSql, ddlDropSql...)
			for _, s := range ddlSql2 {
				_, err := mdb.Exec(s)
				if err != nil {
					log.Fatalln(err)
				}
				if strings.Contains(s,"customer2"){
					_, err := mdb.Exec("rename table customer2 to customer")
					if err != nil {
						log.Fatalln(err)
					}
				}
			}
		}
	}()
	time.Sleep(20 * time.Hour)
}
