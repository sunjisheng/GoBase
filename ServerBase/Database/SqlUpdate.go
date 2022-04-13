package Database

import (
	"fmt"
	"strconv"
	"strings"
)

type SqlUpdate struct{
	Sql
	field_values []string
	conditions []string
}

func(this *SqlUpdate)AddField_Int(field string, value int64) {
	this.field_values = append(this.field_values, fmt.Sprint(field, "=", value))
}

func(this *SqlUpdate)AddField_Float(field string, value float64) {
	this.field_values = append(this.field_values, fmt.Sprint(field, "=", strconv.FormatFloat(value, 'f', 3,64)))
}

func(this *SqlUpdate)AddField_Str(field string, value string) {
	this.field_values = append(this.field_values, fmt.Sprint(field, "='", value, "'"))
}

func(this *SqlUpdate)AddCondition(field string, value int64) {
	this.conditions = append(this.conditions, fmt.Sprint(field, "=", value))
}

func(this *SqlUpdate)AddCondition_Str(field string, value string) {
	this.conditions = append(this.conditions, fmt.Sprint(field, "='", value, "'"))
}

func(this *SqlUpdate)Statement() string{
	this.sql.WriteString("update ")
	this.sql.WriteString(this.table)
	this.sql.WriteString(" set ")
	this.sql.WriteString(strings.Join(this.field_values, ","))
	if len(this.conditions) > 0 {
		this.sql.WriteString(" where ")
		this.sql.WriteString(strings.Join(this.conditions, " and "))
	}
	return this.sql.String()
}