package types

type ValidationResult struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}
