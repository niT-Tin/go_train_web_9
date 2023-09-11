package forms

type PasswordLoginForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,len=11,mobile"`
	// Username string `form:"username" json:"username" binding:"required,min=3,max=20"`
	Password  string `form:"password" json:"password" binding:"required,min=6,max=20"`
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"`
	CaptchaId string `form:"captchaId" json:"captchaId" binding:"required"`
}

type RegisterForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,len=11,mobile"`
	// Username string `form:"username" json:"username" binding:"required,min=3,max=20"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=20"`
	Code     string `form:"code" json:"code" binding:"required,min=5,max=5"`
}
