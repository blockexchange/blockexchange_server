package types

type SchemaScreenshot struct {
	UID       string `json:"uid" gorm:"primarykey;column:uid"`
	SchemaUID string `json:"schema_uid" gorm:"column:schema_uid"`
	Created   int64  `json:"created" gorm:"column:created"`
	Type      string `json:"type" gorm:"column:type"`
	Title     string `json:"title" gorm:"column:title"`
	Data      []byte `json:"data" gorm:"column:data"`
}

func (s *SchemaScreenshot) TableName() string {
	return "schema_screenshot"
}
