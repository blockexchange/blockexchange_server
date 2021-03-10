package types

type Schema struct {
	ID          int64  `json:"id"`
	Created     int64  `json:"created"`
	UserID      int64  `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Complete    bool   `json:"complete"`
	MaxX        int    `json:"max_x"`
	MaxY        int    `json:"max_y"`
	MaxZ        int    `json:"max_z"`
	PartLength  int    `json:"part_length"`
	TotalSize   int    `json:"total_size"`
	TotalParts  int    `json:"total_parts"`
	Downloads   int    `json:"downloads"`
	License     string `json:"license"`
}
