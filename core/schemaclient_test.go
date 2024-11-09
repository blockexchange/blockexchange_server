package core_test

import (
	"blockexchange/core"
	"blockexchange/parser"
	"blockexchange/types"
	"fmt"
	"testing"

	mt "github.com/minetest-go/types"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSchemaClient(t *testing.T) {
	t.SkipNow()

	opts := &core.SchemaClientOpts{
		Pull: &types.SchematicPull{
			Hostname: "127.0.0.1",
			Port:     30000,
			PosX:     0,
			PosY:     100,
			PosZ:     0,
		},
		PullClient: &types.SchematicPullClient{
			Username: "test",
			Password: "test",
		},
		Schema: &types.Schema{
			UID:   uuid.NewString(),
			SizeX: 20,
			SizeY: 2,
			SizeZ: 2,
		},
		SetNode: func(pos *mt.Pos, node *mt.Node) error {
			if node.Name != "air" {
				fmt.Printf("would set node: %v @ %s\n", node, pos)
			}
			return nil
		},
		SetMeta: func(pos *mt.Pos, md *parser.MetadataEntry) error {
			if md.Fields != nil || md.Inventories != nil {
				fmt.Printf("would set metadata: %v @ %s\n", md, pos)
			}
			return nil
		},
	}

	sc := core.NewSchemaClient(opts)
	err := sc.Run()
	assert.NoError(t, err)
	t.FailNow()

}
