package worldedit

func GetBoundaries(entries []*WEEntry) (int, int, int) {
	max_x := 0
	max_y := 0
	max_z := 0
	for _, e := range entries {
		if e.X > max_x {
			max_x = e.X
		}
		if e.Y > max_y {
			max_y = e.Y
		}
		if e.Z > max_z {
			max_z = e.Z
		}
	}

	return max_x, max_y, max_z
}
