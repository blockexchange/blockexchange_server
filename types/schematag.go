package types

type SchemaTag struct {
	UID       string `json:"uid" gorm:"primarykey;column:uid"`
	TagUID    string `json:"tag_uid" gorm:"column:tag_uid"`
	SchemaUID string `json:"schema_uid" gorm:"column:schema_uid"`
}

func (st *SchemaTag) TableName() string {
	return "schematag"
}
