package worldedit

import "fmt"

type WEEntry struct {
	Name   string
	X      int
	Y      int
	Z      int
	Param1 int
	Param2 int
}

func (e WEEntry) String() string {
	return fmt.Sprintf(
		"{WEEntry Name='%s',X=%d,Y=%d,Z=%d,Param1=%d,Param2=%d}",
		e.Name, e.X, e.Y, e.Z, e.Param1, e.Param2,
	)
}
