package types

type AccessToken struct {
	UID      string `json:"uid" ksql:"uid"`
	Name     string `json:"name" ksql:"name"`
	Token    string `json:"token" ksql:"token"`
	UserUID  string `json:"user_uid" ksql:"user_uid"`
	Created  int64  `json:"created" ksql:"created"`
	Expires  int64  `json:"expires" ksql:"expires"`
	UseCount int    `json:"usecount" ksqp:"usecount"`
}
