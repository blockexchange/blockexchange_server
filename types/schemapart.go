package types

import (
	"encoding/base64"
	"encoding/json"
)

type SchemaPartIterator func() (*SchemaPart, error)

type SchemaPart struct {
	UID      string `json:"uid" ksql:"uid"`
	SchemaID int64  `json:"schema_id" ksql:"schema_id"`
	OffsetX  int    `json:"offset_x" ksql:"offset_x"`
	OffsetY  int    `json:"offset_y" ksql:"offset_y"`
	OffsetZ  int    `json:"offset_z" ksql:"offset_z"`
	Mtime    int64  `json:"mtime" ksql:"mtime"`
	Data     []byte `json:"data" ksql:"data"`
	MetaData []byte `json:"metadata" ksql:"metadata"`
}

func (s *SchemaPart) UnmarshalJSON(data []byte) error {
	m := make(map[string]any)
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	s.SchemaID = getInt64(m["schema_id"])
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
	m["schema_id"] = s.SchemaID
	m["offset_x"] = s.OffsetX
	m["offset_y"] = s.OffsetY
	m["offset_z"] = s.OffsetZ
	m["mtime"] = s.Mtime
	m["data"] = base64.RawStdEncoding.EncodeToString(s.Data)
	m["metadata"] = base64.RawStdEncoding.EncodeToString(s.MetaData)

	return json.Marshal(m)
}
