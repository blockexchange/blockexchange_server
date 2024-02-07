package types

type SchemaScreenshot struct {
	ID       int64  `json:"id" ksql:"id"`
	SchemaID int64  `json:"schema_id" ksql:"schema_id"`
	Type     string `json:"type" ksql:"type"`
	Title    string `json:"title" ksql:"title"`
	Data     []byte `json:"data" ksql:"data"`
}
