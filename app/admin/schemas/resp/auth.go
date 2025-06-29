package resp

type LoginResp struct {
	Token string `json:"token" structs:"token"`
}

type CaptchaResp struct {
	CaptchaID string `json:"captcha_id" structs:"captcha_id"`
	Captcha   string `json:"captcha" structs:"captcha"`
}
