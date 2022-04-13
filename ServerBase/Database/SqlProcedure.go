package Database

import (
	"fmt"
	"strconv"
	"strings"
)

type SqlProcedure struct{
	Sql
	procname string
	params []string
}

func(this *SqlProcedure)AddParam_Str(param string) {
	this.params = append(this.params, fmt.Sprint("'", param, "'"))
}

func(this *SqlProcedure)AddParam_Int(param int64) {
	this.params = append(this.params, strconv.FormatInt(param, 10))
}

func(this *SqlProcedure)Statement() string{
	this.sql.WriteString("call ")
	this.sql.WriteString(this.table)
	this.sql.WriteString("(")
	this.sql.WriteString(strings.Join(this.params,","))
	this.sql.WriteString(")")
	return this.sql.String()
}