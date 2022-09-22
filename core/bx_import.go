package core

import (
	"archive/zip"
	"blockexchange/types"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
)

type BXImport struct {
	Parts  []*types.SchemaPart
	Schema *types.Schema
	Mods   []string
}

func ImportBXSchema(r io.ReaderAt, size int64) (*BXImport, error) {
	res := &BXImport{
		Parts:  make([]*types.SchemaPart, 0),
		Schema: &types.Schema{},
		Mods:   make([]string, 0),
	}

	zr, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}

	err = unmarshalFromZip(zr, "schema.json", res.Schema)
	if err != nil {
		return nil, err
	}

	err = unmarshalFromZip(zr, "mods.json", &res.Mods)
	if err != nil {
		return nil, err
	}

	mb_max_x := int(res.Schema.SizeX / 16)
	mb_max_y := int(res.Schema.SizeY / 16)
	mb_max_z := int(res.Schema.SizeZ / 16)
	for x := 0; x <= mb_max_x; x++ {
		for y := 0; y <= mb_max_y; y++ {
			for z := 0; z <= mb_max_z; z++ {
				sp := &types.SchemaPart{
					OffsetX: x * 16,
					OffsetY: y * 16,
					OffsetZ: z * 16,
				}
				err = unmarshalFromZip(zr, formatSchemapartFilename(sp), sp)
				if errors.Is(err, fs.ErrNotExist) {
					// air-only mapblock
					continue
				}
				if err != nil {
					// other error
					return nil, err
				}

				res.Parts = append(res.Parts, sp)
			}
		}
	}

	return res, nil
}

func unmarshalFromZip(zr *zip.Reader, name string, obj any) error {
	f, err := zr.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(obj)
}
