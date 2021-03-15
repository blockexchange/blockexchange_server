package types

import "encoding/json"

// used for the database and GET requests
type Schema struct {
	ID          int64  `db:"id" json:"id"`
	Created     int64  `db:"created" json:"created"`
	UserID      int64  `db:"user_id" json:"user_id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	Complete    bool   `db:"complete" json:"complete"`
	MaxX        int    `db:"max_x" json:"max_x"`
	MaxY        int    `db:"max_y" json:"max_y"`
	MaxZ        int    `db:"max_z" json:"max_z"`
	PartLength  int    `db:"part_length" json:"part_length"`
	TotalSize   int    `db:"total_size" json:"total_size"`
	TotalParts  int    `db:"total_parts" json:"total_parts"`
	Downloads   int    `db:"downloads" json:"downloads"`
	License     string `db:"license" json:"license"`
}

func getInt64(o interface{}) int64 {
	v, _ := o.(float64)
	return int64(v)
}

func getInt(o interface{}) int {
	v, _ := o.(float64)
	return int(v)
}

func getString(o interface{}) string {
	s, _ := o.(string)
	return s
}

func getBool(o interface{}) bool {
	v, _ := o.(bool)
	return v
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
	s.MaxX = getInt(m["max_x"])
	s.MaxY = getInt(m["max_y"])
	s.MaxZ = getInt(m["max_z"])
	s.PartLength = getInt(m["part_length"])
	s.TotalSize = getInt(m["total_size"])
	s.TotalParts = getInt(m["total_parts"])
	s.Downloads = getInt(m["downloads"])

	return nil
}
