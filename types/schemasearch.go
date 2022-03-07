package types

type OrderDirectionType string

const (
	ASC  OrderDirectionType = "asc"
	DESC OrderDirectionType = "desc"
)

type OrderColumnType string

const (
	CREATED OrderColumnType = "created"
)

type SchemaSearch struct {
	UserID         *float64            `json:"user_id"`
	SchemaID       *float64            `json:"schema_id"`
	TagID          *float64            `json:"tag_id"`
	SchemaName     *string             `json:"schema_name"`
	UserName       *string             `json:"user_name"`
	Keywords       *string             `json:"keywords"`
	Complete       *bool               `json:"complete"`
	OrderDirection *OrderDirectionType `json:"order_direction"`
	OrderColumn    *OrderColumnType    `json:"order_column"`
}

type SchemaSearchResult struct {
	Schema
	Stars int         `json:"stars"`
	User  *User       `json:"user"`
	Mods  []string    `json:"mods"`
	Tags  []SchemaTag `json:"tags"`
}
