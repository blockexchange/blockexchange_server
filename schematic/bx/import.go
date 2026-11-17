package bx

import (
	"archive/zip"
	"blockexchange/types"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"regexp"
)

const (
	// MaxSize is the max side-length of an imported schema
	MaxSize = 500
	// MaxFileSize is the max decompressed size of a single file in the archive
	MaxFileSize = 16 * 1024 * 1024
	// MaxTotalSize is the max decompressed size of all files in the archive (zip-bomb protection)
	MaxTotalSize = 512 * 1024 * 1024
	// MaxMods is the max number of entries in mods.json
	MaxMods = 1024
)

// ErrInvalidBX is wrapped by all errors caused by invalid archive content
var ErrInvalidBX = errors.New("invalid blockexchange archive")

var modnameRegex = regexp.MustCompile(`^[a-z0-9_]+$`)

type BXImport struct {
	Parts  []*types.SchemaPart
	Schema *types.Schema
	Mods   []string
}

func ParseBXContent(data []byte) (*BXImport, error) {
	res, err := parseBXContent(data)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidBX, err)
	}
	return res, nil
}

func parseBXContent(data []byte) (*BXImport, error) {
	res := &BXImport{
		Parts:  make([]*types.SchemaPart, 0),
		Schema: &types.Schema{},
		Mods:   make([]string, 0),
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}
	u := &unzipper{zr: zr, budget: MaxTotalSize}

	err = u.unmarshal("schema.json", res.Schema)
	if err != nil {
		return nil, err
	}
	s := res.Schema
	if s.SizeX < 1 || s.SizeY < 1 || s.SizeZ < 1 || s.SizeX > MaxSize || s.SizeY > MaxSize || s.SizeZ > MaxSize {
		return nil, fmt.Errorf("schema size %d/%d/%d out of range [1, %d]", s.SizeX, s.SizeY, s.SizeZ, MaxSize)
	}

	err = u.unmarshal("mods.json", &res.Mods)
	if err != nil {
		return nil, err
	}
	if len(res.Mods) > MaxMods {
		return nil, fmt.Errorf("more than %d mods", MaxMods)
	}
	for _, modname := range res.Mods {
		if len(modname) > 256 || !modnameRegex.MatchString(modname) {
			return nil, fmt.Errorf("invalid mod name %q", modname)
		}
	}

	mb_max_x := (s.SizeX - 1) / 16
	mb_max_y := (s.SizeY - 1) / 16
	mb_max_z := (s.SizeZ - 1) / 16
	for x := 0; x <= mb_max_x; x++ {
		for y := 0; y <= mb_max_y; y++ {
			for z := 0; z <= mb_max_z; z++ {
				sp := &types.SchemaPart{
					OffsetX: x * 16,
					OffsetY: y * 16,
					OffsetZ: z * 16,
				}
				filename := formatSchemapartFilename(sp)
				err = u.unmarshal(filename, sp)
				if errors.Is(err, fs.ErrNotExist) {
					// air-only mapblock
					continue
				}
				if err != nil {
					// other error
					return nil, err
				}

				// the file content must not move the part to another position
				if sp.OffsetX != x*16 || sp.OffsetY != y*16 || sp.OffsetZ != z*16 {
					return nil, fmt.Errorf("%s: offset %d/%d/%d does not match filename", filename, sp.OffsetX, sp.OffsetY, sp.OffsetZ)
				}

				res.Parts = append(res.Parts, sp)
			}
		}
	}

	return res, nil
}

// unzipper reads json files from a zip archive within a total decompressed size budget
type unzipper struct {
	zr     *zip.Reader
	budget int64
}

func (u *unzipper) unmarshal(name string, obj any) error {
	f, err := u.zr.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	limit := min(int64(MaxFileSize), u.budget)
	// read one byte more than allowed to detect oversized files
	data, err := io.ReadAll(io.LimitReader(f, limit+1))
	if err != nil {
		return fmt.Errorf("%s: %v", name, err)
	}
	if int64(len(data)) > limit {
		return fmt.Errorf("%s: decompressed size exceeds limit", name)
	}
	u.budget -= int64(len(data))

	if err = json.Unmarshal(data, obj); err != nil {
		return fmt.Errorf("%s: %v", name, err)
	}
	return nil
}
