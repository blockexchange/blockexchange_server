package types

type Login struct {
	Username string `json:"name"`
	Token    string `json:"access_token"`
}
