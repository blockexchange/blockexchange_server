package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateUsername(t *testing.T) {
	assert.True(t, ValidateName("ok-name"))
	assert.True(t, ValidateName("ok_name"))
	assert.True(t, ValidateName("ok-name123"))
	assert.True(t, ValidateName("ok-name_456"))
	assert.True(t, ValidateName("ok-name--ABC"))
	assert.False(t, ValidateName("bad name"))
	assert.False(t, ValidateName("bad::"))
	assert.False(t, ValidateName("bad/"))
	assert.False(t, ValidateName("bad	"))
}
