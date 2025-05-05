package types

// schematic <- 1:1 -> schematic-pull
type SchematicPull struct {
	SchemaUID string `json:"schema_uid" gorm:"primarykey;column:schema_uid"`
	Enabled   bool   `json:"enabled" gorm:"column:enabled"`
	PosX      int    `json:"pos_x" gorm:"column:pos_x"`
	PosY      int    `json:"pos_y" gorm:"column:pos_y"`
	PosZ      int    `json:"pos_z" gorm:"column:pos_z"`
	Interval  int64  `json:"interval" gorm:"column:interval"` // interval in seconds
	NextRun   int64  `json:"next_run" gorm:"column:next_run"` // time.Now().UnixMilli()
	Hostname  string `json:"hostname" gorm:"column:hostname"`
	Port      int    `json:"port" gorm:"column:port"`
}

func (sp *SchematicPull) TableName() string {
	return "schema_pull"
}

// schematic-pull <- 1:n -> schematic-pull-client
type SchematicPullClient struct {
	UID              string `json:"uid" gorm:"primarykey;column:uid"`
	SchematicPullUID string `json:"schema_pull_uid" gorm:"column:schema_pull_uid"`
	Enabled          bool   `json:"enabled" gorm:"column:enabled"`
	Username         string `json:"username" gorm:"column:username"`
	Password         string `json:"password" gorm:"column:password"`
	LastRun          int64  `json:"last_run" gorm:"column:last_run"`     // time.Now().UnixMilli()
	LastError        bool   `json:"last_error" gorm:"column:last_error"` // true = last run failed
	LastMessage      string `json:"last_message" gorm:"column:last_message"`
}

func (spc *SchematicPullClient) TableName() string {
	return "schema_pull_client"
}
