package core

import (
	"blockexchange/types"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderFeedTemplate(t *testing.T) {
	schema := types.Schema{
		Name:       "my_schema",
		License:    "CC0",
		SizeX:      100,
		SizeY:      200,
		SizeZ:      300,
		TotalSize:  2300,
		TotalParts: 16,
	}
	user := types.User{
		Name: "Somebody",
	}

	buf, err := renderFeedTemplate("http://example.com", &schema, &user)
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	fmt.Print(buf.String())
}
