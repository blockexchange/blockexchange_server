package types

type LuantiMediaConnection struct {
	ID                   string `json:"id" gorm:"primarykey;column:id"`
	Hostname             string `json:"hostname" gorm:"column:hostname"`
	Port                 int    `json:"port" gorm:"column:port"`
	Username             string `json:"username" gorm:"column:username"`
	Password             string `json:"password" gorm:"column:password"`
	LastUpdatedTimestamp int64  `json:"last_updated_timestamp" gorm:"column:last_updated_timestamp"`
	LastUpdatedSuccess   bool   `json:"last_updated_success" gorm:"column:last_updated_success"`
	LastUpdatedMessage   string `json:"last_updated_message" gorm:"column:last_updated_message"`
}

func (c *LuantiMediaConnection) TableName() string {
	return "luanti_media_connection"
}
