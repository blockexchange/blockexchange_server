package types

type SchemaModDependencyType string

const (
	SchemaModDependencyTypeMod  SchemaModDependencyType = "mod"
	SchemaModDependencyTypeGame SchemaModDependencyType = "game"
)

const SchemaModDependencyDefaultVersion = "latest"

type SchemaModDependency struct {
	UID       string                  `json:"uid" gorm:"primarykey;column:uid"`
	SchemaUID string                  `json:"schema_uid" gorm:"column:schema_uid"`
	ModType   SchemaModDependencyType `json:"mod_type" gorm:"column:mod_type"`
	GitURL    string                  `json:"git_url" gorm:"column:git_url"`
	Version   string                  `json:"version" gorm:"column:version"`
}

func (d *SchemaModDependency) TableName() string {
	return "schema_mod_dependency"
}
