package types

type AccessToken struct {
	ID       int64  `json:"id" ksql:"id"`
	Name     string `json:"name" ksql:"name"`
	Token    string `json:"token" ksql:"token"`
	UserID   int64  `json:"user_id" ksql:"user_id"`
	Created  int64  `json:"created" ksql:"created"`
	Expires  int64  `json:"expires" ksql:"expires"`
	UseCount int    `json:"usecount" ksqp:"usecount"`
}
