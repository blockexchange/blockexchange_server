package types

type SchemaMod struct {
	SchemaUID string `json:"schema_uid" gorm:"primarykey;column:schema_uid"`
	ModName   string `json:"mod_name" gorm:"primarykey;column:mod_name"`
}

func (sm *SchemaMod) TableName() string {
	return "schemamod"
}

type SchemaModCount struct {
	ModName string `json:"mod_name" gorm:"column:mod_name"`
	Count   int64  `json:"count" gorm:"column:count"`
}
