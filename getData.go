package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 3030
	user     = "postgres"
	password = "qqee"
	dbname   = "TimeBoard"
)

type Data struct {
	id        int
	name      string
	link      string
	Transport string
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

func GetData() map[string]string {
	m := make(map[string]string)

	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", conn)
	checkErr(err)
	defer db.Close()

	err = db.Ping()
	checkErr(err)

	fmt.Println("Connect")

	rows, err := db.Query("SELECT * FROM public.\"TimeBoard\"")
	checkErr(err)

	defer rows.Close()

	for rows.Next() {
		d := Data{}
		err = rows.Scan(&d.id, &d.name, &d.link, &d.Transport)
		checkErr(err)
		// fmt.Println("*************************************")
		// fmt.Println(d.id, d.name, d.link, d.Transport)
		m[d.name] = d.link
	}

	return m
}
