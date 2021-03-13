package types

type SchemaScreenshots struct {
	ID       int64  `json:"id" db:"id"`
	SchemaID int64  `json:"schema_id" db:"schema_id"`
	Type     string `json:"type" db:"type"`
	Title    string `json:"title" db:"title"`
	Data     []byte `json:"data" db:"data"`
}
