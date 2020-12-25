package sql_to_go

type Struct struct {
	TableName string
	Name      string
	Fields    []Field
}

type Field struct {
	Column  string
	Name    string
	Type    string
	Tag     string
	Comment string
}
