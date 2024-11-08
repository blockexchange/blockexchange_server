package types

import "encoding/json"

type SchemaUpdateError struct {
	NameTaken   bool `json:"name_taken"`
	NameInvalid bool `json:"name_invalid"`
}

// used for the database and GET requests
type Schema struct {
	UID              string  `json:"uid" gorm:"primarykey;column:uid"`
	Created          int64   `json:"created" gorm:"column:created"`
	Mtime            int64   `json:"mtime" gorm:"column:mtime"`
	UserUID          string  `json:"user_uid" gorm:"column:user_uid"`
	CollectionUID    *string `json:"collection_uid" gorm:"column:collection_uid"`
	Name             string  `json:"name" gorm:"column:name"`
	Description      string  `json:"description" gorm:"column:description"`
	ShortDescription string  `json:"short_description" gorm:"column:short_description"`
	CDBCollection    string  `json:"cdb_collection" gorm:"column:cdb_collection"`
	Complete         bool    `json:"complete" gorm:"column:complete"`
	SizeX            int     `json:"size_x" gorm:"column:size_x"`
	SizeY            int     `json:"size_y" gorm:"column:size_y"`
	SizeZ            int     `json:"size_z" gorm:"column:size_z"`
	TotalSize        int     `json:"total_size" gorm:"column:total_size"`
	TotalParts       int     `json:"total_parts" gorm:"column:total_parts"`
	Downloads        int     `json:"downloads" gorm:"column:downloads"`
	Views            int     `json:"views" gorm:"column:views"`
	License          string  `json:"license" gorm:"column:license"`
	Stars            int     `json:"stars" gorm:"column:stars"`
}

func (s *Schema) TableName() string {
	return "schema"
}

func (s *Schema) UnmarshalJSON(data []byte) error {
	m := make(map[string]any)
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	s.UID = getString(m["uid"])
	s.Created = getInt64(m["created"])
	s.Mtime = getInt64(m["mtime"])
	s.UserUID = getString(m["user_uid"])
	s.Name = getString(m["name"])
	s.Description = getString(m["description"])
	s.ShortDescription = getString(m["short_description"])
	s.CDBCollection = getString(m["cdb_collection"])
	cuid := getString(m["collection_uid"])
	if cuid != "" {
		s.CollectionUID = &cuid
	}
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
