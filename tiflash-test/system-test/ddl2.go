package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

/*
循环执行一下步骤
tiflash replica 0
ddl
tiflash replica 1
ap select
 */

func main() {
	// tpcc customer, PRIMARY KEY (`c_w_id`,`c_d_id`,`c_id`)
	// 10k , 10, 3000
	//warehouse := 10000
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
		"alter table t drop index extra",
		"alter table customer drop column extra",
		"rename table customer to customer2",
		"rename table customer2 to customer",

		// recovery table
		//alter table partition is unsupported
		//"truncate table customer",
		//"drop table customer",

	}
	apSql := []string{
		"select count(*) from customer",
		"select sum(c_payment_cnt) from customer",
	}
	mdb, err := sql.Open("mysql", "root:@tcp(172.16.6.27:4019)/tpcc")
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		for true {
			ddlSql2 := append(ddlSql, ddlDropSql...)
			for _, s := range ddlSql2 {
				_, err = mdb.Exec("alter table customer set tiflash replica 0")
				if err != nil {
					log.Fatalln(err)
				}
				_, err := mdb.Exec(s)
				if err != nil {
					log.Fatalln(err)
				}
				_, err = mdb.Exec("alter table customer set tiflash replica 1")
				if err != nil {
					log.Fatalln(err)
				}
				available := 0
				for available == 0 {
					time.Sleep(30 * time.Second)
					available = selectCnt(mdb, `select count(*) from information_schema.tiflash_replica where TABLE_SCHEMa="tpcc" and TABLE_NAME="customer" and AVAILABLE=1`)
				}
				_, err = mdb.Exec("set @@tidb_isolation_read_engines='tiflash'")
				if err != nil {
					log.Println(err)
				}
				for i:=0;i<10;i++{
					for _, aps := range apSql {
						selectCnt(mdb, aps)
					}
				}
			}
		}
	}()
	time.Sleep(20 * time.Hour)
}
