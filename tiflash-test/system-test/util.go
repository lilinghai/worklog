package main

import (
	"database/sql"
	"log"
	"os/exec"
)

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

func shellCommand(cmd string) {
	out, err := exec.Command("bash", "-c", cmd).Output()

	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(out))
}