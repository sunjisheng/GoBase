package Database

import (
	"bytes"
)

type ISql interface{
	Statement() string
}

type Sql struct {
	table string
	sql bytes.Buffer
}

func(this *Sql)SetTable(table string) {
	this.table = table
}

func (this *Sql)SetSql (sql string) {
	this.sql.WriteString(sql)
}

func (this *Sql)Statement() string {
	return this.sql.String()
}

