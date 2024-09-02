package core_test

import (
	"blockexchange/core"
	"blockexchange/types"
	"fmt"
	"testing"

	mt "github.com/minetest-go/types"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSchemaClient(t *testing.T) {

	opts := &core.SchemaClientOpts{
		Pull: &types.SchematicPull{
			Hostname: "127.0.0.1",
			Port:     30000,
			PosX:     0,
			PosY:     0,
			PosZ:     0,
		},
		PullClient: &types.SchematicPullClient{
			Username: "test",
			Password: "test",
		},
		Schema: &types.Schema{
			UID:   uuid.NewString(),
			SizeX: 20,
			SizeY: 20,
			SizeZ: 20,
		},
		SetNode: func(pos *mt.Pos, node *mt.Node) error {
			fmt.Printf("would set node: %v @ %s\n", node, pos)
			return nil
		},
	}

	sc := core.NewSchemaClient(opts)
	err := sc.Run()
	assert.NoError(t, err)
	t.FailNow()

}
