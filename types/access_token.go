package types

type AccessToken struct {
	UID      string `json:"uid" gorm:"primarykey;column:uid"`
	Name     string `json:"name" gorm:"column:name"`
	Token    string `json:"token" gorm:"column:token"`
	UserUID  string `json:"user_uid" gorm:"column:user_uid"`
	Created  int64  `json:"created" gorm:"column:created"`
	Expires  int64  `json:"expires" gorm:"column:expires"`
	UseCount int    `json:"usecount" gorm:"column:usecount"`
}

func (a *AccessToken) TableName() string {
	return "access_token"
}
