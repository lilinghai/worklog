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

/*
循环执行一下步骤
stop tiflash node
ddl
start tiflash node
ap select
*/
func main() {
	var dsn string
	var clusterName string
	flag.StringVar(&clusterName, "cn", "simple", "cluster name")
	flag.StringVar(&dsn, "dsn", "root:@tcp(127.0.0.1:4000)/test", "dsn")
	flag.Parse()
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
		"alter table customer drop index extra",
		"alter table customer drop column extra",
		//"rename table customer to customer2",
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
	mdb, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		for true {
			ddlSql2 := append(ddlSql, ddlDropSql...)
			for _, s := range ddlSql2 {
				log.Println("stop tiflash nodes " + clusterName)
				_, err := shellCommand(fmt.Sprintf("tiup cluster stop %s -R tiflash", clusterName))
				if err != nil {
					log.Fatalln(err)
				}
				time.Sleep(1 * time.Minute)
				log.Println("execute sql " + s)
				_, err = mdb.Exec(s)
				if err != nil {
					log.Fatalln(err)
				}
				log.Println("start tiflash nodes " + clusterName)
				_, err = shellCommand(fmt.Sprintf("tiup cluster start %s -R tiflash", clusterName))
				if err != nil {
					log.Fatalln(err)
				}
				// 是否需要等待 2 分钟？
				//time.Sleep(2 * time.Minute)
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
		}
	}()
	time.Sleep(20 * time.Hour)
}
