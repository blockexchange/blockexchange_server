package worldedit

import "fmt"

type WEFields map[string]string
type WEInventory map[string][]string

type WEMeta struct {
	Fields    WEFields
	Inventory WEInventory
}

type WEEntry struct {
	Name   string
	X      int
	Y      int
	Z      int
	Param1 byte
	Param2 byte
	Meta   *WEMeta
}

func (e WEEntry) String() string {
	return fmt.Sprintf(
		"{WEEntry Name='%s',X=%d,Y=%d,Z=%d,Param1=%d,Param2=%d}",
		e.Name, e.X, e.Y, e.Z, e.Param1, e.Param2,
	)
}
