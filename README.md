# sql-to-go
Database create table statement is converted to go struct

# How to installation
```
    go install github.com/shima-park/sql-to-go/cmd/sql2go
```

# How to use it
```
// generate all table
sql2go --dsn="username:password@protocol(address)/dbname?param=value" --particular=user,proc

// generate by sql file
sql2go --sql=my.sql
```

my.sql
```
CREATE TABLE `user` (
  `Host` char(60) COLLATE utf8_bin NOT NULL DEFAULT '',
  `User` char(32) COLLATE utf8_bin NOT NULL DEFAULT ''
  # ...
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='Users and global privileges';

CREATE TABLE `proc` (
  `db` char(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `name` char(64) NOT NULL DEFAULT ''
  # ...
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Stored Procedures';
```

you will see the following output
```
package unknown

type User struct {
	Host string `gorm:"Host" json:"Host"` //
	User string `gorm:"User" json:"User"` //

}

package unknown

type Proc struct {
	Db   string `gorm:"db" json:"db"`     //
	Name string `gorm:"name" json:"name"` //

}
```

# Environment variable configuration
| Environment variable     | Description     |
| :------------- | :------------- |
| PACKAGE       | Replace default package name       |
| TEMPLATE_FILEPATH       | Replace the default template for generating source code        |
| TAGS       | Generate multiple tags or replace default tags        |


### Use PACKAGE
```
export PACKAGE=model && sql2go --dsn="username:password@protocol(address)/dbname?param=value"
```

You will see the following output
```
package model // Has been changed to your custom package name

type User struct {
	Host                 string        `db:"Host"`                   //
	User                 string        `db:"User"`                   //
    ...
}
```

### Use TEMPLATE_FILEPATH
my.tpl
```
package {{.Package}}

var {{.Struct.Name}}Columns = []string{
    {{range $index, $field := .Struct.Fields}} "{{$field.Column}}", {{end}}
}

type {{.Struct.Name}} struct{
    {{range $index, $field := .Struct.Fields}} {{$field.Name}} {{$field.Type}} {{$field.Tag}} // {{$field.Comment}}
    {{end}}
}
```

```
export TEMPLATE_FILEPATH=my.tpl && sql2go --dsn="username:password@protocol(address)/dbname?param=value" --particular=user
```
You will see the following output
```
package unknown

var UserColumns = []string{
	"Host", "User",
}

type User struct {
	Host string `gorm:"Host" json:"Host"` //
	User string `gorm:"User" json:"User"` //

}
```

### Use TAGS
```
export TAGS=gorm,json && sql2go --dsn="username:password@protocol(address)/dbname?param=value"
```

You will see the following output
```
package unknown

type User struct {
	Host                 string        `gorm:"Host" json:"Host"`                                     
	User                 string        `gorm:"User" json:"User"`                                     
    ...
}
```
