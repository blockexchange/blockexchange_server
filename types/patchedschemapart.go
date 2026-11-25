package types

import (
	"encoding/base64"
	"encoding/json"
)

// A patch schema_part request
type PatchedSchemapart struct {
	*SchemaPart

	// a list of values corresponding to the map-data, they define what happens to that node
	PatchFilter []PatchFilterType `json:"patchfilter"`
}

type PatchFilterType byte

const (
	PatchNoOp            PatchFilterType = 0x00 // nothing happens to this node
	PatchReplaceNodename PatchFilterType = 0x01 // only node name is replaced
	PatchReplaceParam1   PatchFilterType = 0x02 // only param1 is replaced
	PatchReplaceParam2   PatchFilterType = 0x04 // only param2 is replaced
	PatchReplaceMetadata PatchFilterType = 0x08 // only metadata is replaced

	// whole node is replaced with incoming value
	PatchReplaceAll PatchFilterType = PatchReplaceNodename | PatchReplaceParam1 | PatchReplaceParam2 | PatchReplaceMetadata
)

// UnmarshalJSON decodes the embedded SchemaPart fields and the patchfilter, which is
// transferred as unpadded base64 like Data and MetaData. Without it the promoted
// SchemaPart.UnmarshalJSON would drop the patchfilter.
func (p *PatchedSchemapart) UnmarshalJSON(data []byte) error {
	m := make(map[string]any)
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	if p.SchemaPart == nil {
		p.SchemaPart = &SchemaPart{}
	}
	err = p.SchemaPart.fromMap(m)
	if err != nil {
		return err
	}

	pf, err := base64.RawStdEncoding.DecodeString(getString(m["patchfilter"]))
	if err != nil {
		return err
	}
	p.PatchFilter = make([]PatchFilterType, len(pf))
	for i, b := range pf {
		p.PatchFilter[i] = PatchFilterType(b)
	}
	return nil
}

// MarshalJSON encodes the embedded SchemaPart fields and adds the patchfilter as
// unpadded base64.
func (p PatchedSchemapart) MarshalJSON() ([]byte, error) {
	sp := p.SchemaPart
	if sp == nil {
		sp = &SchemaPart{}
	}
	m := sp.toMap()

	pf := make([]byte, len(p.PatchFilter))
	for i, f := range p.PatchFilter {
		pf[i] = byte(f)
	}
	m["patchfilter"] = base64.RawStdEncoding.EncodeToString(pf)

	return json.Marshal(m)
}
