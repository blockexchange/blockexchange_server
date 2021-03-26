package types

type Collection struct {
	ID          int64  `json:"id" db:"id"`
	UserID      int64  `json:"user_id" db:"user_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}
