package worldedit

import (
	"errors"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func Parse(data []byte) ([]*WEEntry, []string, error) {
	if data[0] != byte('5') || data[1] != byte(':') {
		return nil, nil, errors.New("invalid format")
	}

	L := lua.NewState()
	defer L.Close()

	fn, err := L.LoadString(string(data[2:]))
	if err != nil {
		return nil, nil, err
	}

	err = L.CallByParam(lua.P{
		Fn:      fn,
		NRet:    1,
		Protect: true,
	})
	if err != nil {
		return nil, nil, err
	}

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
}

func parseWEEntry(L *lua.LState, tbl lua.LValue) (*WEEntry, error) {
	e := &WEEntry{}

	name, ok := L.GetTable(tbl, lua.LString("name")).(lua.LString)
	if ok {
		e.Name = name.String()
	}

	param2, ok := L.GetTable(tbl, lua.LString("param2")).(lua.LNumber)
	if ok {
		e.Param2 = byte(param2)
	}

	param1, ok := L.GetTable(tbl, lua.LString("param1")).(lua.LNumber)
	if ok {
		e.Param1 = byte(param1)
	}

	x, ok := L.GetTable(tbl, lua.LString("x")).(lua.LNumber)
	if ok {
		e.X = int(x)
	}

	y, ok := L.GetTable(tbl, lua.LString("y")).(lua.LNumber)
	if ok {
		e.Y = int(y)
	}

	z, ok := L.GetTable(tbl, lua.LString("z")).(lua.LNumber)
	if ok {
		e.Z = int(z)
	}

	return e, nil
}
