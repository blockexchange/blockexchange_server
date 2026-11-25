package parser

import (
	"blockexchange/types"
	"fmt"
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
			psp.patchMetadata(patch, i)
		}
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

// patchMetadata replaces the metadata and timer of the node at the given index with the one of the patch-part
func (psp *ParsedSchemaPart) patchMetadata(patch *ParsedSchemaPart, index int) {
	md := psp.GetMetadata()
	key := md.GetKey(psp.GetPos(index))

	var entry *MetadataEntry
	var timer *MetadataTimer
	if patch.Meta.Metadata != nil {
		entry = patch.Meta.Metadata.Meta[key]
		timer = patch.Meta.Metadata.Timers[key]
	}

	if entry != nil {
		md.Meta[key] = entry
	} else {
		delete(md.Meta, key)
	}
	if timer != nil {
		md.Timers[key] = timer
	} else {
		delete(md.Timers, key)
	}
}
