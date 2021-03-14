package types

// used for the database and GET requests
type Schema struct {
	ID          int64  `db:"id"`
	Created     int64  `db:"created"`
	UserID      int64  `db:"user_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Complete    bool   `db:"complete"`
	MaxX        int    `db:"max_y"`
	MaxY        int    `db:"max_y"`
	MaxZ        int    `db:"max_z"`
	PartLength  int    `db:"part_length"`
	TotalSize   int    `db:"total_size"`
	TotalParts  int    `db:"total_parts"`
	Downloads   int    `db:"downloads"`
	License     string `db:"license"`
}

// used for POST/PUT requests
type JsonSchema struct {
	UserID      float64 `json:"user_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Complete    bool    `json:"complete"`
	MaxX        float64 `json:"max_x"`
	MaxY        float64 `json:"max_y"`
	MaxZ        float64 `json:"max_z"`
	PartLength  float64 `json:"part_length"`
	TotalSize   float64 `json:"total_size"`
	TotalParts  float64 `json:"total_parts"`
	License     string  `json:"license"`
}

func MapSchema(j JsonSchema) Schema {
	s := Schema{}
	s.UserID = int64(j.UserID)
	s.Name = j.Name
	s.Description = j.Description
	s.Complete = j.Complete
	s.MaxX = int(j.MaxX)
	s.MaxY = int(j.MaxY)
	s.MaxZ = int(j.MaxZ)
	s.PartLength = int(j.PartLength)
	s.TotalSize = int(j.TotalSize)
	s.TotalParts = int(j.TotalParts)
	s.License = j.License
	return s
}
