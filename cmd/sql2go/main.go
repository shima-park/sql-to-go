package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/shima-park/sql-to-go"
)

var (
	dsn        = flag.String("dsn", "", "Data source name for connect to mysql. e.g. --dsn=\"username:password@protocol(address)/dbname?param=value\"")
	particular = flag.String("particular", "", "To generate a go structure table, if this option is not set, a go structure of all tables will be generated. e.g. --particular=user,proc")
	sql        = flag.String("sql", "", "Use sql file content to generate go struct. e.g. --sql=my.sql")
)

// go run cmd/sql2go/main.go --dsn="username:password@protocol(address)/dbname?param=value"
func main() {
	flag.Parse()

	c := sql_to_go.NewDefaultSQL2Go()
	if *dsn == "" && *sql == "" {
		fmt.Println("Please use flag --sql=my.sql or --dsn=\"username:password@protocol(address)/dbname?param=value\"")
		os.Exit(1)
	}

	var sqls []string
	if *dsn != "" {
		var particulars []string
		if *particular != "" {
			for _, t := range strings.Split(*particular, ",") {
				particulars = append(particulars, strings.ToLower(strings.TrimSpace(t)))
			}
		}

		tables, err := sql_to_go.ScanTable(*dsn, particulars...)
		if err != nil {
			fmt.Println("ScanTable error:", err)
			os.Exit(1)
		}

		for _, table := range tables {
			sqls = append(sqls, table.CreateTableSQL)
		}
	} else if *sql != "" {
		content, err := ioutil.ReadFile(*sql)
		if err != nil {
			fmt.Println("Read file:", *sql, " error:", err)
			os.Exit(1)
		}
		sqls = append(sqls, string(content))
	}

	for _, sql := range sqls {
		err := c.Convert(strings.NewReader(sql), os.Stdout)
		if err != nil {
			fmt.Println("Convert error:", err)
			os.Exit(1)
		}
	}
}
