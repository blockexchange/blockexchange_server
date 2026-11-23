package types

import (
	"errors"
	"regexp"
	"strings"
)

type SchemaModDependencyType string

const (
	SchemaModDependencyTypeMod  SchemaModDependencyType = "mod"
	SchemaModDependencyTypeGame SchemaModDependencyType = "game"
)

const SchemaModDependencyDefaultVersion = "latest"

// branch, tag or full commit hash, must not look like a command line option (passed to git)
var GitVersionRegex = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._/-]*$`)

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

// trims the fields and sets the default version, then checks the editable fields
func (d *SchemaModDependency) NormalizeAndValidate() error {
	d.GitURL = strings.TrimSpace(d.GitURL)
	d.Version = strings.TrimSpace(d.Version)
	if d.Version == "" {
		d.Version = SchemaModDependencyDefaultVersion
	}

	switch {
	case d.ModType != SchemaModDependencyTypeMod && d.ModType != SchemaModDependencyTypeGame:
		return errors.New("mod_type must be 'mod' or 'game'")
	case d.GitURL == "":
		return errors.New("git_url is required")
	case !strings.HasPrefix(d.GitURL, "https://") && !strings.HasPrefix(d.GitURL, "http://"):
		return errors.New("git_url must be a http(s) url")
	case len(d.GitURL) > 256:
		return errors.New("git_url is too long (max 256)")
	case len(d.Version) > 64:
		return errors.New("version is too long (max 64)")
	case !GitVersionRegex.MatchString(d.Version):
		return errors.New("version must be 'latest', a branch, a tag or a full commit hash")
	}
	return nil
}
