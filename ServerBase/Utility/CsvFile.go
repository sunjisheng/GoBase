package Utility

import (
	"../Log"
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

const (
	Header_Row_Count = 1
)
type CsvFile struct {
	ids []uint32
	content map[uint32][]string
}

func (this *CsvFile)  Open(fileName string) bool {
	if this.IsOpen() {
		return true
	}
	this.ids = make([]uint32, 0)
	this.content = make(map[uint32][]string)

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
		//忽略Header Row
		if row >= Header_Row_Count && len(record) > 0 {
			id, err:= strconv.Atoi(record[0])
			if err != nil {
				break
			}
			this.ids =  append(this.ids, uint32(id))
			this.content[uint32(id)] = record
		}
		row++
	}
	return true
}

func (this *CsvFile) GetContent() map[uint32][]string {
	return this.content
}

func (this *CsvFile) IsOpen() bool {
	if len(this.content) > 0 {
		return true
	} else {
		return false
	}
}

func (this *CsvFile) RowCount() uint32 {
	return uint32(len(this.content))
}

func (this *CsvFile) GetRowID(row uint32) uint32 {
	if row < 0 || row >= uint32(len(this.ids)) {
		return 0
	}
	return this.ids[row]
}

func (this *CsvFile) GetString(id uint32, col uint32) string {
	row, ok := this.content[id]
	if !ok {
		return ""
	}
	if col >= uint32(len(row)) {
		return ""
	}
	return row[col]
}

func (this *CsvFile) GetInt(id uint32, col uint32) int32 {
	row, ok := this.content[id]
	if !ok {
		return 0
	}
	if uint32(len(row)) <= col {
		return 0
	}
	value, _:= strconv.Atoi(row[col])
	return int32(value)
}