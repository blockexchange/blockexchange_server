package types

type RegisterRequest struct {
	Name          string `json:"name"`
	Password      string `json:"password"`
	CaptchaID     string `json:"captcha_id"`
	CaptchaAnswer string `json:"captcha_answer"`
}

type CheckRegisterResponse struct {
	Success             bool `json:"success"`
	ErrPasswordTooShort bool `json:"error_password_too_short"`
	ErrCaptcha          bool `json:"error_captcha"`
	ErrInvalidUsername  bool `json:"error_invalid_username"`
	ErrUsernameTaken    bool `json:"error_username_taken"`
}
