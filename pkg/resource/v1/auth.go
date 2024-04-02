package v1

/*
For apiserver struct
*/

type AuthUserLoginReq struct {
	Username  string `json:"username"`   //	用户名
	Password  string `json:"password"`   //	密码
	CaptchaId string `json:"captcha_id"` // 验证码ID
	Captcha   string `json:"captcha"`    //	验证码
}

type AuthData struct {
	Roles  map[string]string `json:"roles"`
	Token  string            `json:"token"`
	Status AuthStatus        `json:"status"`
}

type AuthStatus struct {
	AuthSts string `json:"authSts"`
}

/*
For controller backend struct
*/

type RestAuthData struct {
	Password *RestAuthPassword `json:"password,omitempty"`
	Token    *RestAuthToken    `json:"Token,omitempty"`
	ClientIP string            `json:"client_ip"`
}

type RestAuthPassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RestAuthToken struct {
	Token    string `json:"token"`
	State    string `json:"state"`
	Redirect string `json:"redirect_endpoint"`
}
