package types

type Tag struct {
	UID         string `json:"uid" gorm:"primarykey;column:uid"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Restricted  bool   `json:"restricted" gorm:"column:restricted"`
}

func (t *Tag) TableName() string {
	return "tag"
}
