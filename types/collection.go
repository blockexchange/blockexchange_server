package types

type Collection struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (cs *Collection) Table() string {
	return "collection"
}

func (cs *Collection) Columns(action string) []string {
	cols := []string{"user_id", "name", "description"}
	switch action {
	case "insert", "update":
		return cols
	default:
		return append([]string{"id"}, cols...)
	}
}

func (cs *Collection) Values(action string) []any {
	vals := []any{cs.UserID, cs.Name, cs.Description}
	switch action {
	case "insert", "update":
		return vals
	default:
		return append([]any{cs.ID}, vals...)
	}
}

func (cs *Collection) Scan(action string, r func(dest ...any) error) error {
	return r(&cs.ID, &cs.UserID, &cs.Name, &cs.Description)
}
