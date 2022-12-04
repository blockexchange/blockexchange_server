package types

import "strconv"

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

type JsonInt int

func (ji *JsonInt) UnmarshalJSON(data []byte) error {
	var f float64
	f, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return err
	}
	i := int(f)
	*ji = JsonInt(i)
	return nil
}

type JsonInt64 int64

func (ji *JsonInt64) UnmarshalJSON(data []byte) error {
	var f float64
	f, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return err
	}
	i := int(f)
	*ji = JsonInt64(i)
	return nil
}
