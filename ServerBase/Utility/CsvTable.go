package Utility

import (
	"../Log"
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

const (
	FixHeader_Row_Name = 0
	FixHeader_Row_Count = 1
)

type CsvTable struct {
	FieldNames []string
	Rows map[uint32]map[string]string
}

func (this *CsvTable)  Open(fileName string) bool {
	this.FieldNames = make([]string, 0)
	this.Rows = make(map[uint32]map[string]string)

	fs, err := os.Open(fileName)
	if err != nil {
		Log.WriteLog(Log.Log_Level_Error, "CsvFile Open %s Failure", fileName)
	}
	reader := csv.NewReader(fs)
	reader.FieldsPerRecord = -1
	var row int
	for {
		record, err := reader.Read()
		if err != nil {
			if err != io.EOF {
				Log.WriteLog(Log.Log_Level_Error, "CsvFile Read %s Failure", fileName)
			}
			break
		}
		if row == FixHeader_Row_Name{
			if len(record) == 0 {
				break
			}
			for col := 0; col < len(record); col++{
				this.FieldNames = append(this.FieldNames, record[col])
			}
		} else if row >= FixHeader_Row_Count && len(record) > 0 {
			id, err:= strconv.Atoi(record[0])
			if err != nil {
				break
			}
			row := make(map[string]string)
			for col := 0; col < len(record); col++{
				row[this.FieldNames[col]] = record[col]
			}
			this.Rows[uint32(id)] = row
		}
		row++
	}
	return true
}

func (this *CsvTable) RowCount() int {
	return len(this.Rows)
}

func (this *CsvTable) FieldCount() int {
	return len(this.FieldNames)
}