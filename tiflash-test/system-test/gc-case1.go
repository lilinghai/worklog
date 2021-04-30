package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)
/*
 br restore tpcc 10K warehouses
 gc test
1. 一边删除，一边查询对比 kv 和 cs 引擎的查询结果
2. gc 能够正常工作，tiflash metric 确认
3. gc 不影响查询的性能，如 cpu 资源消耗，latency
 */
func main() {
	mdb, err := sql.Open("mysql", "root:@tcp(172.16.4.204:4009)/tpcc")
	if err != nil {
		log.Fatalln(err)
	}
	for i := 0; i < 1; i++ {
		go func() {
			for j := 0; j < 10000; j++ {
				_, err = mdb.Exec("set @@tidb_isolation_read_engines='tiflash'")
				if err != nil {
					log.Println(err)
				}
				csc := selectCnt(mdb, "select count(*) from customer")
				_, err = mdb.Exec("set @@tidb_isolation_read_engines='tikv'")
				if err != nil {
					log.Println(err)
				}
				kvc := selectCnt(mdb, "select count(*) from customer")
				if csc != kvc {
					log.Fatalf("tics %d,tikv %d ,not equal", csc, kvc)
				}
				if kvc!=0{
					_, err := mdb.Exec("delete from customer limit 300000")
					if err != nil {
						log.Println(err)
					}
				}
			}
		}()
	}
	time.Sleep(20 * time.Hour)
}

func selectCnt(d *sql.DB, sql string) int {
	rows, err := d.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
	}
	var cnt int
	for rows.Next() {
		if rows.Err() != nil {
			log.Fatalln(err)
		}
		if err := rows.Scan(&cnt); err != nil {
			log.Fatalln(err)
		}
	}
	log.Println(cnt)
	return cnt
}
