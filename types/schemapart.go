package types

import (
	"encoding/base64"
	"encoding/json"
)

func GetSchemaPartOrderID(offset_x, offset_y, offset_z int) int {
	return offset_x + (offset_z * 2000) + (offset_y * 2000 * 2000)
}

type SchemaPartIterator func() (*SchemaPart, error)

type SchemaPart struct {
	OrderID   int64  `json:"order_id" gorm:"column:order_id"`
	SchemaUID string `json:"schema_uid" gorm:"primarykey;column:schema_uid"`
	OffsetX   int    `json:"offset_x" gorm:"primarykey;column:offset_x"`
	OffsetY   int    `json:"offset_y" gorm:"primarykey;column:offset_y"`
	OffsetZ   int    `json:"offset_z" gorm:"primarykey;column:offset_z"`
	Mtime     int64  `json:"mtime" gorm:"column:mtime"`
	Data      []byte `json:"data" gorm:"column:data"`
	MetaData  []byte `json:"metadata" gorm:"column:metadata"`
}

func (sp *SchemaPart) TableName() string {
	return "schemapart"
}

func (s *SchemaPart) UnmarshalJSON(data []byte) error {
	m := make(map[string]any)
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

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
	if err != nil {
		return err
	}

	return nil
}

func (s SchemaPart) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["schema_uid"] = s.SchemaUID
	m["offset_x"] = s.OffsetX
	m["offset_y"] = s.OffsetY
	m["offset_z"] = s.OffsetZ
	m["mtime"] = s.Mtime
	m["data"] = base64.RawStdEncoding.EncodeToString(s.Data)
	m["metadata"] = base64.RawStdEncoding.EncodeToString(s.MetaData)

	return json.Marshal(m)
}
