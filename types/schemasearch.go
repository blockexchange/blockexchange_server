package types

import "github.com/lib/pq"

var OrderDirections = map[string]bool{
	"asc":  true,
	"desc": true,
}

var OrderColumns = map[string]bool{
	"s.created": true,
}

type SchemaSearchRequest struct {
	UserID         *int64  `json:"user_id"`
	SchemaID       *int64  `json:"schema_id"`
	SchemaIDList   []int64 `json:"schema_id_list"`
	TagID          *int64  `json:"tag_id"`
	SchemaName     *string `json:"schema_name"`
	UserName       *string `json:"user_name"`
	Keywords       *string `json:"keywords"`
	Complete       *bool   `json:"complete"`
	OrderDirection *string `json:"order_direction"`
	OrderColumn    *string `json:"order_column"`
	Limit          *int    `json:"limit"`
	Offset         *int    `json:"offset"`
}

type SchemaSearchResult struct {
	Schema
	Stars    int            `json:"stars"`
	UserName string         `json:"username"`
	Mods     pq.StringArray `json:"mods"`
	Tags     pq.StringArray `json:"tags"`
}

func (s *SchemaSearchResult) Table() string {
	return "schema s join public.user u on s.user_id = u.id"
}

func (s *SchemaSearchResult) Columns(action string) []string {
	return []string{
		"s.id",
		"s.created",
		"s.mtime",
		"s.user_id",
		"s.name",
		"s.description",
		"s.short_description",
		"s.cdb_collection",
		"s.complete",
		"s.size_x",
		"s.size_y",
		"s.size_z",
		"s.total_size",
		"s.total_parts",
		"s.downloads",
		"s.views",
		"s.license",
		"u.name",
		"array(select name from schematag st join tag t on st.tag_id = t.id where schema_id = s.id)",
		"array(select mod_name from schemamod where schema_id = s.id)",
		"(select count(*) from user_schema_star where schema_id = s.id)",
	}
}

func (s *SchemaSearchResult) Scan(action string, r func(dest ...any) error) error {
	return r(
		&s.ID,
		&s.Created,
		&s.Mtime,
		&s.UserID,
		&s.Name,
		&s.Description,
		&s.ShortDescription,
		&s.CDBCollection,
		&s.Complete,
		&s.SizeX,
		&s.SizeY,
		&s.SizeZ,
		&s.TotalSize,
		&s.TotalParts,
		&s.Downloads,
		&s.Views,
		&s.License,
		&s.UserName,
		&s.Tags,
		&s.Mods,
		&s.Stars,
	)
}
