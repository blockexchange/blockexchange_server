package types

type Schema struct {
	ID          int64  `json:"id" db:"id"`
	Created     int64  `json:"created" db:"created"`
	UserID      int64  `json:"user_id" db:"user_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Complete    bool   `json:"complete" db:"complete"`
	MaxX        int    `json:"max_x" db:"max_y"`
	MaxY        int    `json:"max_y" db:"max_y"`
	MaxZ        int    `json:"max_z" db:"max_z"`
	PartLength  int    `json:"part_length" db:"part_length"`
	TotalSize   int    `json:"total_size" db:"total_size"`
	TotalParts  int    `json:"total_parts" db:"total_parts"`
	Downloads   int    `json:"downloads" db:"downloads"`
	License     string `json:"license" db:"license"`
}
