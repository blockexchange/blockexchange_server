package parser

import (
	"blockexchange/types"
	"fmt"
	"strconv"
	"strings"
)

// GetPos returns the relative position of the node at the given index (inverse of GetIndex)
func (psp *ParsedSchemaPart) GetPos(index int) (int, int, int) {
	sy, sz := psp.Meta.Size.Y, psp.Meta.Size.Z
	return index / (sy * sz), (index / sz) % sy, index % sz
}

// GetOrCreateNodeID returns the node-id for the given name, a new mapping is added if it doesn't exist yet
func (psp *ParsedSchemaPart) GetOrCreateNodeID(name string) int16 {
	if id, found := psp.Meta.NodeMapping[name]; found {
		return id
	}

	maxnodeid := int16(0)
	for id := range psp.NodeNameLookup {
		maxnodeid = max(id, maxnodeid)
	}
	id := maxnodeid + 1
	psp.NodeNameLookup[id] = name
	psp.Meta.NodeMapping[name] = id
	return id
}

// removeUnusedNodeIDs drops node-mappings that aren't referenced by any node anymore
func (psp *ParsedSchemaPart) removeUnusedNodeIDs() {
	used := map[int16]bool{}
	for _, id := range psp.NodeIDS {
		used[id] = true
	}
	for id, name := range psp.NodeNameLookup {
		if !used[id] {
			delete(psp.NodeNameLookup, id)
			delete(psp.Meta.NodeMapping, name)
		}
	}
}

// Patch applies the nodes of the given patch-part to this part, the filter defines per node what gets replaced
func (psp *ParsedSchemaPart) Patch(patch *ParsedSchemaPart, filter []types.PatchFilterType) error {
	if psp.Meta.Size != patch.Meta.Size {
		return fmt.Errorf("size mismatch: %v vs %v", psp.Meta.Size, patch.Meta.Size)
	}
	if len(filter) != len(psp.NodeIDS) {
		return fmt.Errorf("patchfilter length mismatch: got %d, expected %d", len(filter), len(psp.NodeIDS))
	}

	patch_metadata := false
	for i, f := range filter {
		if f&types.PatchReplaceNodename != 0 {
			name, found := patch.NodeNameLookup[patch.NodeIDS[i]]
			if !found {
				return fmt.Errorf("unmapped node-id %d at index %d", patch.NodeIDS[i], i)
			}
			psp.NodeIDS[i] = psp.GetOrCreateNodeID(name)
		}
		if f&types.PatchReplaceParam1 != 0 {
			psp.Param1[i] = patch.Param1[i]
		}
		if f&types.PatchReplaceParam2 != 0 {
			psp.Param2[i] = patch.Param2[i]
		}
		if f&types.PatchReplaceMetadata != 0 {
			patch_metadata = true
		}
	}

	if patch_metadata {
		psp.patchMetadata(patch, filter)
	}

	psp.removeUnusedNodeIDs()
	return nil
}

// GetMetadata returns the node-metadata of the part, missing fields are initialized
func (psp *ParsedSchemaPart) GetMetadata() *Metadata {
	if psp.Meta.Metadata == nil {
		psp.Meta.Metadata = &Metadata{}
	}
	md := psp.Meta.Metadata
	if md.Meta == nil {
		md.Meta = map[string]*MetadataEntry{}
	}
	if md.Timers == nil {
		md.Timers = map[string]*MetadataTimer{}
	}
	return md
}

// parseKey is the inverse of Metadata.GetKey, it returns the node-index of a "(x,y,z)" key
func (psp *ParsedSchemaPart) parseKey(key string) (int, bool) {
	inner, found := strings.CutPrefix(key, "(")
	if !found {
		return 0, false
	}
	inner, found = strings.CutSuffix(inner, ")")
	if !found {
		return 0, false
	}
	parts := strings.Split(inner, ",")
	if len(parts) != 3 {
		return 0, false
	}
	var pos [3]int
	for i, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil {
			return 0, false
		}
		pos[i] = v
	}
	size := psp.Meta.Size
	if pos[0] < 0 || pos[1] < 0 || pos[2] < 0 || pos[0] >= size.X || pos[1] >= size.Y || pos[2] >= size.Z {
		return 0, false
	}
	return psp.GetIndex(pos[0], pos[1], pos[2]), true
}

// patchMetadata replaces the metadata and timers of all nodes selected for metadata replacement
// with the ones of the patch-part. Only the existing entries are visited, instead of building
// a key for every selected node.
func (psp *ParsedSchemaPart) patchMetadata(patch *ParsedSchemaPart, filter []types.PatchFilterType) {
	md := psp.GetMetadata()
	selected := func(key string) bool {
		index, ok := psp.parseKey(key)
		return ok && filter[index]&types.PatchReplaceMetadata != 0
	}

	// selected nodes get the patch values, so drop their current entries
	for key := range md.Meta {
		if selected(key) {
			delete(md.Meta, key)
		}
	}
	for key := range md.Timers {
		if selected(key) {
			delete(md.Timers, key)
		}
	}

	if patch.Meta.Metadata == nil {
		return
	}
	for key, entry := range patch.Meta.Metadata.Meta {
		if entry != nil && selected(key) {
			md.Meta[key] = entry
		}
	}
	for key, timer := range patch.Meta.Metadata.Timers {
		if timer != nil && selected(key) {
			md.Timers[key] = timer
		}
	}
}
