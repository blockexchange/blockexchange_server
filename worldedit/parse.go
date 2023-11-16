package worldedit

import (
	"fmt"

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
			L.PushString("fields")
			L.Table(L.Top() - 1)

			L.Pop(1) // entry.meta.fields
		}

		L.Pop(1) // entry.meta
		L.Pop(1) // entry
	}

	modnames := make([]string, 0)
	for m := range modname_map {
		modnames = append(modnames, m)
	}

	return entries, modnames, nil

	/*
		lv := L.Get(-1)
		tbl, ok := lv.(*lua.LTable)
		if !ok {
			return nil, nil, errors.New("root table not found")
		}

		entry_count := L.ObjLen(tbl)
		entries := make([]*WEEntry, entry_count)
		modname_map := make(map[string]bool)

		for i := 1; i <= entry_count; i++ {
			entrytbl := L.GetTable(tbl, lua.LNumber(i))

			// common entries
			entry, err := parseWEEntry(L, entrytbl)
			if err != nil {
				return nil, nil, err
			}
			entries[i-1] = entry

			// set modname flag
			modname_map[strings.Split(entry.Name, ":")[0]] = true

			metatbl := L.GetTable(entrytbl, lua.LString("meta"))
			if metatbl != lua.LNil {
				entry.Meta = &WEMeta{
					Fields:    make(WEFields),
					Inventory: make(WEInventory),
				}

				// fields
				fieldvalue := L.GetTable(metatbl, lua.LString("fields"))
				if fieldvalue != lua.LNil {
					fieldstbl := fieldvalue.(*lua.LTable)
					fieldstbl.ForEach(func(key, value lua.LValue) {
						entry.Meta.Fields[key.String()] = value.String()
					})
				}

				invvalue := L.GetTable(metatbl, lua.LString("inventory"))
				if invvalue != lua.LNil {
					invtbl := invvalue.(*lua.LTable)
					invtbl.ForEach(func(key, value lua.LValue) {
						stacks := make([]string, 0)
						stacktbl := value.(*lua.LTable)
						stacktbl.ForEach(func(_, stack lua.LValue) {
							stacks = append(stacks, stack.String())
						})
						entry.Meta.Inventory[key.String()] = stacks
					})
				}
			}
		}

		modnames := make([]string, 0)
		for m := range modname_map {
			modnames = append(modnames, m)
		}

		return entries, modnames, nil
	*/
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
