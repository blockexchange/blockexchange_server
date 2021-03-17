package render

import (
	"blockexchange/db"
	"blockexchange/types"
	"math"
)

type Renderer struct {
	SchemaPartRepo db.SchemaPartRepository
	Colormapping   map[string]*Color
}

func NewRenderer(spr db.SchemaPartRepository, cm map[string]*Color) *Renderer {
	return &Renderer{
		SchemaPartRepo: spr,
		Colormapping:   cm,
	}
}

func (r *Renderer) RenderSchema(schema *types.Schema) ([]byte, error) {
	//TODO: create canvas
	start_block_x := int(math.Ceil(float64(schema.MaxX)/16)) - 1
	start_block_z := int(math.Ceil(float64(schema.MaxZ)/16)) - 1
	end_block_y := int(math.Ceil(float64(schema.MaxY)/16)) - 1

	for block_x := start_block_x; block_x >= 0; block_x-- {
		for block_z := start_block_z; block_z >= 0; block_z-- {
			for block_y := 0; block_y <= end_block_y; block_y++ {
				x := block_x * 16
				y := block_y * 16
				z := block_z * 16

				schemapart, err := r.SchemaPartRepo.GetBySchemaIDAndOffset(0, x, y, z)
				if err != nil {
					return nil, err
				}

				if schemapart == nil {
					continue
				}

				mapblock, err := ParseSchemaPart(schemapart)
				if err != nil {
					return nil, err
				}

				pr := NewPartRenderer(schemapart, mapblock, r.Colormapping)
				//TODO: paint onto canvas
				//TODO: calculate canvas-offset
				pr.RenderSchemaPart()
			}
		}
	}

	return nil, nil
}
