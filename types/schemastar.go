package types

type SchemaStar struct {
	UserUID   string `json:"user_uid" gorm:"primarykey;column:user_uid"`
	SchemaUID string `json:"schema_uid" gorm:"primarykey;column:schema_uid"`
}

func (st *SchemaStar) TableName() string {
	return "user_schema_star"
}

type SchemaStarResponse struct {
	Count   int  `json:"count"`
	Starred bool `json:"starred"`
}
