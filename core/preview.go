package core

import (
	"blockexchange/parser"
	"blockexchange/types"
	"bytes"
	"fmt"
	"image/png"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/minetest-go/colormapping"
	"github.com/minetest-go/maprenderer"
	mtypes "github.com/minetest-go/types"
)

var cm *colormapping.ColorMapping = colormapping.NewColorMapping()
var zeroPos = mtypes.NewPos(0, 0, 0)
var airOnlyMapblock = &parser.ParsedSchemaPart{}

func init() {
	err := cm.LoadDefaults()
	if err != nil {
		panic(err)
	}
}

func (c *Core) getNodeAccessor(schema *types.Schema) mtypes.NodeAccessor {
	cache := expirable.NewLRU[int64, *parser.ParsedSchemaPart](1000, nil, time.Second*10)

	return func(p *mtypes.Pos) (*mtypes.Node, error) {
		po := p.Divide(16).Multiply(16)

		key := int64(po[0]) +
			int64(po[1])<<16 +
			int64(po[2])<<32

		mapblock, found := cache.Get(key)
		if found && mapblock == airOnlyMapblock {
			return nil, nil
		}

		if !found {
			schemapart, err := c.repos.SchemaPartRepo.GetBySchemaUIDAndOffset(schema.UID, po[0], po[1], po[2])
			if err != nil {
				return nil, fmt.Errorf("get schemapart error @ %s: %v", p, err)
			}

			if schemapart != nil {
				mapblock, err = parser.ParseSchemaPart(schemapart)
				if err != nil {
					return nil, fmt.Errorf("parse error @ %s: %v", p, err)
				}

				if len(mapblock.Meta.NodeMapping) == 1 && mapblock.Meta.NodeMapping["air"] > 0 {
					// mark as air-only
					cache.Add(key, airOnlyMapblock)
					return nil, nil
				}

				cache.Add(key, mapblock)
			} else {
				// not found, mark as air-only
				cache.Add(key, airOnlyMapblock)
				return nil, nil
			}
		}

		rel_pos := p.Subtract(po)
		index := mapblock.GetIndex(rel_pos.X(), rel_pos.Y(), rel_pos.Z())
		if index >= len(mapblock.NodeIDS) {
			return nil, fmt.Errorf("index mismatch: got %d, length: %d, rel_pos: %s, abs_pos: %s", index, len(mapblock.NodeIDS), rel_pos, p)
		}
		nodeid := mapblock.NodeIDS[index]

		return &mtypes.Node{
			Name:   mapblock.NodeNameLookup[nodeid],
			Param1: int(mapblock.Param1[index]),
			Param2: int(mapblock.Param2[index]),
			Pos:    p,
		}, nil
	}
}

func (c *Core) UpdatePreview(schema *types.Schema) (*types.SchemaScreenshot, error) {
	max_pos := mtypes.NewPos(schema.SizeX-1, schema.SizeY-1, schema.SizeZ-1)
	max_axis := max(max_pos.X(), max_pos.Y(), max_pos.Z())

	opts := maprenderer.NewDefaultIsoRenderOpts()
	if max_axis < 10 {
		opts.CubeLen = 48
	} else if max_axis < 50 {
		opts.CubeLen = 16
	}

	img, err := maprenderer.RenderIsometric(c.getNodeAccessor(schema), cm.GetColor, zeroPos, max_pos, opts)
	if err != nil {
		return nil, fmt.Errorf("render error: %v", err)
	}

	buf := bytes.NewBuffer([]byte{})
	err = png.Encode(buf, img)
	if err != nil {
		return nil, fmt.Errorf("encode error: %v", err)
	}

	screenshots, err := c.repos.SchemaScreenshotRepo.GetBySchemaUID(schema.UID)
	if err != nil {
		return nil, err
	}

	var screenshot *types.SchemaScreenshot

	if len(screenshots) >= 1 {
		// update existing
		screenshot = screenshots[0]
		screenshot.Data = buf.Bytes()

		err = c.repos.SchemaScreenshotRepo.Update(screenshot)
		if err != nil {
			return nil, err
		}
	} else {
		// create a new one
		screenshot = &types.SchemaScreenshot{
			SchemaUID: schema.UID,
			Type:      "image/png",
			Title:     "Isometric preview",
			Data:      buf.Bytes(),
		}

		err = c.repos.SchemaScreenshotRepo.Create(screenshot)
		if err != nil {
			return nil, err
		}
	}

	return screenshot, nil
}
