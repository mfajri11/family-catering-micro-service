package HandlerHttp

type LoginRequest struct {
	Email           string `json:"username" validate:"required,email"`
	Password        string `json:"password" validate:"required,alphanumunicode,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,alphanumunicode,min=8"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
