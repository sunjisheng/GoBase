package Database

import (
	"fmt"
	"strconv"
	"strings"
)

type SqlInsert struct{
	Sql
	fields []string
	values []string
}

func(this *SqlInsert)AddField_Int(field string, value int64) {
	this.fields = append(this.fields, field)
	this.values = append(this.values, strconv.FormatInt(value, 10))
}

func(this *SqlInsert)AddField_Float(field string, value float64) {
	this.fields = append(this.fields, field)
	this.values = append(this.values, strconv.FormatFloat(value, 'f', 3,64))
}

func(this *SqlInsert)AddField_Str(field string, value string) {
	this.fields = append(this.fields, field)
	this.values = append(this.values, fmt.Sprint("'",value,"'"))
}

func(this *SqlInsert)AddField_Binary(field string) {
	this.fields = append(this.fields, field)
	this.values = append(this.values, "?")
}

func(this *SqlInsert)Statement() string{
	this.sql.WriteString("insert into ")
	this.sql.WriteString(this.table)
	this.sql.WriteString("(")
	this.sql.WriteString(strings.Join(this.fields,","))
	this.sql.WriteString(")values(")
	this.sql.WriteString(strings.Join(this.values, ","))
	this.sql.WriteString(")")
	return this.sql.String()
}