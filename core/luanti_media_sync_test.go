package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGuessMediaModName(t *testing.T) {
	modnames := map[string]bool{"default": true, "mesecons": true, "mesecons_wires": true}
	assert.Equal(t, "default", guessMediaModName("default_stone.png", modnames))
	assert.Equal(t, "mesecons_wires", guessMediaModName("mesecons_wires_on.png", modnames))
	assert.Equal(t, "mesecons", guessMediaModName("mesecons_x.png", modnames))
	assert.Equal(t, "default", guessMediaModName("default.png", modnames))
	assert.Equal(t, "", guessMediaModName("unknown_texture.png", modnames))
	assert.Equal(t, "", guessMediaModName("defaultx_stone.png", modnames))
}
