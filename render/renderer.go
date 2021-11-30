package render

import (
	"blockexchange/core"
	"blockexchange/types"
	"bytes"
	"math"

	"github.com/fogleman/gg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Renderer struct {
	SchemaPartProvider SchemaPartProvider
	Colormapping       map[string]*Color
}

type SchemaPartProvider func(schema_id int64, offset_x, offset_y, offset_z int) (*types.SchemaPart, error)

func NewRenderer(spp SchemaPartProvider, cm map[string]*Color) *Renderer {
	return &Renderer{
		SchemaPartProvider: spp,
		Colormapping:       cm,
	}
}

const img_size_x = 800
const img_size_y = 600

var renderHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
	Name:    "bx_renderschema_hist",
	Help:    "Histogram for the schema render time",
	Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10, 30, 60},
})

func (r *Renderer) RenderSchema(schema *types.Schema) ([]byte, error) {
	timer := prometheus.NewTimer(renderHistogram)
	defer timer.ObserveDuration()

	size_x := schema.SizeX
	size_y := schema.SizeY
	size_z := schema.SizeZ

	img_center_x := img_size_x / (size_x + size_z) * size_z
	img_center_y := img_size_y

	max_size := Max(size_x, Max(size_y, size_z))
	size := float64(img_size_x) / float64(max_size) / 2.5

	dc := gg.NewContext(img_size_x, img_size_y)

	start_block_x := int(math.Ceil(float64(size_x)/16)) - 1
	start_block_z := int(math.Ceil(float64(size_z)/16)) - 1
	end_block_y := int(math.Ceil(float64(size_y)/16)) - 1

	for block_x := start_block_x; block_x >= 0; block_x-- {
		for block_z := start_block_z; block_z >= 0; block_z-- {
			for block_y := 0; block_y <= end_block_y; block_y++ {
				x := block_x * 16
				y := block_y * 16
				z := block_z * 16

				schemapart, err := r.SchemaPartProvider(schema.ID, x, y, z)
				if err != nil {
					return nil, err
				}

				if schemapart == nil {
					continue
				}

				mapblock, err := core.ParseSchemaPart(schemapart)
				if err != nil {
					return nil, err
				}
				x_offset := float64(img_center_x) + (size * float64(x)) - (size * float64(z))
				y_offset := float64(img_center_y) - (size * tan30 * float64(x)) - (size * tan30 * float64(z)) - (size * sqrt3div2 * float64(y))

				pr := NewPartRenderer(mapblock, r.Colormapping, size, x_offset, y_offset)
				pr.RenderSchemaPart(dc)
			}
		}
	}

	buf := bytes.NewBuffer([]byte{})
	err := dc.EncodePNG(buf)

	return buf.Bytes(), err
}
