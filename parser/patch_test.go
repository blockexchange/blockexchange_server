package parser

import (
	"blockexchange/types"
	"fmt"
	"maps"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

// newPart creates a 2x1x1 part filled with the given node
func newPart(name string, param1, param2 byte) *ParsedSchemaPart {
	return &ParsedSchemaPart{
		NodeIDS:        []int16{1, 1},
		Param1:         []byte{param1, param1},
		Param2:         []byte{param2, param2},
		NodeNameLookup: map[int16]string{1: name},
		Meta: &SchemaPartMetadata{
			NodeMapping: map[string]int16{name: 1},
			Size:        SchemaPartSize{X: 2, Y: 1, Z: 1},
		},
	}
}

func TestGetPos(t *testing.T) {
	psp := &ParsedSchemaPart{Meta: &SchemaPartMetadata{Size: SchemaPartSize{X: 3, Y: 4, Z: 5}}}
	for x := 0; x < 3; x++ {
		for y := 0; y < 4; y++ {
			for z := 0; z < 5; z++ {
				px, py, pz := psp.GetPos(psp.GetIndex(x, y, z))
				assert.Equal(t, []int{x, y, z}, []int{px, py, pz})
			}
		}
	}
}

func TestPatch(t *testing.T) {
	base := newPart("default:stone", 1, 2)
	base.Meta.Metadata = &Metadata{
		Meta: map[string]*MetadataEntry{"(1,0,0)": {Fields: map[string]string{"a": "b"}}},
	}

	patch := newPart("default:dirt", 10, 20)
	patch.Meta.Metadata = &Metadata{
		Meta: map[string]*MetadataEntry{"(0,0,0)": {Fields: map[string]string{"x": "y"}}},
	}

	// replace everything on the first node, only param2 on the second
	err := base.Patch(patch, []types.PatchFilterType{types.PatchReplaceAll, types.PatchReplaceParam2})
	assert.NoError(t, err)

	assert.Equal(t, "default:dirt", base.NodeNameLookup[base.NodeIDS[0]])
	assert.Equal(t, "default:stone", base.NodeNameLookup[base.NodeIDS[1]])
	assert.Equal(t, []byte{10, 1}, base.Param1)
	assert.Equal(t, []byte{20, 20}, base.Param2)
	assert.Equal(t, "y", base.Meta.Metadata.Meta["(0,0,0)"].Fields["x"])
	assert.Equal(t, "b", base.Meta.Metadata.Meta["(1,0,0)"].Fields["a"])
	assert.Len(t, base.Meta.NodeMapping, 2)

	// replace the remaining stone node, mapping and metadata get removed
	err = base.Patch(patch, []types.PatchFilterType{types.PatchNoOp, types.PatchReplaceNodename | types.PatchReplaceMetadata})
	assert.NoError(t, err)
	assert.Equal(t, map[string]int16{"default:dirt": base.NodeIDS[0]}, base.Meta.NodeMapping)
	assert.Equal(t, base.NodeIDS[0], base.NodeIDS[1])
	assert.Nil(t, base.Meta.Metadata.Meta["(1,0,0)"])
	assert.Equal(t, []byte{10, 1}, base.Param1)

	// roundtrip
	sp, err := base.Convert()
	assert.NoError(t, err)
	parsed, err := ParseSchemaPart(sp)
	assert.NoError(t, err)
	assert.Equal(t, base.NodeIDS, parsed.NodeIDS)
	assert.Equal(t, base.Meta.NodeMapping, parsed.Meta.NodeMapping)
}

func TestPatchMismatch(t *testing.T) {
	base := newPart("default:stone", 0, 0)
	patch := newPart("default:dirt", 0, 0)

	assert.Error(t, base.Patch(patch, []types.PatchFilterType{types.PatchReplaceAll}))

	patch.Meta.Size.X = 1
	assert.Error(t, base.Patch(patch, []types.PatchFilterType{types.PatchReplaceAll, types.PatchReplaceAll}))
}

func TestIsAirOnly(t *testing.T) {
	assert.True(t, newPart("air", 0, 0).IsAirOnly())
	assert.False(t, newPart("default:stone", 0, 0).IsAirOnly())

	psp := newPart("air", 0, 0)
	psp.Meta.NodeMapping["air"] = 0
	assert.True(t, psp.IsAirOnly())
	psp.Meta.NodeMapping["default:stone"] = 1
	assert.False(t, psp.IsAirOnly())
}

// referencePatchMetadata is the straightforward per-node metadata patch, used to verify patchMetadata
func referencePatchMetadata(psp, patch *ParsedSchemaPart, filter []types.PatchFilterType) {
	md := psp.GetMetadata()
	for i, f := range filter {
		if f&types.PatchReplaceMetadata == 0 {
			continue
		}
		key := md.GetKey(psp.GetPos(i))
		delete(md.Meta, key)
		delete(md.Timers, key)
		if patch.Meta.Metadata != nil {
			if e := patch.Meta.Metadata.Meta[key]; e != nil {
				md.Meta[key] = e
			}
			if t := patch.Meta.Metadata.Timers[key]; t != nil {
				md.Timers[key] = t
			}
		}
	}
}

func TestPatchMetadataMatchesReference(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	size := SchemaPartSize{X: 4, Y: 3, Z: 5}
	volume := size.X * size.Y * size.Z

	randomPart := func() *ParsedSchemaPart {
		md := &Metadata{Meta: map[string]*MetadataEntry{}, Timers: map[string]*MetadataTimer{}}
		psp := &ParsedSchemaPart{Meta: &SchemaPartMetadata{Size: size, Metadata: md}}
		for i := 0; i < volume; i++ {
			key := md.GetKey(psp.GetPos(i))
			if rnd.Intn(3) == 0 {
				md.Meta[key] = &MetadataEntry{Fields: map[string]string{"n": fmt.Sprint(rnd.Int())}}
			}
			if rnd.Intn(4) == 0 {
				md.Timers[key] = &MetadataTimer{Timeout: rnd.Float64()}
			}
		}
		// out of bounds and malformed keys are left untouched
		md.Meta["(99,0,0)"] = &MetadataEntry{}
		md.Meta["garbage"] = &MetadataEntry{}
		return psp
	}
	clone := func(psp *ParsedSchemaPart) *ParsedSchemaPart {
		md := &Metadata{Meta: maps.Clone(psp.Meta.Metadata.Meta), Timers: maps.Clone(psp.Meta.Metadata.Timers)}
		return &ParsedSchemaPart{Meta: &SchemaPartMetadata{Size: size, Metadata: md}}
	}

	for run := 0; run < 50; run++ {
		base, patch := randomPart(), randomPart()
		filter := make([]types.PatchFilterType, volume)
		for i := range filter {
			if rnd.Intn(2) == 0 {
				filter[i] = types.PatchReplaceMetadata
			}
		}
		expected := clone(base)
		referencePatchMetadata(expected, patch, filter)
		base.patchMetadata(patch, filter)
		assert.Equal(t, expected.Meta.Metadata, base.Meta.Metadata)
	}

	// patch without metadata clears the selected entries
	base := randomPart()
	expected := clone(base)
	empty := &ParsedSchemaPart{Meta: &SchemaPartMetadata{Size: size}}
	all := make([]types.PatchFilterType, volume)
	for i := range all {
		all[i] = types.PatchReplaceAll
	}
	referencePatchMetadata(expected, empty, all)
	base.patchMetadata(empty, all)
	assert.Equal(t, expected.Meta.Metadata, base.Meta.Metadata)
}
