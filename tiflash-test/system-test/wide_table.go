package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"time"
)

// 宽表 ，热点写入
func main() {
	for i := 0; i < 10; i++ {
		fmt.Println(RandString(10))

	}
	var dsn string
	var threads int
	flag.IntVar(&threads, "threads", 1, "threads to insert")
	flag.StringVar(&dsn, "dsn", "root:@tcp(127.0.0.1:4000)/test", "dsn")
	flag.Parse()
	mdb, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	maxLen := 102400
	createSql := `create table t(a bigint not null auto_increment,b varchar(10240),c text,d mediumtext,e longtext,primary key(a))`
	insertSql := `insert into t(b,c,d,e) values("%s","%s","%s","%s")`
	_, err = mdb.Exec(createSql)
	if err != nil {
		log.Println(err)
	}
	_, err = mdb.Exec("alter table t set tiflash replica 2")
	if err != nil {
		log.Println(err)
	}
	for i := 0; i < threads; i++ {
		go func() {
			rand.Seed(time.Now().UnixNano())
			_, err = mdb.Exec(fmt.Sprintf(insertSql,
				RandString(rand.Intn(maxLen)),
				RandString(rand.Intn(maxLen)),
				RandString(rand.Intn(maxLen)),
				RandString(rand.Intn(maxLen))))
			if err != nil {
				log.Println(err)
			}
		}()
	}

	time.Sleep(20 * time.Hour)
}
