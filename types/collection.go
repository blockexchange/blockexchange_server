package types

type Collection struct {
	UID         string `json:"uid" ksql:"uid"`
	UserUID     string `json:"user_uid" ksql:"user_uid"`
	Name        string `json:"name" ksql:"name"`
	Description string `json:"description" ksql:"description"`
}
