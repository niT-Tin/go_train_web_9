package forms

type SendSmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,len=11,mobile"`
	Type   string `form:"type" json:"type" binding:"required,oneof=1 2"`
}
