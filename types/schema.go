package types

import "encoding/json"

type SchemaUpdateError struct {
	NameTaken   bool `json:"name_taken"`
	NameInvalid bool `json:"name_invalid"`
}

// used for the database and GET requests
type Schema struct {
	ID           int64  `json:"id"`
	Created      int64  `json:"created"`
	Mtime        int64  `json:"mtime"`
	UserID       int64  `json:"user_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Complete     bool   `json:"complete"`
	SizeX        int    `json:"size_x"`
	SizeY        int    `json:"size_y"`
	SizeZ        int    `json:"size_z"`
	TotalSize    int    `json:"total_size"`
	TotalParts   int    `json:"total_parts"`
	Downloads    int    `json:"downloads"`
	Views        int    `json:"views"`
	License      string `json:"license"`
	SearchTokens string `json:"-"`
}

func (s *Schema) UnmarshalJSON(data []byte) error {
	m := make(map[string]any)
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	s.ID = getInt64(m["id"])
	s.Created = getInt64(m["created"])
	s.Mtime = getInt64(m["mtime"])
	s.UserID = getInt64(m["user_id"])
	s.Name = getString(m["name"])
	s.Description = getString(m["description"])
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

func (a *Schema) Table() string {
	return "schema"
}

func (a *Schema) Columns(action string) []string {
	cols := []string{"created", "mtime", "user_id", "name", "description", "complete", "size_x", "size_y", "size_z", "total_size", "total_parts", "downloads", "views", "license"}
	switch action {
	case "insert", "update":
		return cols
	default:
		return append([]string{"id"}, cols...)
	}
}

func (a *Schema) Values(action string) []any {
	vals := []any{a.Created, a.Mtime, a.UserID, a.Name, a.Description, a.Complete, a.SizeX, a.SizeY, a.SizeZ, a.TotalSize, a.TotalParts, a.Downloads, a.Views, a.License}
	switch action {
	case "insert", "update":
		return vals
	default:
		return append([]any{a.ID}, vals...)
	}
}

func (a *Schema) Scan(action string, r func(dest ...any) error) error {
	return r(&a.ID, &a.Created, &a.Mtime, &a.UserID, &a.Name, &a.Description, &a.Complete, &a.SizeX, &a.SizeY, &a.SizeZ, &a.TotalSize, &a.TotalParts, &a.Downloads, &a.Views, &a.License)
}
