package render

import (
	"blockexchange/db"
	"blockexchange/types"
	"bytes"
	"math"

	"github.com/fogleman/gg"
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

const img_size_x = 800
const img_size_y = 600

func (r *Renderer) RenderSchema(schema *types.Schema) ([]byte, error) {

	img_center_x := img_size_x / (schema.SizeZ + schema.SizeX) * schema.SizeZ
	img_center_y := img_size_y

	max_size := Max(schema.SizeX, Max(schema.SizeY, schema.SizeZ))
	size := float64(img_size_x) / float64(max_size) / 2.5

	dc := gg.NewContext(img_size_x, img_size_y)

	start_block_x := int(math.Ceil(float64(schema.SizeX)/16)) - 1
	start_block_z := int(math.Ceil(float64(schema.SizeZ)/16)) - 1
	end_block_y := int(math.Ceil(float64(schema.SizeY)/16)) - 1

	for block_x := start_block_x; block_x >= 0; block_x-- {
		for block_z := start_block_z; block_z >= 0; block_z-- {
			for block_y := 0; block_y <= end_block_y; block_y++ {
				x := block_x * 16
				y := block_y * 16
				z := block_z * 16

				schemapart, err := r.SchemaPartRepo.GetBySchemaIDAndOffset(schema.ID, x, y, z)
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
				x_offset := float64(img_center_x) + (size * float64(x)) - (size * float64(z))
				y_offset := float64(img_center_y) - (size * tan30 * float64(x)) - (size * tan30 * float64(z)) - (size * sqrt3div2 * float64(y))

				pr := NewPartRenderer(schemapart, mapblock, r.Colormapping, size, x_offset, y_offset)
				pr.RenderSchemaPart(dc)
			}
		}
	}

	buf := bytes.NewBuffer([]byte{})
	err := dc.EncodePNG(buf)

	return buf.Bytes(), err
}
