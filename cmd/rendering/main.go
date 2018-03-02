package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

var (
	defaultConnStr = "postgresql://localhost?sslmode=disable"
)

func main() {
	connStr, ok := os.LookupEnv("CONNECTION_STRING")
	if !ok {
		log.Println("CONNECTION_STRING not set. Using default.")
		connStr = defaultConnStr
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(10)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			http.ServeFile(w, r, "./static/index.html")
		} else if r.Method == "POST" {
			w.Header().Set("Content-Type", "application/json;   charset=UTF-8")
			rows, err := db.Query("select * from words")
			if err != nil {
				panic(err)
			}
			defer rows.Close()
			var buf bytes.Buffer
			buf.WriteString("[")
			isFirstRecord := true
			for rows.Next() {
				var name string
				var count int64
				err = rows.Scan(&name, &count)
				if err != nil {
					fmt.Println(err) // Handle scan error
				} else {
					if !isFirstRecord {
						buf.WriteString(",")
					}
					buf.WriteString("{")
					buf.WriteString("\"text\":\"")
					buf.WriteString(name)
					buf.WriteString("\",\"size\":")
					buf.WriteString(strconv.FormatInt(count, 10))
					buf.WriteString("}")
					isFirstRecord = false
				}
			}
			buf.WriteString("]")
			fmt.Fprintf(w, buf.String())
		}
	})

	log.Fatal(http.ListenAndServe(":80", nil))
}
