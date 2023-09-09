package api

import (
	"context"
	"math/rand"
	"net/http"
	"gotrains/userpassenger_web/user-web/forms"
	"gotrains/userpassenger_web/user-web/global"
	"strconv"
	"strings"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

func GenerateSmsCode(width int) string {
	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	for i := 0; i < width; i++ {
		sb.WriteString(strconv.Itoa(rand.Intn(10)))
	}
	return sb.String()
}

func SendSms(c *gin.Context) {

	sendSmsForm := forms.SendSmsForm{}
	if err := c.ShouldBind(&sendSmsForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	client, _err := CreateClient(&global.ServerConfig.AliSmsConfig.AccessKeyId, &global.ServerConfig.AliSmsConfig.AccessSecret)
	if _err != nil {
		zap.S().Errorw("Create Sms Client failed", "msg", _err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "发送短信失败",
		})
		return
	}
	code := GenerateSmsCode(global.ServerConfig.AliSmsConfig.CodeLen)

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String(global.ServerConfig.AliSmsConfig.SignName),
		TemplateCode:  tea.String(global.ServerConfig.AliSmsConfig.TemplateCode),
		PhoneNumbers:  tea.String(sendSmsForm.Mobile),
		TemplateParam: tea.String("{\"code\":\"" + code + "\"}"),
	}
	runtime := &util.RuntimeOptions{}
	msg, _err := client.SendSmsWithOptions(sendSmsRequest, runtime)
	if _err != nil {
		zap.S().Errorw("Send Sms failed", "msg", _err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "发送短信失败",
		})
		return
	}
	zap.S().Infow("Send Sms success", "msg", msg)
	// 如果有多个不同Redis数据库，可以单独创建一个Redis客户端
	// 因为懒，只有一个Redis数据库，所以直接用全局变量
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr: global.ServerConfig.RedisConfig.Host + ":" + strconv.Itoa(global.ServerConfig.RedisConfig.Port),
	// })
	global.Rdb.Set(context.Background(), sendSmsForm.Mobile, code, time.Duration(global.ServerConfig.AliSmsConfig.Expire)*time.Second)
	c.JSON(http.StatusOK, gin.H{
		"message": "发送短信成功",
	})
}
