package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	var dsn string
	var replicaNum int
	flag.IntVar(&replicaNum, "replica", 1, "tiflash replica number")
	flag.StringVar(&dsn, "dsn", "root:@tcp(127.0.0.1:4000)/test", "dsn example root:@tcp(127.0.0.1:4000)/test")
	mdb, err := sql.Open("mysql", "root:@tcp(172.16.6.27:4019)/tpcc")
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
		_, err := mdb.Exec(fmt.Sprintf("ALTER TABLE %s SET TIFLASH REPLICA %d", tn, replicaNum))
		if err != nil {
			log.Fatalln(err)
		}
	}
}
