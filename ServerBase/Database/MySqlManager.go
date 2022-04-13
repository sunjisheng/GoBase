package Database

import "../Log"

var mysql_instances []*MySql

func InitMySqlInstance(maxDbCount uint32) {
	mysql_instances = make([]*MySql, maxDbCount)
}

func MySql_Instance(index uint32) *MySql {
	if mysql_instances[index] == nil {
		mysql_instances[index] = new(MySql)
	}
	return mysql_instances[index]
}

func ExecuteSQL(dbIndex uint32, sql ISql, args...interface{}) bool{
	mysql := MySql_Instance(dbIndex)
	if mysql == nil {
		Log.WriteLog(Log.Log_Level_Error, "ExecuteSQL MySqlPool_Instance().Get return nil")
		return false
	}
	result := mysql.Execute(sql.Statement(), args...)
	if result == nil {
		return false
	}
	_, err := result.RowsAffected()
	if err == nil {
		return true
	}
	return false
}

func ExecuteTx(dbIndex uint32, sqls []string, argsList...[]interface{}) bool{
	mysql := MySql_Instance(dbIndex)
	if mysql == nil {
		Log.WriteLog(Log.Log_Level_Error, "ExecuteTx MySqlPool_Instance().Get return nil")
		return false
	}
	return mysql.ExecuteTx(sqls, argsList...)
}