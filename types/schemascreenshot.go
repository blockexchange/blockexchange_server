package types

type SchemaScreenshot struct {
	UID      string `json:"uid" ksql:"uid"`
	SchemaID int64  `json:"schema_id" ksql:"schema_id"`
	Type     string `json:"type" ksql:"type"`
	Title    string `json:"title" ksql:"title"`
	Data     []byte `json:"data" ksql:"data"`
}
