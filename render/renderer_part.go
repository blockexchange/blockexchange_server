package render

import (
	"blockexchange/types"
	"fmt"
	"math"
)

var tan30 = math.Tan(30 * math.Pi / 180)
var sqrt3div2 = 2 / math.Sqrt(3)

func (r *Renderer) RenderSchemaPart(schemapart *types.SchemaPart) error {
	mapblock, err := ParseSchemaPart(schemapart)
	if err != nil {
		return err
	}
	fmt.Println(mapblock)
	return nil
}
