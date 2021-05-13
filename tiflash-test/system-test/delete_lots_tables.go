package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

/*
 br restore 20k+ 小表（100 rows）或使用 sql coverage 直接生成
 apply snapshot 或 ingest sst 后删除所有的表数据，观察 gc
*/
func main() {
	var dsn string
	flag.StringVar(&dsn, "dsn", "root:@tcp(127.0.0.1:4000)/test", "dsn")
	flag.Parse()
	mdb, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	rows, err := mdb.Query(`show tables`)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	tnames := []string{}
	for rows.Next() {
		if rows.Err() != nil {
			log.Fatalln(err)
		}
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatalln(err)
		}
		tnames = append(tnames, tableName)
	}
	for i, tn := range tnames {
		log.Println(i, len(tnames))
		_, err := mdb.Exec(fmt.Sprintf("delete from %s", tn, tn))
		if err != nil {
			log.Fatalln(err)
		}
	}
}
