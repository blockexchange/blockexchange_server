package types

import (
	"encoding/base64"
	"encoding/json"
)

func GetSchemaPartOrderID(offset_x, offset_y, offset_z int) int {
	return offset_x + (offset_z * 2000) + (offset_y * 2000 * 2000)
}

type SchemaPartIterator func() (*SchemaPart, error)

// SchemaPart is one chunk of a schematic, stored in the "schemapart" table.
// A schematic is split into parts of up to 16x16x16 nodes; each part is
// identified by its schematic and its offset inside that schematic.
//
// Custom MarshalJSON/UnmarshalJSON methods transfer Data and MetaData as
// unpadded base64 (RawStdEncoding) and leave OrderID out of the JSON.
type SchemaPart struct {
	// OrderID is a sort key derived from the offset (see GetSchemaPartOrderID).
	// The repository sets it on save; it is used to page through the parts of a
	// schematic in a stable order. Not part of the JSON payload.
	OrderID int64 `json:"order_id" gorm:"column:order_id"`
	// SchemaUID is the UID of the schematic this part belongs to.
	SchemaUID string `json:"schema_uid" gorm:"primarykey;column:schema_uid"`
	// OffsetX, OffsetY and OffsetZ are the part's position inside the
	// schematic, in nodes. They are multiples of 16 (divide by 16 for the part
	// index). Together with SchemaUID they form the primary key.
	OffsetX int `json:"offset_x" gorm:"primarykey;column:offset_x"`
	OffsetY int `json:"offset_y" gorm:"primarykey;column:offset_y"`
	OffsetZ int `json:"offset_z" gorm:"primarykey;column:offset_z"`
	// Mtime is the last modification time in Unix milliseconds. The API sets it
	// on upload; clients use it to fetch only parts changed since a given time.
	Mtime int64 `json:"mtime" gorm:"column:mtime"`
	// Data is the zlib-compressed node data for size = X*Y*Z nodes, laid out as
	// size big-endian uint16 node IDs (stored +32768, mapped to names through
	// the metadata's node_mapping), then size param1 bytes, then size param2
	// bytes.
	Data []byte `json:"data" gorm:"column:data"`
	// MetaData is zlib-compressed JSON with the part's "size" {x,y,z}, its
	// "node_mapping" (node name to ID) and optional per-node "metadata"
	// (inventories, fields, timers).
	MetaData []byte `json:"metadata" gorm:"column:metadata"`
}

func (sp *SchemaPart) TableName() string {
	return "schemapart"
}

// fromMap fills the SchemaPart from a decoded JSON object (see UnmarshalJSON).
func (s *SchemaPart) fromMap(m map[string]any) error {
	var err error
	s.SchemaUID = getString(m["schema_uid"])
	s.OffsetX = getInt(m["offset_x"])
	s.OffsetY = getInt(m["offset_y"])
	s.OffsetZ = getInt(m["offset_z"])
	s.Mtime = getInt64(m["mtime"])
	s.Data, err = base64.RawStdEncoding.DecodeString(getString(m["data"]))
	if err != nil {
		return err
	}
	s.MetaData, err = base64.RawStdEncoding.DecodeString(getString(m["metadata"]))
	return err
}

// toMap returns the JSON object representation of the SchemaPart (see MarshalJSON).
func (s SchemaPart) toMap() map[string]any {
	return map[string]any{
		"schema_uid": s.SchemaUID,
		"offset_x":   s.OffsetX,
		"offset_y":   s.OffsetY,
		"offset_z":   s.OffsetZ,
		"mtime":      s.Mtime,
		"data":       base64.RawStdEncoding.EncodeToString(s.Data),
		"metadata":   base64.RawStdEncoding.EncodeToString(s.MetaData),
	}
}

func (s *SchemaPart) UnmarshalJSON(data []byte) error {
	m := make(map[string]any)
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	return s.fromMap(m)
}

func (s SchemaPart) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.toMap())
}
