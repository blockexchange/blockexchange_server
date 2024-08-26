package core_test

import (
	"blockexchange/core"
	"blockexchange/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchemaClient(t *testing.T) {

	p := &types.SchematicPull{
		Hostname: "127.0.0.1",
		Port:     30000,
	}
	pc := &types.SchematicPullClient{
		Username: "test",
		Password: "test",
	}
	err := core.SchemaClient(p, pc)
	assert.NoError(t, err)
	t.FailNow()

}
