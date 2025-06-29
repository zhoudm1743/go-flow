package req

type LoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Captcha  string `json:"captcha" validate:"required"`
	AppID    string `json:"app_id" validate:"required"`
}
