package core_test

import (
	"blockexchange/testutils"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWEImport(t *testing.T) {
	c, repos := getTestCore(t)

	u := testutils.CreateUser(repos.UserRepo, t, nil)
	assert.NotNil(t, u)

	data, err := os.ReadFile("testdata/plain_chest.we")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	s, err := c.ImportWE(data, u.Name, "we_import")
	assert.NoError(t, err)
	assert.NotNil(t, s)
}

func TestMTSImport(t *testing.T) {
	c, repos := getTestCore(t)

	u := testutils.CreateUser(repos.UserRepo, t, nil)
	assert.NotNil(t, u)

	data, err := os.ReadFile("../schematic/mtschem/testdata/Statue_of_Liberty(New_York,USA).mts")
	assert.NoError(t, err)

	s, err := c.ImportMTS(data, u.Name, "mts_import")
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, "mts_import", s.Name)
	assert.Equal(t, 44, s.SizeX)
	assert.Equal(t, 130, s.SizeY)
	assert.Equal(t, 42, s.SizeZ)
	assert.True(t, s.TotalParts > 0)

	mods, err := repos.SchemaModRepo.GetSchemaModsBySchemaUID(s.UID)
	assert.NoError(t, err)
	modnames := []string{}
	for _, m := range mods {
		modnames = append(modnames, m.ModName)
	}
	assert.Contains(t, modnames, "default")

	_, err = c.ImportMTS([]byte("invalid"), u.Name, "mts_import")
	assert.Error(t, err)
}
