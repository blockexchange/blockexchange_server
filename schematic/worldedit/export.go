package worldedit

import (
	"blockexchange/parser"
	"blockexchange/types"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Exporter struct {
	w          io.Writer
	intialized bool
}

func NewExporter(w io.Writer) *Exporter {
	return &Exporter{w: w, intialized: false}
}

func (e *Exporter) Export(schemapart *types.SchemaPart) error {
	if !e.intialized {
		_, err := e.w.Write([]byte("5:return {"))
		if err != nil {
			return err
		}
		e.intialized = true
	}

	return exportSchemaPart(e.w, schemapart)
}

func (e *Exporter) Close() error {
	// add footer
	_, err := e.w.Write([]byte("}"))
	return err
}

func exportSchemaPart(w io.Writer, schemapart *types.SchemaPart) error {
	mapblock, err := parser.ParseSchemaPart(schemapart)
	if err != nil {
		return err
	}

	// create reverse lookup table
	nodeid_names := make(map[int16]string)
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

func exportNode(w io.Writer, x, y, z int, mapblock *parser.ParsedSchemaPart, schemapart *types.SchemaPart, nodeid_names map[int16]string) (bool, error) {
	index := mapblock.GetIndex(x, y, z)
	nodeid := mapblock.NodeIDS[index]
	if nodeid == mapblock.Meta.NodeMapping["air"] {
		return false, nil
	}

	nodename := nodeid_names[nodeid]
	if nodename == "" {
		return false, errors.New("Nodename not found for " + strconv.Itoa(int(nodeid)))
	}

	param1 := mapblock.Param1[index]
	param2 := mapblock.Param2[index]

	parts := []string{
		fmt.Sprintf("x=%d", x+schemapart.OffsetX),
		fmt.Sprintf("y=%d", y+schemapart.OffsetY),
		fmt.Sprintf("z=%d", z+schemapart.OffsetZ),
		fmt.Sprintf("name=\"%s\"", nodename),
	}

	if param1 != 0 {
		parts = append(parts, fmt.Sprintf("param1=%d", param1))
	}

	if param2 != 0 {
		parts = append(parts, fmt.Sprintf("param2=%d", param2))
	}

	if mapblock.Meta != nil && mapblock.Meta.Metadata != nil {
		fieldsStr := ""
		key := mapblock.Meta.Metadata.GetKey(x, y, z)
		if mapblock.Meta.Metadata.Meta != nil && mapblock.Meta.Metadata.Meta[key] != nil {
			fields := mapblock.Meta.Metadata.Meta[key].Fields
			fieldparts := make([]string, 0)
			for k, v := range fields {
				rv := strings.ReplaceAll(v, "\\", "\\\\")
				rv = strings.ReplaceAll(rv, "\"", "\\\"")
				rv = strings.ReplaceAll(rv, "\n", "\\n")
				fieldparts = append(fieldparts, fmt.Sprintf("[\"%s\"]=\"%s\"", k, rv))
			}
			fieldsStr = strings.Join(fieldparts, ",")

			inventories := mapblock.Meta.Metadata.Meta[key].Inventories
			invparts := make([]string, 0)
			for invname, inventry := range inventories {
				stacks := make([]string, len(inventry))
				for i, stack := range inventry {
					rv := strings.ReplaceAll(stack, "\\", "\\\\")
					rv = strings.ReplaceAll(rv, "\"", "\\\"")
					rv = strings.ReplaceAll(rv, "\n", "\\n")

					stacks[i] = fmt.Sprintf("\"%s\"", rv)
				}

				invparts = append(invparts, fmt.Sprintf("[\"%s\"]={%s}", invname, strings.Join(stacks, ",")))
			}
			invStr := strings.Join(invparts, ",")

			if fieldsStr != "" || invStr != "" {
				parts = append(parts, fmt.Sprintf("meta={fields={%s},inventory={%s}}", fieldsStr, invStr))
			}
		}

	}

	str := fmt.Sprintf("{%s}", strings.Join(parts, ","))
	_, err := w.Write([]byte(str))
	return true, err
}
