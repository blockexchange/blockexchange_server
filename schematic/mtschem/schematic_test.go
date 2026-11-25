package mtschem_test

import (
	"blockexchange/schematic/mtschem"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testfile struct {
	filename            string
	sizeX, sizeY, sizeZ int
	names               int
	mods                []string
}

var testfiles = []testfile{
	{"testdata/Notre_Dame(Paris,France).mts", 60, 64, 143, 18, []string{"mcl_core", "mcl_doors", "mcl_fences", "mcl_stairs", "mcl_torches", "mcl_wool"}},
	{"testdata/Statue_of_Liberty(New_York,USA).mts", 44, 130, 42, 11, nil},
	{"testdata/St.Peters_Basilica(Vatican,Italy).mts", 73, 60, 138, 15, nil},
}

func parseFile(t *testing.T, filename string) *mtschem.Schematic {
	data, err := os.ReadFile(filename)
	require.NoError(t, err)
	s, err := mtschem.Parse(data)
	require.NoError(t, err)
	return s
}

func TestParse(t *testing.T) {
	for _, tf := range testfiles {
		t.Run(tf.filename, func(t *testing.T) {
			s := parseFile(t, tf.filename)
			assert.Equal(t, tf.sizeX, s.SizeX)
			assert.Equal(t, tf.sizeY, s.SizeY)
			assert.Equal(t, tf.sizeZ, s.SizeZ)
			assert.Equal(t, tf.names, len(s.NodeNames))
			assert.Equal(t, "air", s.NodeNames[0])

			count := tf.sizeX * tf.sizeY * tf.sizeZ
			assert.Equal(t, count, len(s.Content))
			assert.Equal(t, count, len(s.Param1))
			assert.Equal(t, count, len(s.Param2))
			assert.Equal(t, tf.sizeY, len(s.SliceProb))
			for _, p := range s.SliceProb {
				assert.Equal(t, mtschem.ProbAlways, p)
			}
			if tf.mods != nil {
				assert.Equal(t, tf.mods, s.ModNames())
			}
		})
	}
}

func TestParseModNames(t *testing.T) {
	s := parseFile(t, "testdata/Statue_of_Liberty(New_York,USA).mts")
	mods := s.ModNames()
	assert.Contains(t, mods, "default")
	assert.NotContains(t, mods, "air")
}

func TestWriteRoundtrip(t *testing.T) {
	for _, tf := range testfiles {
		t.Run(tf.filename, func(t *testing.T) {
			s := parseFile(t, tf.filename)

			buf := bytes.NewBuffer([]byte{})
			require.NoError(t, s.Write(buf))

			s2, err := mtschem.Parse(buf.Bytes())
			require.NoError(t, err)
			assert.Equal(t, s, s2)
		})
	}
}

// builds a raw schematic file with the given version and node data
func buildRaw(t *testing.T, version uint16, sx, sy, sz uint16, slices []byte, names []string, content []uint16, param1, param2 []byte) []byte {
	buf := bytes.NewBuffer([]byte{})
	w := func(v any) { require.NoError(t, binary.Write(buf, binary.BigEndian, v)) }
	w(mtschem.Signature)
	w(version)
	w(sx)
	w(sy)
	w(sz)
	buf.Write(slices)
	w(uint16(len(names)))
	for _, n := range names {
		w(uint16(len(n)))
		buf.WriteString(n)
	}
	zw := zlib.NewWriter(buf)
	for _, c := range content {
		require.NoError(t, binary.Write(zw, binary.BigEndian, c))
	}
	zw.Write(param1)
	zw.Write(param2)
	require.NoError(t, zw.Close())
	return buf.Bytes()
}

func TestParseOldVersions(t *testing.T) {
	// version 1: no slice probabilities, "ignore" nodes, full-byte probabilities
	data := buildRaw(t, 1, 2, 1, 1, nil,
		[]string{"ignore", "default:stone"},
		[]uint16{0, 1}, []byte{0xFF, 0xFF}, []byte{0, 3})
	s, err := mtschem.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, []string{"air", "default:stone"}, s.NodeNames)
	assert.Equal(t, []byte{mtschem.ProbAlways}, s.SliceProb)
	assert.Equal(t, []byte{mtschem.ProbNever, mtschem.ProbAlways}, s.Param1)
	assert.Equal(t, []byte{0, 3}, s.Param2)

	// version 3: slice probabilities, full-byte probabilities
	data = buildRaw(t, 3, 1, 2, 1, []byte{0xFF, 0x80},
		[]string{"air", "default:stone"},
		[]uint16{1, 1}, []byte{0xFF, 0x80}, []byte{0, 0})
	s, err = mtschem.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, []byte{0x7F, 0x40}, s.SliceProb)
	assert.Equal(t, []byte{0x7F, 0x40}, s.Param1)
}

func TestParseInvalid(t *testing.T) {
	valid := buildRaw(t, 4, 1, 1, 1, []byte{0x7F}, []string{"air"}, []uint16{0}, []byte{0x7F}, []byte{0})
	_, err := mtschem.Parse(valid)
	require.NoError(t, err)

	cases := map[string][]byte{
		"empty":            {},
		"bad signature":    append([]byte("XXXX"), valid[4:]...),
		"truncated header": valid[:8],
		"truncated data":   valid[:len(valid)-4],
		"version 5":        buildRaw(t, 5, 1, 1, 1, []byte{0x7F}, []string{"air"}, []uint16{0}, []byte{0x7F}, []byte{0}),
		"zero size":        buildRaw(t, 4, 0, 1, 1, []byte{0x7F}, []string{"air"}, nil, nil, nil),
		"oversized":        buildRaw(t, 4, 501, 1, 1, []byte{0x7F}, []string{"air"}, nil, nil, nil),
		"no names":         buildRaw(t, 4, 1, 1, 1, []byte{0x7F}, nil, []uint16{0}, []byte{0x7F}, []byte{0}),
		"invalid name":     buildRaw(t, 4, 1, 1, 1, []byte{0x7F}, []string{"a\nb"}, []uint16{0}, []byte{0x7F}, []byte{0}),
		"content id oob":   buildRaw(t, 4, 1, 1, 1, []byte{0x7F}, []string{"air"}, []uint16{1}, []byte{0x7F}, []byte{0}),
	}

	for name, data := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := mtschem.Parse(data)
			assert.ErrorIs(t, err, mtschem.ErrInvalidSchematic)
		})
	}
}
