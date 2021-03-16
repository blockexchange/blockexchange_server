package types

func getInt64(o interface{}) int64 {
	v, _ := o.(float64)
	return int64(v)
}

func getInt(o interface{}) int {
	v, _ := o.(float64)
	return int(v)
}

func getString(o interface{}) string {
	s, _ := o.(string)
	return s
}

func getBool(o interface{}) bool {
	v, _ := o.(bool)
	return v
}
