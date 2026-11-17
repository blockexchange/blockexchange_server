package worldedit

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	// MaxCoord is the largest allowed node coordinate (inclusive),
	// matches the max schema side-length of 500 nodes
	MaxCoord = 499
	// MaxEntries is the max number of nodes in a schematic
	MaxEntries = 5_000_000
	// MaxNameLength is the max length of a node name
	MaxNameLength = 256
	// MaxMetaFields is the max number of metadata fields per node
	MaxMetaFields = 1024
	// MaxFieldKeyLength is the max length of a metadata field or inventory name
	MaxFieldKeyLength = 256
	// MaxInventoryLists is the max number of inventory lists per node
	MaxInventoryLists = 64
	// MaxInventorySize is the max number of stacks per inventory list
	MaxInventorySize = 4096
)

// ErrInvalidSchematic is wrapped by all errors returned from Parse
var ErrInvalidSchematic = errors.New("invalid worldedit schematic")

// nodeNameRegex matches "modname:itemname" as allowed by the engine
var nodeNameRegex = regexp.MustCompile(`^[a-z0-9_]+:[a-zA-Z0-9_]+$`)

var formatHeader = []byte("5:")

// Parse parses a worldedit schematic (format version 5) without executing it.
// It returns the node entries and the list of referenced mod names.
// Entries for "air" are skipped.
func Parse(data []byte) ([]*WEEntry, []string, error) {
	entries, modnames, err := parse(data)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", ErrInvalidSchematic, err)
	}
	return entries, modnames, nil
}

func parse(data []byte) ([]*WEEntry, []string, error) {
	if !bytes.HasPrefix(data, formatHeader) {
		return nil, nil, fmt.Errorf("unsupported format, expected header %q, got %q", formatHeader, data[:min(len(data), 8)])
	}

	p := newLuaParser(data, len(formatHeader))
	if err := p.expectKeyword("return"); err != nil {
		return nil, nil, err
	}

	entries := make([]*WEEntry, 0)
	modnameMap := make(map[string]bool)
	index := 0

	err := p.parseTableItems(func(key, value any) error {
		index++
		if key != nil {
			return p.errorf("entry %d: root table must be a list", index)
		}
		if index > MaxEntries {
			return p.errorf("more than %d entries", MaxEntries)
		}
		t, ok := value.(*luaTable)
		if !ok {
			return p.errorf("entry %d: expected a table, got %s", index, typeName(value))
		}
		entry, err := toEntry(t)
		if err != nil {
			return p.errorf("entry %d: %v", index, err)
		}
		if entry == nil {
			// air
			return nil
		}

		entries = append(entries, entry)
		modnameMap[strings.SplitN(entry.Name, ":", 2)[0]] = true
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	if err = p.expectEOF(); err != nil {
		return nil, nil, err
	}
	if len(entries) == 0 {
		return nil, nil, errors.New("schematic contains no nodes")
	}

	modnames := make([]string, 0, len(modnameMap))
	for m := range modnameMap {
		modnames = append(modnames, m)
	}
	sort.Strings(modnames)

	return entries, modnames, nil
}

// toEntry converts and validates a single node table, returns nil for air
func toEntry(t *luaTable) (*WEEntry, error) {
	name, err := getString(t, "name", true)
	if err != nil {
		return nil, err
	}
	if name == "air" {
		return nil, nil
	}
	if len(name) > MaxNameLength {
		return nil, fmt.Errorf("name longer than %d bytes", MaxNameLength)
	}
	if !nodeNameRegex.MatchString(name) {
		return nil, fmt.Errorf("invalid node name %q", name)
	}

	entry := &WEEntry{Name: name}

	if entry.X, err = getInt(t, "x", true, 0, MaxCoord); err != nil {
		return nil, err
	}
	if entry.Y, err = getInt(t, "y", true, 0, MaxCoord); err != nil {
		return nil, err
	}
	if entry.Z, err = getInt(t, "z", true, 0, MaxCoord); err != nil {
		return nil, err
	}

	param1, err := getInt(t, "param1", false, 0, 255)
	if err != nil {
		return nil, err
	}
	entry.Param1 = byte(param1)

	param2, err := getInt(t, "param2", false, 0, 255)
	if err != nil {
		return nil, err
	}
	entry.Param2 = byte(param2)

	entry.Meta, err = toMeta(t.get("meta"))
	if err != nil {
		return nil, fmt.Errorf("meta: %v", err)
	}

	return entry, nil
}

// toMeta converts the optional "meta" table, returns nil if empty
func toMeta(v any) (*WEMeta, error) {
	if v == nil {
		return nil, nil
	}
	t, ok := v.(*luaTable)
	if !ok {
		return nil, fmt.Errorf("expected a table, got %s", typeName(v))
	}

	fields, err := toFields(t.get("fields"))
	if err != nil {
		return nil, fmt.Errorf("fields: %v", err)
	}
	inventory, err := toInventory(t.get("inventory"))
	if err != nil {
		return nil, fmt.Errorf("inventory: %v", err)
	}

	if len(fields) == 0 && len(inventory) == 0 {
		return nil, nil
	}
	return &WEMeta{Fields: fields, Inventory: inventory}, nil
}

func toFields(v any) (WEFields, error) {
	fields := WEFields{}
	if v == nil {
		return fields, nil
	}
	t, ok := v.(*luaTable)
	if !ok {
		return nil, fmt.Errorf("expected a table, got %s", typeName(v))
	}
	if len(t.arr) > 0 || len(t.idx) > 0 {
		return nil, errors.New("keys must be strings")
	}
	if len(t.hash) > MaxMetaFields {
		return nil, fmt.Errorf("more than %d fields", MaxMetaFields)
	}

	for key, value := range t.hash {
		if len(key) > MaxFieldKeyLength {
			return nil, fmt.Errorf("key longer than %d bytes", MaxFieldKeyLength)
		}
		switch val := value.(type) {
		case string:
			fields[key] = val
		case float64:
			// numbers are stored as strings by the engine
			fields[key] = formatNumber(val)
		default:
			return nil, fmt.Errorf("field %q: expected a string, got %s", key, typeName(value))
		}
	}
	return fields, nil
}

func toInventory(v any) (WEInventory, error) {
	inventory := WEInventory{}
	if v == nil {
		return inventory, nil
	}
	t, ok := v.(*luaTable)
	if !ok {
		return nil, fmt.Errorf("expected a table, got %s", typeName(v))
	}
	if len(t.arr) > 0 || len(t.idx) > 0 {
		return nil, errors.New("list names must be strings")
	}
	if len(t.hash) > MaxInventoryLists {
		return nil, fmt.Errorf("more than %d lists", MaxInventoryLists)
	}

	for listname, value := range t.hash {
		if len(listname) > MaxFieldKeyLength {
			return nil, fmt.Errorf("list name longer than %d bytes", MaxFieldKeyLength)
		}
		lt, ok := value.(*luaTable)
		if !ok {
			return nil, fmt.Errorf("list %q: expected a table, got %s", listname, typeName(value))
		}
		items, err := lt.list()
		if err != nil {
			return nil, fmt.Errorf("list %q: %v", listname, err)
		}
		if len(items) > MaxInventorySize {
			return nil, fmt.Errorf("list %q: more than %d stacks", listname, MaxInventorySize)
		}

		stacks := make([]string, len(items))
		for i, item := range items {
			s, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("list %q, stack %d: expected a string, got %s", listname, i+1, typeName(item))
			}
			stacks[i] = s
		}
		inventory[listname] = stacks
	}
	return inventory, nil
}

func getString(t *luaTable, key string, required bool) (string, error) {
	v := t.get(key)
	if v == nil {
		if required {
			return "", fmt.Errorf("missing required field %q", key)
		}
		return "", nil
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("field %q: expected a string, got %s", key, typeName(v))
	}
	return s, nil
}

// getInt returns an integer field within [minValue, maxValue], 0 if optional and missing
func getInt(t *luaTable, key string, required bool, minValue, maxValue int) (int, error) {
	v := t.get(key)
	if v == nil {
		if required {
			return 0, fmt.Errorf("missing required field %q", key)
		}
		return 0, nil
	}
	f, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("field %q: expected a number, got %s", key, typeName(v))
	}
	if math.IsNaN(f) || math.IsInf(f, 0) || f != math.Trunc(f) {
		return 0, fmt.Errorf("field %q: expected an integer, got %v", key, f)
	}
	if f < float64(minValue) || f > float64(maxValue) {
		return 0, fmt.Errorf("field %q: value %v out of range [%d, %d]", key, f, minValue, maxValue)
	}
	return int(f), nil
}

func formatNumber(f float64) string {
	if f == math.Trunc(f) && math.Abs(f) < 1e15 {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'g', 14, 64)
}

func typeName(v any) string {
	switch v.(type) {
	case nil:
		return "nil"
	case bool:
		return "boolean"
	case float64:
		return "number"
	case string:
		return "string"
	case *luaTable:
		return "table"
	default:
		return fmt.Sprintf("%T", v)
	}
}
