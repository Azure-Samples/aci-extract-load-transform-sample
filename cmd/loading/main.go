package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var (
	defaultConnStr  = "postgresql://localhost?sslmode=disable"
	defaultFilePath = "./result.csv"
)

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func main() {
	filePath, ok := os.LookupEnv("FILE_PATH")
	if !ok {
		log.Println("FILE_PATH not set. Using default.")
		filePath = defaultFilePath
	}
	connStr, ok := os.LookupEnv("CONNECTION_STRING")
	if !ok {
		log.Println("CONNECTION_STRING not set. Using default.")
		connStr = defaultConnStr
	}

	for {
		exists, _ := pathExists(filePath)
		if exists {
			break
		}
		log.Println("Source file is not exist, retry in 1s")
		time.Sleep(1 * time.Second)
	}

	in, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(strings.NewReader(string(in)))

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(10)
	_, err = db.Query("do $do$ begin IF ( to_regclass('public.words') is null ) then create table words(name text not null primary key, count integer not null); end if; end $do$")
	if err != nil {
		panic(err)
	}
	for i := 1; i < len(records); i++ {
		query := fmt.Sprintf(`insert into "words" values('%s',%s) on conflict(name) do update set "count"=excluded."count"`, records[i][0], records[i][1])
		rows, err := db.Query(query)
		log.Printf("%d\n", i)
		if err != nil {
			log.Fatal(err)
		}
		rows.Close()
	}
}
