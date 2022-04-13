package Database

import (
"fmt"
"strings"
)

type SqlSelect struct{
	Sql
	fields []string
	conditions []string
	orderby string
}

func(this *SqlSelect)AddField(field string) {
	this.fields = append(this.fields, field)
}

func(this *SqlSelect)AddFields(fields string) {
	this.fields = strings.Split(fields, ",")
}

func(this *SqlSelect)AddCondition(field string, value int64) {
	this.conditions = append(this.conditions, fmt.Sprint(field, "=", value))
}

func(this *SqlSelect)AddCondition_Str(field string, value string) {
	this.conditions = append(this.conditions, fmt.Sprint(field, "='", value, "'"))
}

func(this *SqlSelect)AddOrderBy(orderby string) {
	this.orderby = orderby
}

func(this *SqlSelect)Statement() string{
	this.sql.WriteString("select ")
	this.sql.WriteString(strings.Join(this.fields, ","))
	this.sql.WriteString(" from ")
	this.sql.WriteString(this.table)
	if len(this.conditions) > 0 {
		this.sql.WriteString(" where ")
		this.sql.WriteString(strings.Join(this.conditions, " and "))
	}
	if this.orderby != "" {
		this.sql.WriteString(this.orderby)
	}
	return this.sql.String()
}