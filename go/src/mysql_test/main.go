// mysql_test project main.go
package main

import (
	"database/sql"
	"fmt"
	_ "mysql"
)

func main() {
	db, e := sql.Open("mysql", "root:@tcp(192.168.100.167:3388)/kslave?charset=utf8") //[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]

	defer db.Close()

	if e != nil {
		fmt.Printf("%v\n", e)
	}

	rows, error := db.Query("select mid, mtkey from ks_membertable where svid > ?", 0)

	if error != nil {
		fmt.Printf("%v", error)
	}
	defer rows.Close()

	cols, _ := rows.Columns()

	for k, v := range cols {
		if k > 0 {

		}

		fmt.Printf("%v\t", v)
	}

	fmt.Println("")

	var mid uint32
	var mtkey string

	for rows.Next() {
		if err := rows.Scan(&mid, &mtkey); err == nil {
			fmt.Printf("%v\t", mid)
			fmt.Printf("%v\n", mtkey)
		}
	}

}
