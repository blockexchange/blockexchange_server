package types

func getInt64(o any) int64 {
	v, _ := o.(float64)
	return int64(v)
}

func getInt(o any) int {
	v, _ := o.(float64)
	return int(v)
}

func getString(o any) string {
	s, _ := o.(string)
	return s
}

func getBool(o any) bool {
	v, _ := o.(bool)
	return v
}
