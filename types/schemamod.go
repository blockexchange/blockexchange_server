package types

type SchemaMod struct {
	ID       int64  `json:"id" db:"id"`
	SchemaID int64  `json:"schema_id" db:"schema_id"`
	ModName  string `json:"mod_name" db:"mod_name"`
}

func (u *SchemaMod) Columns(action string) []string {
	cols := []string{}
	if action != "insert" {
		cols = append(cols, "id")
	}
	cols = append(cols, "schema_id", "mod_name")
	return cols
}

func (u *SchemaMod) Table() string {
	return "schemamod"
}

func (u *SchemaMod) Scan(action string, r func(dest ...any) error) error {
	return r(&u.ID, &u.SchemaID, &u.ModName)
}

func (u *SchemaMod) Values(action string) []any {
	vals := []any{}
	if action != "insert" {
		vals = append(vals, u.ID)
	}
	vals = append(vals, u.SchemaID, u.ModName)
	return vals
}
