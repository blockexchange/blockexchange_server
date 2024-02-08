package render

import (
	"blockexchange/parser"
	"blockexchange/types"
	"bytes"
	"math"

	"github.com/minetest-go/colormapping"

	"github.com/fogleman/gg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type ISORenderer struct {
	SchemaPartProvider SchemaPartProvider
	Colormapping       *colormapping.ColorMapping
}

type SchemaPartProvider func(schema_id int64, offset_x, offset_y, offset_z int) (*types.SchemaPart, error)

func NewISORenderer(spp SchemaPartProvider, cm *colormapping.ColorMapping) *ISORenderer {
	return &ISORenderer{
		SchemaPartProvider: spp,
		Colormapping:       cm,
	}
}

const DefaultImageSizeX = 800
const DefaultImageSizeY = 600

var renderHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
	Name:    "bx_renderschema_hist",
	Help:    "Histogram for the schema render time",
	Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10, 30, 60},
})

func (r *ISORenderer) RenderIsometricPreview(schema *types.Schema) ([]byte, error) {
	timer := prometheus.NewTimer(renderHistogram)
	defer timer.ObserveDuration()

	size_x := schema.SizeX
	size_y := schema.SizeY
	size_z := schema.SizeZ

	img_center_x := DefaultImageSizeX / (size_x + size_z) * size_z
	img_center_y := DefaultImageSizeY

	max_size := Max(size_x, Max(size_y, size_z))
	size := float64(DefaultImageSizeX) / float64(max_size) / 2.5

	dc := gg.NewContext(DefaultImageSizeX, DefaultImageSizeY)

	start_block_x := int(math.Ceil(float64(size_x)/16)) - 1
	start_block_z := int(math.Ceil(float64(size_z)/16)) - 1
	end_block_y := int(math.Ceil(float64(size_y)/16)) - 1

	for block_x := start_block_x; block_x >= 0; block_x-- {
		for block_z := start_block_z; block_z >= 0; block_z-- {
			for block_y := 0; block_y <= end_block_y; block_y++ {
				x := block_x * 16
				y := block_y * 16
				z := block_z * 16

				schemapart, err := r.SchemaPartProvider(*schema.ID, x, y, z)
				if err != nil {
					return nil, err
				}

				if schemapart == nil {
					continue
				}

				mapblock, err := parser.ParseSchemaPart(schemapart)
				if err != nil {
					return nil, err
				}
				x_offset := float64(img_center_x) + (size * float64(x)) - (size * float64(z))
				y_offset := float64(img_center_y) - (size * tan30 * float64(x)) - (size * tan30 * float64(z)) - (size * sqrt3div2 * float64(y))

				pr := NewISOPartRenderer(mapblock, r.Colormapping, size, x_offset, y_offset)
				pr.RenderSchemaPart(dc)
			}
		}
	}

	buf := bytes.NewBuffer([]byte{})
	err := dc.EncodePNG(buf)

	return buf.Bytes(), err
}
