package types

type Tag struct {
	ID          *int64 `json:"id" ksql:"id"`
	Name        string `json:"name" ksql:"name"`
	Description string `json:"description" ksql:"description"`
	Restricted  bool   `json:"restricted" ksql:"restricted"`
}
