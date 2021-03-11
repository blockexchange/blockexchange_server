package types

type AccessToken struct {
	ID       int64  `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Token    string `json:"token" db:"name"`
	UserID   int64  `json:"user_id" db:"user_id"`
	Created  int64  `json:"created" db:"created"`
	Expires  int64  `json:"expires" db:"expires"`
	UseCount int    `json:"usecount" db:"usecount"`
}
