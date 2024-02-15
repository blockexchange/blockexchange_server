package core

import (
	"blockexchange/parser"
	"blockexchange/types"
	"strings"
)

func (c *Core) ExtractModnames(schema_uid string) ([]string, error) {

	modname_map := map[string]bool{}

	err := c.SchemapartCallback(schema_uid, func(sp *types.SchemaPart) error {
		psp, err := parser.ParseSchemaPart(sp)
		if err != nil {
			return err
		}

		for name := range psp.Meta.NodeMapping {
			nameparts := strings.Split(name, ":")
			if len(nameparts) == 2 {
				modname_map[nameparts[0]] = true
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// convert to list
	modnames := []string{}
	for n := range modname_map {
		modnames = append(modnames, n)
	}
	return modnames, nil
}
