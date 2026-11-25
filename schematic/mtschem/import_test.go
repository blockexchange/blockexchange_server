package mtschem_test

import (
	"blockexchange/parser"
	"blockexchange/schematic/mtschem"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// returns the effectively placed node name at the given position ("air" for never-placed nodes)
func placedName(s *mtschem.Schematic, x, y, z int) string {
	index := s.Index(x, y, z)
	if s.Param1[index]&mtschem.ProbMask == mtschem.ProbNever {
		return "air"
	}
	return s.NodeNames[s.Content[index]]
}

// asserts that the parts contain exactly the placed nodes of the schematic
func assertPartsMatch(t *testing.T, s *mtschem.Schematic, parts []*parser.ParsedSchemaPart) {
	nodeCount := 0
	for _, part := range parts {
		lookup := map[int16]string{}
		for name, id := range part.Meta.NodeMapping {
			lookup[id] = name
		}

		for x := 0; x < part.Meta.Size.X; x++ {
			for y := 0; y < part.Meta.Size.Y; y++ {
				for z := 0; z < part.Meta.Size.Z; z++ {
					sx, sy, sz := part.PosX*16+x, part.PosY*16+y, part.PosZ*16+z
					partIndex := part.GetIndex(x, y, z)
					name := lookup[part.NodeIDS[partIndex]]
					require.Equal(t, placedName(s, sx, sy, sz), name, "node at %d/%d/%d", sx, sy, sz)
					if name != "air" {
						nodeCount++
						require.Equal(t, s.Param2[s.Index(sx, sy, sz)], part.Param2[partIndex])
					}
				}
			}
		}
	}

	expectedCount := 0
	for x := 0; x < s.SizeX; x++ {
		for y := 0; y < s.SizeY; y++ {
			for z := 0; z < s.SizeZ; z++ {
				if placedName(s, x, y, z) != "air" {
					expectedCount++
				}
			}
		}
	}
	assert.Equal(t, expectedCount, nodeCount)
}

func TestImport(t *testing.T) {
	for _, tf := range testfiles {
		t.Run(tf.filename, func(t *testing.T) {
			s := parseFile(t, tf.filename)
			parts, err := mtschem.Import(s)
			require.NoError(t, err)
			assert.NotEmpty(t, parts)

			maxParts := ((s.SizeX + 15) / 16) * ((s.SizeY + 15) / 16) * ((s.SizeZ + 15) / 16)
			assert.LessOrEqual(t, len(parts), maxParts)

			for _, part := range parts {
				assert.Equal(t, int16(0), part.Meta.NodeMapping["air"])
				assert.True(t, part.Meta.Size.X >= 1 && part.Meta.Size.X <= 16)
				assert.True(t, part.Meta.Size.Y >= 1 && part.Meta.Size.Y <= 16)
				assert.True(t, part.Meta.Size.Z >= 1 && part.Meta.Size.Z <= 16)
				for _, p1 := range part.Param1 {
					assert.Equal(t, byte(0), p1)
				}
			}

			assertPartsMatch(t, s, parts)
		})
	}
}

func TestImportConvertible(t *testing.T) {
	s := parseFile(t, "testdata/Statue_of_Liberty(New_York,USA).mts")
	parts, err := mtschem.Import(s)
	require.NoError(t, err)

	for _, part := range parts {
		sp, err := part.Convert()
		require.NoError(t, err)

		parsed, err := parser.ParseSchemaPart(sp)
		require.NoError(t, err)
		assert.Equal(t, part.NodeIDS, parsed.NodeIDS)
		assert.Equal(t, part.Param2, parsed.Param2)
		assert.Equal(t, part.PosX, parsed.PosX)
		assert.Equal(t, part.PosY, parsed.PosY)
		assert.Equal(t, part.PosZ, parsed.PosZ)
	}
}

func TestImportNeverPlaced(t *testing.T) {
	data := buildRaw(t, 4, 3, 1, 1, []byte{0x7F},
		[]string{"air", "default:stone"},
		[]uint16{1, 1, 0}, []byte{mtschem.ProbNever, mtschem.ProbAlways | mtschem.ForcePlace, mtschem.ProbAlways}, []byte{0, 2, 0})
	s, err := mtschem.Parse(data)
	require.NoError(t, err)

	parts, err := mtschem.Import(s)
	require.NoError(t, err)
	require.Equal(t, 1, len(parts))
	part := parts[0]
	assert.Equal(t, []int16{0, 1, 0}, part.NodeIDS)
	assert.Equal(t, []byte{0, 2, 0}, part.Param2)
}

func TestImportEmpty(t *testing.T) {
	s := mtschem.NewSchematic(20, 20, 20)
	parts, err := mtschem.Import(s)
	require.NoError(t, err)
	assert.Empty(t, parts)

	// a single node only creates one part
	s.NodeNames = append(s.NodeNames, "default:stone")
	s.Content[s.Index(17, 0, 3)] = 1
	parts, err = mtschem.Import(s)
	require.NoError(t, err)
	require.Equal(t, 1, len(parts))
	assert.Equal(t, 1, parts[0].PosX)
	assert.Equal(t, parser.SchemaPartSize{X: 4, Y: 16, Z: 16}, parts[0].Meta.Size)

	buf := bytes.NewBuffer([]byte{})
	require.NoError(t, s.Write(buf))
}
