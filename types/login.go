package types

type Login struct {
	Username string `json:"name"`
	Password string `json:"password"`
	Token    string `json:"access_token"`
}
