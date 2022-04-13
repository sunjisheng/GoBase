package Orm

import "fmt"

type DBError struct {
	error string
}

func(this *DBError) Error() string {
	return fmt.Sprintf("Databse Error:%s", this.error)
}

func NewDBError(error string)  *DBError {
	return &DBError{error: error}
}
