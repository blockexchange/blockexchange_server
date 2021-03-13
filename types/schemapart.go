package types

type SchemaPart struct {
	ID       int64  `json:"id" db:"id"`
	SchemaID int64  `json:"schema_id" db:"schema_id"`
	OffsetX  int    `json:"offset_x" db:"offset_x"`
	OffsetY  int    `json:"offset_y" db:"offset_y"`
	OffsetZ  int    `json:"offset_z" db:"offset_z"`
	Mtime    int64  `json:"mtime" db:"mtime"`
	Data     []byte `json:"data" db:"data"`
	MetaData []byte `json:"metadata" db:"metadata"`
}
