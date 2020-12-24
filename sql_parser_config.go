package sql_to_go

import (
	"fmt"
	"os"
	"strings"

	"github.com/gobeam/stringy"
	"github.com/pingcap/parser/mysql"
)

var DefaultSQLParserConfig = SQLParserConfig{
	IsStrictMode: func() bool {
		m := strings.ToLower(os.Getenv("STRICT_MODE"))
		return m == "true" || m == "1" || m == "on"
	}(),
	StructNameMap: func(tableName string) string {
		return stringy.New(tableName).CamelCase("?", "")
	},
	FieldNameMap: func(column string) string {
		return stringy.New(column).CamelCase("?", "")
	},
	FieldTypeMap: func(typ byte, opts ColumnOptions) string {
		isNullable := opts.IsNullable
		switch typ {
		case mysql.TypeTiny, mysql.TypeShort, mysql.TypeInt24, mysql.TypeLong, mysql.TypeLonglong,
			mysql.TypeBit, mysql.TypeYear:
			if isNullable {
				return "sql.NullInt64"
			}
			return "int64"
		case mysql.TypeFloat, mysql.TypeDouble:
			if isNullable {
				return "sql.NullFloat64"
			}
			return "float64"
		case mysql.TypeNewDecimal:
			if isNullable {
				return "sql.NullFloat64"
			}
			return "float64"
		case mysql.TypeDate, mysql.TypeDatetime:
			if isNullable {
				return "sql.NullString"
			}
			return "string"
		case mysql.TypeTimestamp:
			if isNullable {
				return "sql.NullInt64"
			}
			return "int64"
		case mysql.TypeDuration:
			if isNullable {
				return "sql.NullInt64"
			}
			return "int64"
		case mysql.TypeJSON:
			if isNullable {
				return "sql.NullString"
			}
			return "string"
		}

		if isNullable {
			return "sql.NullString"
		}
		return "string"
	},
	FieldTagMap: func(column string) string {
		tags := os.Getenv("TAGS")
		if tags == "" {
			tags = "db"
		}

		var tagArr []string
		for _, t := range strings.Split(tags, ",") {
			t = strings.TrimSpace(t)
			tagArr = append(tagArr, fmt.Sprintf("%s:\"%s\"", t, column))
		}

		return fmt.Sprintf("`%s`", strings.Join(tagArr, " "))
	},
}

type SQLParserConfig struct {
	IsStrictMode  bool // 出现warn会返回error
	StructNameMap func(tableName string) string
	FieldNameMap  func(column string) string
	FieldTypeMap  func(typ byte, opts ColumnOptions) string
	FieldTagMap   func(column string) string
}
