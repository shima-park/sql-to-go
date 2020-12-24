package sql_to_go

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

type Table struct {
	Name           string
	CreateTableSQL string
}

func ScanTable(dsn string, particularTables ...string) ([]Table, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("show tables")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []Table
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			return nil, err
		}

		if len(particularTables) > 0 &&
			!stringInSlice(strings.ToLower(table), particularTables) {
			continue
		}

		ct := db.QueryRow("show create table " + table)
		var t Table
		err = ct.Scan(&t.Name, &t.CreateTableSQL)
		if err != nil {
			return nil, err
		}

		tables = append(tables, t)
	}

	return tables, nil
}

func stringInSlice(s string, slice []string) bool {
	for _, v := range slice {
		if s == v {
			return true
		}
	}
	return false
}
