package types

type Tag struct {
	UID         string `json:"uid" ksql:"uid"`
	Name        string `json:"name" ksql:"name"`
	Description string `json:"description" ksql:"description"`
	Restricted  bool   `json:"restricted" ksql:"restricted"`
}
