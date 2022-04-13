package Database

import (
	_ "github.com/go-sql-driver/mysql"
	"../Log"
	"database/sql"
	"fmt"
	"time"
)

const (
	Max_Open_Conn_Count= 20
	Max_Idle_Conn_Count= 10
)

type MySql struct {
	db *sql.DB
}

func (this *MySql) Open(addr string, dbName string, userName string, pwd string)  bool{
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", userName, pwd, addr, dbName)
	var err error
	this.db, err = sql.Open("mysql", dns)
	if err != nil {
		Log.WriteLog(Log.Log_Level_Error, "MySql Open(%s %s %s %s) error %s", addr, dbName, userName, pwd, err.Error())
		return false
	}

	err = this.db.Ping()
	if err != nil {
		Log.WriteLog(Log.Log_Level_Error, "MySql Open(%s %s %s %s) Failure, %s", addr, dbName, userName, pwd, err.Error())
		return false
	}

	this.db.SetMaxOpenConns(Max_Open_Conn_Count)
	this.db.SetMaxIdleConns(Max_Idle_Conn_Count)
	this.db.SetConnMaxIdleTime(time.Minute * 30)
	this.db.SetConnMaxLifetime(0)
	return true
}

func (this *MySql) Query(sql string, args ...interface{}) *sql.Rows{
	rows, err := this.db.Query(sql, args...)
	if err != nil {
		Log.WriteLog(Log.Log_Level_Error, "MySql Query %s error %s", sql, err.Error())
		return nil
	}
	return rows
}

func (this *MySql) Execute(sql string, args ...interface{} ) sql.Result{
	result, err := this.db.Exec(sql, args...)
	if err != nil {
		Log.WriteLog(Log.Log_Level_Error, "MySql Execute %s error %s", sql, err.Error())
		return nil
	}
	return result
}

//执行事务
func (this *MySql) ExecuteTx(sqls []string, argsList... []interface{}) bool{
	tx, err := this.db.Begin()
	if err != nil {
		Log.WriteLog(Log.Log_Level_Error, "MySql ExecuteTx Begin error %s", err.Error())
		return false
	}

	argsNum := len(argsList)
	for i, sqlstr := range sqls {
		var err error
		if i < argsNum {
			_, err = tx.Exec(sqlstr, argsList[i]...)
		} else {
			_, err = tx.Exec(sqlstr)
		}
		if err != nil{
			Log.WriteLog(Log.Log_Level_Error, "MySql ExecuteTx %s error %s", sqlstr, err.Error())
			err = tx.Rollback()
			if err != nil {
				Log.WriteLog(Log.Log_Level_Error, "MySql ExecuteTx RollBack error %s", err.Error())
			}
			return false
		}
	}

	err = tx.Commit()
	if err != nil {
		Log.WriteLog(Log.Log_Level_Error, "MySql ExecuteTx Commit error %s", err.Error())
		return false
	}
	return true
}
