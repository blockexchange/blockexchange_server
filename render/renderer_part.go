package render

import (
	"blockexchange/types"
	"sort"
)

type Block struct {
	X     int
	Y     int
	Z     int
	Color *Color
	Order int
}

type PartRenderer struct {
	Schemapart          *types.SchemaPart
	Mapblock            *ParsedSchemaPart
	Colormapping        map[string]*Color
	NodeIDStringMapping map[int]string
	Blocks              []Block
	MaxX                int
	MaxY                int
	MaxZ                int
	YMultiplier         int
	XMultiplier         int
}

func NewPartRenderer(schemapart *types.SchemaPart, mapblock *ParsedSchemaPart, cm map[string]*Color) *PartRenderer {
	// reverse index
	idm := make(map[int]string)
	for k, v := range mapblock.Meta.NodeMapping {
		idm[v] = k
	}
	return &PartRenderer{
		Schemapart:          schemapart,
		Mapblock:            mapblock,
		Blocks:              make([]Block, 0),
		Colormapping:        cm,
		NodeIDStringMapping: idm,
		MaxX:                mapblock.Meta.Size.X - 1,
		MaxY:                mapblock.Meta.Size.Y - 1,
		MaxZ:                mapblock.Meta.Size.Z - 1,
		YMultiplier:         mapblock.Meta.Size.Z,
		XMultiplier:         mapblock.Meta.Size.Y * mapblock.Meta.Size.Z,
	}
}

func (r *PartRenderer) GetColorAtPos(x, y, z int) *Color {
	if x > r.MaxX || y > r.MaxY || z > r.MaxZ || x < 0 || y < 0 || z < 0 {
		return nil
	}

	index := z + (y * r.YMultiplier) + (x * r.XMultiplier)
	nodeid := int(r.Mapblock.NodeIDS[index])
	nodename := r.NodeIDStringMapping[nodeid]
	color := r.Colormapping[nodename]
	return color
}

func (r *PartRenderer) ProbePosition(x, y, z int) {
	color := r.GetColorAtPos(x, y, z)
	if color != nil {
		block := Block{
			X:     x,
			Y:     y,
			Z:     z,
			Color: color,
			Order: y + ((r.MaxX - x) * r.MaxX) + ((r.MaxZ - z) + r.MaxZ),
		}

		r.Blocks = append(r.Blocks, block)
		return
	}

	next_x := x + 1
	next_y := y + 1
	next_z := z + 1

	if next_x > r.MaxX || next_z > r.MaxZ || next_y < 0 {
		// mapblock ends
		return
	}

	r.ProbePosition(next_x, next_y, next_z)
}

func (r *PartRenderer) RenderSchemaPart() error {

	for y := 0; y < r.MaxY; y++ {
		// right side
		for x := r.MaxX; x >= 1; x-- {
			r.ProbePosition(x, y, 0)
		}

		// left side
		for z := r.MaxZ; z >= 0; z-- {
			r.ProbePosition(0, y, z)
		}
	}

	// top side
	for z := r.MaxZ; z >= 0; z-- {
		for x := r.MaxX; x >= 0; x-- {
			r.ProbePosition(x, r.MaxY, z)
		}
	}

	sort.Slice(r.Blocks, func(i int, j int) bool {
		return r.Blocks[i].Order > r.Blocks[j].Order
	})

	//TODO: draw isometric rectangles

	return nil
}
