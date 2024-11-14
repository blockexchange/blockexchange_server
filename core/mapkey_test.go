package core_test

import (
	"blockexchange/core"
	"testing"

	mt "github.com/minetest-go/types"
	"github.com/stretchr/testify/assert"
)

func TestMapKey(t *testing.T) {
	p := mt.NewPos(-10, 20, -40)
	i := core.MapKey(p)
	p1 := core.ParseMapKey(i)
	assert.Equal(t, p, p1)
}
