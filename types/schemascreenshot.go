package types

type SchemaScreenshot struct {
	UID       string `json:"uid" ksql:"uid"`
	SchemaUID string `json:"schema_uid" ksql:"schema_uid"`
	Created   int64  `json:"created" ksql:"created"`
	Type      string `json:"type" ksql:"type"`
	Title     string `json:"title" ksql:"title"`
	Data      []byte `json:"data" ksql:"data"`
}
