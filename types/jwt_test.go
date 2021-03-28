package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenInfoString(t *testing.T) {
	token := TokenInfo{
		Username: "user",
	}

	assert.NotNil(t, token.String())
}
