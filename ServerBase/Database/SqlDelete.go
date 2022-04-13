package Database

import (
	"fmt"
	"strings"
)

type SqlDelete struct{
	Sql
	conditions []string
}

func(this *SqlDelete)AddCondition(field string, value int64) {
	this.conditions = append(this.conditions, fmt.Sprint(field, "=", value))
}

func(this *SqlDelete)AddCondition_Str(field string, value string) {
	this.conditions = append(this.conditions, fmt.Sprint(field, "='", value, "'"))
}

func(this *SqlDelete)Statement() string{
	this.sql.WriteString("delete from ")
	this.sql.WriteString(this.table)
	if len(this.conditions) > 0 {
		this.sql.WriteString(" where ")
		this.sql.WriteString(strings.Join(this.conditions, " and "))
	}
	return this.sql.String()
}