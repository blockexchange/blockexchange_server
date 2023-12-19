package worldedit

import (
	"fmt"
	"strings"

	"github.com/Shopify/go-lua"
)

func Parse(data []byte) ([]*WEEntry, []string, error) {
	if len(data) == 0 {
		return nil, nil, fmt.Errorf("data length is zero")
	}
	if data[0] != byte('5') || data[1] != byte(':') {
		return nil, nil, fmt.Errorf("invalid format: %d", data[0])
	}

	L := lua.NewState()
	err := lua.DoString(L, string(data[2:]))
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't execute we file: %v", err)
	}

	if !L.IsTable(L.Top()) {
		return nil, nil, fmt.Errorf("root table not found")
	}

	L.Length(L.Top())
	if !L.IsNumber(L.Top()) {
		return nil, nil, fmt.Errorf("get length error")
	}
	entry_count, ok := L.ToNumber(L.Top())
	if !ok {
		return nil, nil, fmt.Errorf("couldn't get length number")
	}
	L.Pop(1)

	entries := make([]*WEEntry, int(entry_count))
	modname_map := make(map[string]bool)

	for i := 1; i <= int(entry_count); i++ {
		L.PushNumber(float64(i))
		L.Table(L.Top() - 1)

		x, err := get_table_integer(L, "x")
		if err != nil {
			return nil, nil, fmt.Errorf("error in x-entry %d: %v", i, err)
		}
		if x == nil {
			return nil, nil, fmt.Errorf("nil entry for x in %d", i)
		}

		y, err := get_table_integer(L, "y")
		if err != nil {
			return nil, nil, fmt.Errorf("error in y-entry %d: %v", i, err)
		}
		if y == nil {
			return nil, nil, fmt.Errorf("nil entry for y in %d", i)
		}

		z, err := get_table_integer(L, "z")
		if err != nil {
			return nil, nil, fmt.Errorf("error in z-entry %d: %v", i, err)
		}
		if z == nil {
			return nil, nil, fmt.Errorf("nil entry for z in %d", i)
		}

		name, err := get_table_string(L, "name")
		if err != nil {
			return nil, nil, fmt.Errorf("error in name-entry %d: %v", i, err)
		}
		if name == nil {
			return nil, nil, fmt.Errorf("nil entry for name in %d", i)
		}

		// set modname flag
		modname_map[strings.Split(*name, ":")[0]] = true

		entry := &WEEntry{
			Name: *name,
			X:    *x,
			Y:    *y,
			Z:    *z,
		}
		entries[i-1] = entry

		// optional entries
		param1, err := get_table_integer(L, "param1")
		if err != nil {
			return nil, nil, fmt.Errorf("error in param1-entry %d: %v", i, err)
		}
		if param1 != nil {
			entry.Param1 = byte(*param1)
		}

		param2, err := get_table_integer(L, "param2")
		if err != nil {
			return nil, nil, fmt.Errorf("error in param2-entry %d: %v", i, err)
		}
		if param2 != nil {
			entry.Param2 = byte(*param2)
		}

		L.PushString("meta")
		L.Table(L.Top() - 1)
		if L.IsTable(L.Top()) {
			entry.Meta = &WEMeta{
				Fields:    WEFields{},
				Inventory: WEInventory{},
			}

			L.PushString("fields")
			L.Table(L.Top() - 1) // entry.meta.fields

			L.PushNil() // value
			for L.Next(-2) {
				key := lua.CheckString(L, -2)
				value := lua.CheckString(L, -1)
				entry.Meta.Fields[key] = value
				L.Pop(1) // value
			}
			L.Pop(1) // entry.meta.fields

			L.PushString("inventory")
			L.Table(L.Top() - 1) // entry.meta.inventory
			L.PushNil()          // value
			for L.Next(-2) {
				key := lua.CheckString(L, -2)
				inv := []string{}

				L.PushNil() // {}
				for L.Next(-2) {
					entry := lua.CheckString(L, -1)
					inv = append(inv, entry)
					L.Pop(1) // {}
				}
				entry.Meta.Inventory[key] = inv

				L.Pop(1) // value
			}

			L.Pop(1) // entry.meta.inventory
		}

		L.Pop(1) // entry.meta
		L.Pop(1) // entry
	}

	modnames := make([]string, 0)
	for m := range modname_map {
		modnames = append(modnames, m)
	}

	return entries, modnames, nil
}

func get_table_integer(L *lua.State, key string) (*int, error) {
	if L.IsNil(L.Top()) {
		return nil, fmt.Errorf("table is nil")
	}

	L.PushString(key)
	L.Table(L.Top() - 1)
	if L.IsNil(L.Top()) {
		L.Pop(1)
		return nil, nil
	}

	if !L.IsNumber(L.Top()) {
		return nil, fmt.Errorf("not a number")
	}
	n, ok := L.ToNumber(L.Top())
	if !ok {
		return nil, fmt.Errorf("couldn't get number")
	}
	L.Pop(1)

	i := int(n)

	return &i, nil
}

func get_table_string(L *lua.State, key string) (*string, error) {
	if L.IsNil(L.Top()) {
		return nil, fmt.Errorf("table is nil")
	}

	L.PushString(key)
	L.Table(L.Top() - 1)
	if L.IsNil(L.Top()) {
		L.Pop(1)
		return nil, nil
	}

	if !L.IsString(L.Top()) {
		return nil, fmt.Errorf("not a string")
	}
	s, ok := L.ToString(L.Top())
	if !ok {
		return nil, fmt.Errorf("couldn't get string")
	}
	L.Pop(1)

	return &s, nil
}
