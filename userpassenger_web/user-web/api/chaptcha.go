package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

var store = base64Captcha.DefaultMemStore

// CaptchaGet 获取验证码
func CaptchaGet(c *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := cp.Generate()
	if err != nil {
		zap.S().Errorw("生成验证码失败", "msg", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "获取验证码失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"captchaId": id,
		"picPath":   b64s,
	})
}

func CaptchaVerify(c *gin.Context) {
	var captcha struct {
		CaptchaId string `json:"captchaId" binding:"required"`
		Value     string `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&captcha); err != nil {
		zap.S().Errorw("参数错误", "msg", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
			"success": false,
		})
		return
	}
	if !store.Verify(captcha.CaptchaId, captcha.Value, true) {
		c.JSON(http.StatusOK, gin.H{
			"message": "验证码错误",
			"success": false,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "验证码正确",
		"success": true,
	})
}
