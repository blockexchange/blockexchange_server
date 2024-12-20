package types_test

import (
	"blockexchange/types"
	"testing"

	mt "github.com/minetest-go/types"

	"github.com/stretchr/testify/assert"
)

func TestParsePos(t *testing.T) {
	pos, err := types.ParsePos("1,2,3")
	assert.NoError(t, err)
	assert.Equal(t, mt.NewPos(1, 2, 3), pos)

	pos, err = types.ParsePos("-1,-2,-3")
	assert.NoError(t, err)
	assert.Equal(t, mt.NewPos(-1, -2, -3), pos)

	pos, err = types.ParsePos("0,0,0")
	assert.NoError(t, err)
	assert.Equal(t, mt.NewPos(0, 0, 0), pos)

	pos, err = types.ParsePos("1   ,   2 , 03")
	assert.NoError(t, err)
	assert.Equal(t, mt.NewPos(1, 2, 3), pos)

	pos, err = types.ParsePos("1,2,3,4")
	assert.Error(t, err)
	assert.Nil(t, pos)

	pos, err = types.ParsePos("")
	assert.Error(t, err)
	assert.Nil(t, pos)

	pos, err = types.ParsePos("1,2,")
	assert.Error(t, err)
	assert.Nil(t, pos)

	pos, err = types.ParsePos("garbage")
	assert.Error(t, err)
	assert.Nil(t, pos)

	pos, err = types.ParsePos("garbage,x,y")
	assert.Error(t, err)
	assert.Nil(t, pos)
}
