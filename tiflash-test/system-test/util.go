package main

import (
	"database/sql"
	"os/exec"
)

// filter some known issues
// https://github.com/pingcap/tics/issues/1947
func selectCnt(d *sql.DB, sql string) (int, error) {
	rows, err := d.Query(sql)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		//if strings.Contains(err.Error(), "newer than query schema version")
		return 0, err
	}
	var cnt int
	for rows.Next() {
		if rows.Err() != nil {
			return 0, err
		}
		if err := rows.Scan(&cnt); err != nil {
			return 0, err
		}
	}
	return cnt, nil
}

func shellCommand(cmd string) ([]byte, error) {
	return exec.Command("bash", "-c", cmd).Output()
}
