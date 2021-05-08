package core

import (
	"blockexchange/types"
	"errors"
	"io"
	"strconv"
)

type SchemaPartIterator func() (*types.SchemaPart, error)

func ExportSchema(w io.Writer, it SchemaPartIterator) error {
	// add header
	_, err := w.Write([]byte("5:return {"))
	if err != nil {
		return err
	}

	for {
		schemapart, err := it()
		if err != nil {
			return err
		}
		if schemapart == nil {
			// done
			break
		}
		exportSchemaPart(w, schemapart)
	}

	// add footer
	_, err = w.Write([]byte("}"))
	return err
}

func exportSchemaPart(w io.Writer, schemapart *types.SchemaPart) error {
	mapblock, err := ParseSchemaPart(schemapart)
	if err != nil {
		return err
	}

	// create reverse lookup table
	nodeid_names := make(map[int]string)
	for name, nodeid := range mapblock.Meta.NodeMapping {
		nodeid_names[nodeid] = name
	}

	for x := 0; x < mapblock.Meta.Size.X; x++ {
		for y := 0; y < mapblock.Meta.Size.Y; y++ {
			for z := 0; z < mapblock.Meta.Size.Z; z++ {
				exported, err := exportNode(w, x, y, z, mapblock, schemapart, nodeid_names)
				if err != nil {
					return err
				}
				if exported {
					// add delimiter
					_, err = w.Write([]byte(","))
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func exportNode(w io.Writer, x, y, z int, mapblock *ParsedSchemaPart, schemapart *types.SchemaPart, nodeid_names map[int]string) (bool, error) {
	// TODO: account for negative schema offset (size_n_minus)
	index := mapblock.GetIndex(x, y, z)
	nodeid := mapblock.NodeIDS[index]
	nodename := nodeid_names[int(nodeid)]
	if nodename == "" {
		return false, errors.New("Nodename not found for " + strconv.Itoa(int(nodeid)))
	} else if nodename == "air" {
		return false, nil
	}

	param1 := mapblock.Param1[index]
	param2 := mapblock.Param2[index]

	s := `{` +
		`["x"]=` + strconv.Itoa(x+schemapart.OffsetX) + `,` +
		`["y"]=` + strconv.Itoa(y+schemapart.OffsetY) + `,` +
		`["z"]=` + strconv.Itoa(z+schemapart.OffsetZ) + `,` +
		`["name"]="` + nodename + `"`

	if param1 != 0 {
		s += `,["param1"]=` + strconv.Itoa(int(param1))
	}
	if param2 != 0 {
		s += `,["param2"]=` + strconv.Itoa(int(param2))
	}

	s += `}`

	//TODO: export metadata

	_, err := w.Write([]byte(s))
	return true, err
}
