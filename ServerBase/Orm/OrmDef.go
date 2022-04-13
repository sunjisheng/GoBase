package Orm

import "encoding/xml"

const (
	OrmType_KV = "kv"
	OrmType_KMap = "kmap"
)

const (
	FieldType_String = "string"
	FieldType_Int = "int"
	FieldType_UInt32 = "uint32"
	FieldType_UInt64 = "uint64"
)

const (
	RelationType_Key = "key"
	RelationType_Field = "field"
)

type SolutionDef struct{
	Solution xml.Name  `xml:"solution"`
	PackageName string `xml:"package"`
	Tables []TableDef `xml:"table"`
}

type TableDef struct {
	TableName string `xml:"name"`
	TableType string `xml:"type"`
	Cache int `xml:"cache"`
	Fields []FieldDef `xml:"field"`
}

type FieldDef struct {
	FieldName string `xml:"name"`
	FieldType string `xml:"type"`
	RelationType string `xml:"relation"`
}