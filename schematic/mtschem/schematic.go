// Package mtschem reads and writes the binary luanti schematic format (.mts)
//
// Format reference: src/mapgen/mg_schematic.h in https://github.com/luanti-org/luanti
//
//	All values are stored in big-endian byte order.
//	[u32] signature: 'MTSM'
//	[u16] version: 4
//	[u16] size X
//	[u16] size Y
//	[u16] size Z
//	For each Y:
//		[u8] slice probability value
//	[Name-ID table] Name ID Mapping Table
//		[u16] name-id count
//		For each name-id mapping:
//			[u16] name length
//			[u8[]] name
//	ZLib deflated {
//	For each node in schematic:  (for z, y, x)
//		[u16] content
//	For each node in schematic:
//		[u8] param1
//		  bit 0-6: probability
//		  bit 7:   specific node force placement
//	For each node in schematic:
//		[u8] param2
//	}
package mtschem

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
)

const (
	// Signature is the magic number at the start of every schematic file ('MTSM')
	Signature uint32 = 0x4d54534d
	// Version is the highest supported (and the written) format version
	Version uint16 = 4

	// ProbMask masks the probability bits of param1 (and the slice probabilities)
	ProbMask byte = 0x7F
	// ProbNever means the node is never placed
	ProbNever byte = 0x00
	// ProbAlways means the node is always placed
	ProbAlways byte = 0x7F
	// ForcePlace is the per-node force placement bit in param1
	ForcePlace byte = 0x80

	// MaxSize is the max side-length of a parsed schematic (matches the max blockexchange schema size)
	MaxSize = 500
	// MaxNameLength is the max length of a node name
	MaxNameLength = 256

	probAlwaysOld byte = 0xFF
)

// ErrInvalidSchematic is wrapped by all errors returned from Parse
var ErrInvalidSchematic = errors.New("invalid mts schematic")

// Schematic is a decoded mts schematic, normalized to the semantics of format version 4
type Schematic struct {
	SizeX int
	SizeY int
	SizeZ int
	// SliceProb contains one probability (0-127) per y-layer
	SliceProb []byte
	// NodeNames maps the content ids to node names
	NodeNames []string
	// Content, Param1 and Param2 are indexed in (z, y, x) order, see Index()
	Content []uint16
	Param1  []byte
	Param2  []byte
}

// NewSchematic creates an empty schematic of the given size, filled with "air" which is always placed
func NewSchematic(sizeX, sizeY, sizeZ int) *Schematic {
	count := sizeX * sizeY * sizeZ
	s := &Schematic{
		SizeX:     sizeX,
		SizeY:     sizeY,
		SizeZ:     sizeZ,
		SliceProb: bytes.Repeat([]byte{ProbAlways}, sizeY),
		NodeNames: []string{"air"},
		Content:   make([]uint16, count),
		Param1:    bytes.Repeat([]byte{ProbAlways}, count),
		Param2:    make([]byte, count),
	}
	return s
}

// Index returns the buffer index of the given position
func (s *Schematic) Index(x, y, z int) int {
	return x + (y * s.SizeX) + (z * s.SizeX * s.SizeY)
}

// ModNames returns the sorted list of mods referenced by the node names
func (s *Schematic) ModNames() []string {
	mods := map[string]bool{}
	for _, name := range s.NodeNames {
		modname, _, found := strings.Cut(name, ":")
		if found && modname != "" {
			mods[modname] = true
		}
	}
	res := make([]string, 0, len(mods))
	for modname := range mods {
		res = append(res, modname)
	}
	sort.Strings(res)
	return res
}

// Parse decodes a binary schematic, older format versions (1-3) are converted to version 4 semantics
func Parse(data []byte) (*Schematic, error) {
	s, err := parse(data)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidSchematic, err)
	}
	return s, nil
}

type header struct {
	Signature uint32
	Version   uint16
	SizeX     uint16
	SizeY     uint16
	SizeZ     uint16
}

func parse(data []byte) (*Schematic, error) {
	r := bytes.NewReader(data)

	h := header{}
	err := binary.Read(r, binary.BigEndian, &h)
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	if h.Signature != Signature {
		return nil, fmt.Errorf("invalid signature: 0x%08x", h.Signature)
	}
	if h.Version < 1 || h.Version > Version {
		return nil, fmt.Errorf("unsupported version: %d", h.Version)
	}
	for _, size := range []uint16{h.SizeX, h.SizeY, h.SizeZ} {
		if size < 1 || size > MaxSize {
			return nil, fmt.Errorf("schematic size %d/%d/%d out of range [1, %d]", h.SizeX, h.SizeY, h.SizeZ, MaxSize)
		}
	}

	s := &Schematic{
		SizeX:     int(h.SizeX),
		SizeY:     int(h.SizeY),
		SizeZ:     int(h.SizeZ),
		SliceProb: make([]byte, h.SizeY),
	}

	// y-slice probabilities (since version 3)
	if h.Version >= 3 {
		_, err = io.ReadFull(r, s.SliceProb)
		if err != nil {
			return nil, fmt.Errorf("read slice probabilities: %w", err)
		}
	} else {
		for i := range s.SliceProb {
			s.SliceProb[i] = probAlwaysOld
		}
	}

	// name-id table
	var nameCount uint16
	err = binary.Read(r, binary.BigEndian, &nameCount)
	if err != nil {
		return nil, fmt.Errorf("read name count: %w", err)
	}
	if nameCount == 0 {
		return nil, errors.New("empty name-id table")
	}
	ignoreID := -1
	s.NodeNames = make([]string, nameCount)
	for i := range s.NodeNames {
		var nameLen uint16
		err = binary.Read(r, binary.BigEndian, &nameLen)
		if err != nil {
			return nil, fmt.Errorf("read name length: %w", err)
		}
		if nameLen == 0 || nameLen > MaxNameLength {
			return nil, fmt.Errorf("node name length %d out of range [1, %d]", nameLen, MaxNameLength)
		}
		name := make([]byte, nameLen)
		_, err = io.ReadFull(r, name)
		if err != nil {
			return nil, fmt.Errorf("read name: %w", err)
		}
		for _, c := range name {
			if c <= ' ' || c >= 0x7F {
				return nil, fmt.Errorf("invalid character 0x%02x in node name %q", c, name)
			}
		}
		s.NodeNames[i] = string(name)
		if s.NodeNames[i] == "ignore" {
			// same as the engine: "ignore" (version 1) is converted to "air"
			s.NodeNames[i] = "air"
			ignoreID = i
		}
	}

	// node data, the read is bounded by the expected size (zip-bomb safe)
	count := s.SizeX * s.SizeY * s.SizeZ
	zr, err := zlib.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("open node data: %w", err)
	}
	defer zr.Close()
	raw := make([]byte, count*4)
	_, err = io.ReadFull(zr, raw)
	if err != nil {
		return nil, fmt.Errorf("read node data: %w", err)
	}
	// reaching the end of the stream verifies the checksum
	n, err := zr.Read(make([]byte, 1))
	if n > 0 {
		return nil, errors.New("excess node data")
	}
	if err != io.EOF {
		return nil, fmt.Errorf("read node data: %w", err)
	}

	s.Content = make([]uint16, count)
	for i := range s.Content {
		id := binary.BigEndian.Uint16(raw[i*2:])
		if int(id) >= len(s.NodeNames) {
			return nil, fmt.Errorf("content id %d at index %d not in name-id table (size %d)", id, i, len(s.NodeNames))
		}
		s.Content[i] = id
	}
	s.Param1 = raw[count*2 : count*3 : count*3]
	s.Param2 = raw[count*3:]

	if h.Version < 2 && ignoreID >= 0 {
		// former "ignore" nodes are never placed
		for i, id := range s.Content {
			if int(id) == ignoreID {
				s.Param1[i] = ProbNever
			}
		}
	}

	if h.Version < 4 {
		// probabilities used the full byte range before version 4
		for i := range s.Param1 {
			s.Param1[i] >>= 1
		}
		for i := range s.SliceProb {
			s.SliceProb[i] >>= 1
		}
	}

	return s, nil
}

// Write encodes the schematic in the binary format (version 4)
func (s *Schematic) Write(w io.Writer) error {
	count := s.SizeX * s.SizeY * s.SizeZ
	if s.SizeX < 1 || s.SizeY < 1 || s.SizeZ < 1 || s.SizeX > 0xFFFF || s.SizeY > 0xFFFF || s.SizeZ > 0xFFFF {
		return fmt.Errorf("invalid schematic size %d/%d/%d", s.SizeX, s.SizeY, s.SizeZ)
	}
	if len(s.SliceProb) != s.SizeY || len(s.Content) != count || len(s.Param1) != count || len(s.Param2) != count {
		return errors.New("schematic buffer sizes do not match the schematic size")
	}
	if len(s.NodeNames) == 0 || len(s.NodeNames) > 0xFFFF {
		return fmt.Errorf("invalid name-id table size: %d", len(s.NodeNames))
	}

	bw := bufio.NewWriter(w)
	err := binary.Write(bw, binary.BigEndian, header{
		Signature: Signature,
		Version:   Version,
		SizeX:     uint16(s.SizeX),
		SizeY:     uint16(s.SizeY),
		SizeZ:     uint16(s.SizeZ),
	})
	if err != nil {
		return err
	}
	_, err = bw.Write(s.SliceProb)
	if err != nil {
		return err
	}

	err = binary.Write(bw, binary.BigEndian, uint16(len(s.NodeNames)))
	if err != nil {
		return err
	}
	for _, name := range s.NodeNames {
		if len(name) == 0 || len(name) > 0xFFFF {
			return fmt.Errorf("invalid node name length: %d", len(name))
		}
		err = binary.Write(bw, binary.BigEndian, uint16(len(name)))
		if err != nil {
			return err
		}
		_, err = bw.WriteString(name)
		if err != nil {
			return err
		}
	}

	zw := zlib.NewWriter(bw)
	zbw := bufio.NewWriterSize(zw, 64*1024)
	buf := make([]byte, 2)
	for _, id := range s.Content {
		binary.BigEndian.PutUint16(buf, id)
		_, err = zbw.Write(buf)
		if err != nil {
			return err
		}
	}
	_, err = zbw.Write(s.Param1)
	if err != nil {
		return err
	}
	_, err = zbw.Write(s.Param2)
	if err != nil {
		return err
	}
	err = zbw.Flush()
	if err != nil {
		return err
	}
	err = zw.Close()
	if err != nil {
		return err
	}

	return bw.Flush()
}
