package types

import "encoding/json"

// used for the database and GET requests
type Schema struct {
	ID           int64  `db:"id" json:"id"`
	Created      int64  `db:"created" json:"created"`
	UserID       int64  `db:"user_id" json:"user_id"`
	Name         string `db:"name" json:"name"`
	Description  string `db:"description" json:"description"`
	Complete     bool   `db:"complete" json:"complete"`
	SizeX        int    `db:"size_x" json:"size_x"`
	SizeY        int    `db:"size_y" json:"size_y"`
	SizeZ        int    `db:"size_z" json:"size_z"`
	PartLength   int    `db:"part_length" json:"part_length"`
	TotalSize    int    `db:"total_size" json:"total_size"`
	TotalParts   int    `db:"total_parts" json:"total_parts"`
	Downloads    int    `db:"downloads" json:"downloads"`
	License      string `db:"license" json:"license"`
	SearchTokens string `db:"search_tokens"`
}

func (s *Schema) UnmarshalJSON(data []byte) error {
	m := make(map[string]interface{})
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	s.ID = getInt64(m["id"])
	s.Created = getInt64(m["created"])
	s.UserID = getInt64(m["user_id"])
	s.Name = getString(m["name"])
	s.Description = getString(m["description"])
	s.License = getString(m["license"])
	s.Complete = getBool(m["complete"])
	s.SizeX = getInt(m["size_x"])
	s.SizeY = getInt(m["size_y"])
	s.SizeZ = getInt(m["size_z"])
	s.PartLength = getInt(m["part_length"])
	s.TotalSize = getInt(m["total_size"])
	s.TotalParts = getInt(m["total_parts"])
	s.Downloads = getInt(m["downloads"])

	return nil
}
