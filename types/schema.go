package types

import "encoding/json"

type SchemaUpdateError struct {
	NameTaken   bool `json:"name_taken"`
	NameInvalid bool `json:"name_invalid"`
}

// used for the database and GET requests
type Schema struct {
	ID               *int64 `json:"id" ksql:"id"`
	Created          int64  `json:"created" ksql:"created"`
	Mtime            int64  `json:"mtime" ksql:"mtime"`
	UserID           int64  `json:"user_id" ksql:"user_id"`
	Name             string `json:"name" ksql:"name"`
	Description      string `json:"description" ksql:"description"`
	ShortDescription string `json:"short_description" ksql:"short_description"`
	CDBCollection    string `json:"cdb_collection" ksql:"cdb_collection"`
	Complete         bool   `json:"complete" ksql:"complete"`
	SizeX            int    `json:"size_x" ksql:"size_x"`
	SizeY            int    `json:"size_y" ksql:"size_y"`
	SizeZ            int    `json:"size_z" ksql:"size_z"`
	TotalSize        int    `json:"total_size" ksql:"total_size"`
	TotalParts       int    `json:"total_parts" ksql:"total_parts"`
	Downloads        int    `json:"downloads" ksql:"downloads"`
	Views            int    `json:"views" ksql:"views"`
	License          string `json:"license" ksql:"license"`
	Stars            int    `json:"stars" ksql:"stars"`
}

func (s *Schema) UnmarshalJSON(data []byte) error {
	m := make(map[string]any)
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	id := getInt64(m["id"])
	s.ID = &id
	s.Created = getInt64(m["created"])
	s.Mtime = getInt64(m["mtime"])
	s.UserID = getInt64(m["user_id"])
	s.Name = getString(m["name"])
	s.Description = getString(m["description"])
	s.ShortDescription = getString(m["short_description"])
	s.CDBCollection = getString(m["cdb_collection"])
	s.License = getString(m["license"])
	s.Complete = getBool(m["complete"])
	s.SizeX = getInt(m["size_x"])
	s.SizeY = getInt(m["size_y"])
	s.SizeZ = getInt(m["size_z"])
	s.TotalSize = getInt(m["total_size"])
	s.TotalParts = getInt(m["total_parts"])
	s.Downloads = getInt(m["downloads"])

	return nil
}
