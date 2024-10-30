package types

type Collection struct {
	UID         string `json:"uid" gorm:"primarykey;column:uid"`
	UserUID     string `json:"user_uid" gorm:"column:user_uid"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
}

func (c *Collection) TableName() string {
	return "collection"
}
