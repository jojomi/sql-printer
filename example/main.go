package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jojomi/sqlprinter"
	"github.com/jojomi/strtpl"
	"github.com/juju/errors"
)

func main() {
	db, closer, err := getDatabaseConnection()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer closer()

	sql := strtpl.MustEval(`
	SELECT * FROM {{ .tablename }}
	`, map[string]interface{}{
		"tablename": "mytable",
	})
	fmt.Println(sql)

	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	err = sqlprinter.Table(rows)
	if err != nil {
		log.Fatal(err)
	}
}

func getDatabaseConnection() (*sql.DB, func(), error) {
	definition := struct {
		Username string
		Password string
		Hostname string
		Port     int
		Database string
	}{
		Username: "root",
		Password: "root1234",
		Hostname: "localhost",
		Port:     63306,
		Database: "sqlprinter",
	}
	addr := strtpl.MustEval("{{ .Username }}:{{ .Password }}@tcp({{ .Hostname }}:{{ .Port }})/{{ .Database }}?charset=utf8&parseTime=True&readTimeout=600s&writeTimeout=600s", definition)
	fmt.Println(addr)
	db, err := sql.Open("mysql", addr)
	/// db.SetConnMaxLifetime(0 * time.Second)
	if err != nil {
		return nil, nil, errors.Annotatef(err, "Error happend when connecting to DB %s?", definition.Database)
	}
	db.SetMaxIdleConns(0)
	if err = db.Ping(); err != nil {
		return nil, nil, errors.Annotatef(err, "Error happend when pinging to DB %s?", definition.Database)
	}
	return db, func() {
		db.Close()
	}, nil
}
