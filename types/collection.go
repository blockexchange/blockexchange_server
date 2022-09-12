package types

type Collection struct {
	ID          int64  `json:"id" db:"id"`
	UserID      int64  `json:"user_id" db:"user_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

func (cs *Collection) Table() string {
	return "collection"
}

func (cs *Collection) Columns(action string) []string {
	return []string{"id", "user_id", "name", "description"}
}

func (cs *Collection) Values(action string) []any {
	switch action {
	case "insert", "update":
		return []any{cs.UserID, cs.Name, cs.Description}
	}
	return []any{cs.ID, cs.UserID, cs.Name, cs.Description}
}

func (cs *Collection) Scan(action string, r func(dest ...any) error) error {
	switch action {
	case "insert", "update":
		return r(&cs.UserID, &cs.Name, &cs.Description)
	}
	return r(&cs.ID, &cs.UserID, &cs.Name, &cs.Description)
}
