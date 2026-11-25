package mtschem_test

import (
	"blockexchange/schematic/mtschem"
	"blockexchange/types"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportImportRoundtrip(t *testing.T) {
	for _, tf := range testfiles {
		t.Run(tf.filename, func(t *testing.T) {
			s := parseFile(t, tf.filename)

			// mts -> blockexchange parts
			parsedParts, err := mtschem.Import(s)
			require.NoError(t, err)
			parts := make([]*types.SchemaPart, len(parsedParts))
			for i, pp := range parsedParts {
				parts[i], err = pp.Convert()
				require.NoError(t, err)
			}

			// blockexchange parts -> mts
			buf := bytes.NewBuffer([]byte{})
			e, err := mtschem.NewExporter(buf, s.SizeX, s.SizeY, s.SizeZ)
			require.NoError(t, err)
			for _, part := range parts {
				require.NoError(t, e.Export(part))
			}
			require.NoError(t, e.Close())

			s2, err := mtschem.Parse(buf.Bytes())
			require.NoError(t, err)
			assert.Equal(t, s.SizeX, s2.SizeX)
			assert.Equal(t, s.SizeY, s2.SizeY)
			assert.Equal(t, s.SizeZ, s2.SizeZ)
			assert.Equal(t, "air", s2.NodeNames[0])

			for x := 0; x < s.SizeX; x++ {
				for y := 0; y < s.SizeY; y++ {
					for z := 0; z < s.SizeZ; z++ {
						name := placedName(s, x, y, z)
						index := s2.Index(x, y, z)
						require.Equal(t, name, s2.NodeNames[s2.Content[index]], "node at %d/%d/%d", x, y, z)
						require.Equal(t, mtschem.ProbAlways, s2.Param1[index])
						if name != "air" {
							require.Equal(t, s.Param2[index], s2.Param2[index])
						}
					}
				}
			}

			// re-import yields the same parts
			parsedParts2, err := mtschem.Import(s2)
			require.NoError(t, err)
			assert.Equal(t, len(parsedParts), len(parsedParts2))
			assertPartsMatch(t, s2, parsedParts2)
		})
	}
}

func TestExportOutOfBounds(t *testing.T) {
	s := parseFile(t, "testdata/Statue_of_Liberty(New_York,USA).mts")
	parsedParts, err := mtschem.Import(s)
	require.NoError(t, err)
	part, err := parsedParts[0].Convert()
	require.NoError(t, err)

	e, err := mtschem.NewExporter(bytes.NewBuffer([]byte{}), 1, 1, 1)
	require.NoError(t, err)
	assert.Error(t, e.Export(part))
}

func TestNewExporterInvalidSize(t *testing.T) {
	_, err := mtschem.NewExporter(bytes.NewBuffer([]byte{}), 0, 1, 1)
	assert.Error(t, err)
	_, err = mtschem.NewExporter(bytes.NewBuffer([]byte{}), 1, 501, 1)
	assert.Error(t, err)
}
