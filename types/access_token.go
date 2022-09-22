package types

type AccessToken struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Token    string `json:"token"`
	UserID   int64  `json:"user_id"`
	Created  int64  `json:"created"`
	Expires  int64  `json:"expires"`
	UseCount int    `json:"usecount"`
}

func (a *AccessToken) Table() string {
	return "access_token"
}

func (a *AccessToken) Columns(action string) []string {
	cols := []string{"name", "token", "user_id", "created", "expires", "usecount"}
	switch action {
	case "insert", "update":
		return cols
	default:
		return append([]string{"id"}, cols...)
	}
}

func (a *AccessToken) Values(action string) []any {
	vals := []any{a.Name, a.Token, a.UserID, a.Created, a.Expires, a.UseCount}
	switch action {
	case "insert", "update":
		return vals
	default:
		return append([]any{a.ID}, vals...)
	}
}

func (a *AccessToken) Scan(action string, r func(dest ...any) error) error {
	return r(&a.ID, &a.Name, &a.Token, &a.UserID, &a.Created, &a.Expires, &a.UseCount)
}
