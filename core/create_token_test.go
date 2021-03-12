package core

import (
	"testing"
)

func TestCreateToken(t *testing.T) {
	s := CreateToken(6)
	if len(s) != 6 {
		t.Fatal("wrong length!")
	}
}
