package types

// schematic <- 1:1 -> schematic-pull
type SchematicPull struct {
	SchemaUID string `json:"schema_uid" ksql:"schema_uid"`
	PosX      string `json:"pos_x"`
	PosY      string `json:"pos_y"`
	PosZ      string `json:"pos_z"`
	Interval  int64  `json:"interval"` // interval in seconds
}

// schematic-pull <- 1:n -> schematic-pull-client
type SchematicPullClient struct {
	UID              string `json:"uid" ksql:"uid"`
	SchematicPullUID string `json:"schema_pull_uid" ksql:"schema_pull_uid"`
	Enabled          bool   `json:"enabled"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	Hostname         string `json:"hostname"`
	Port             int    `json:"port"`
	LastRun          int64  `json:"last_run"`
	LastError        bool   `json:"last_error"` // true = last run failed
	LastMessage      string `json:"last_message"`
}
