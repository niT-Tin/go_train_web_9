package api

import (
	"context"
	"net/http"
	"gotrains/userpassenger_web/user-web/forms"
	"gotrains/userpassenger_web/user-web/global"
	"gotrains/userpassenger_web/user-web/global/response"
	"gotrains/userpassenger_web/user-web/middlewares"
	"gotrains/userpassenger_web/user-web/models"
	"gotrains/userpassenger_web/user-web/proto"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RemoveTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err == nil {
		return
	}
	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "参数错误",
			})
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"message": e.Message(),
			})
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "内部错误",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "其他错误",
			})
		}
		return
	}
	zap.S().Errorw("grpc 请求失败", "msg", err.Error())
	c.JSON(500, gin.H{
		"message": err.Error(),
	})
}

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": RemoveTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func GetUserList(c *gin.Context) {
	zap.S().Debug("获取用户列表")
	claims, _ := c.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户: %d", currentUser.ID)

	userSrvClient := global.UserClient
	pn := c.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	ps := c.DefaultQuery("ps", "10")
	psInt, _ := strconv.Atoi(ps)
	ulr, err2 := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn: uint32(pnInt),
		Ps: uint32(psInt),
	})
	if err2 != nil {
		zap.S().Errorw("GetUserList 获取用户列表失败", "msg", err2.Error())
		HandleGrpcErrorToHttp(err2, c)
		return
	} else {
		result := make([]response.UserResponse, 0)
		for _, value := range ulr.Data {
			user := response.UserResponse{
				Id:       value.Id,
				NickName: value.NickName,
				Birthday: response.JsonTime(time.Unix(int64(value.BirthDay), 0)),
				Gender:   value.Gender,
				Mobile:   value.Mobile,
			}
			result = append(result, user)
		}
		c.JSON(http.StatusOK, result)
	}
}

func PasswordLogin(c *gin.Context) {
	// 表单验证
	plf := forms.PasswordLoginForm{}
	if err := c.ShouldBind(&plf); err != nil {
		HandleValidatorError(c, err)
		return
	}

	if !store.Verify(plf.CaptchaId, plf.Captcha, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "验证码错误",
		})
		return
	}
	userSrvClient := global.UserClient

	if rsp, err := userSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: plf.Mobile,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.NotFound {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "用户不存在",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "登录失败",
			})
			return
		}
	} else {
		if chkpass, passErr := userSrvClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          plf.Password,
			EncryptedPassword: rsp.Password,
		}); passErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "登录失败",
			})
			return
		} else {
			if chkpass.Success {
				// 生成token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),               // 签名生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30, // 过期时间 30天
						Issuer:    "liuzehao",                      // 签名的发行者
					},
				}
				token, err := j.CreateToken(claims)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "生成token失败",
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"id":         rsp.Id,
					"token":      token,
					"nick_name":  rsp.NickName,
					"message":    "登录成功",
					"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
				})
				return
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "密码错误",
				})
				return
			}
		}
	}
}

func Register(c *gin.Context) {
	registerForm := forms.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(c, err)
		return
	}
	// 从redis中取出验证码
	code, err := global.Rdb.Get(context.Background(), registerForm.Mobile).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "验证码错误",
		})
		return
	} else {
		if code != registerForm.Code {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "验证码错误",
			})
			return
		}
	}
	userSrvClient := global.UserClient
	u, e := userSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: registerForm.Mobile,
		Password: registerForm.Password,
		Mobile:   registerForm.Mobile,
	})
	if err != nil {
		zap.S().Errorw("Register 注册失败", "msg", e.Error())
		HandleGrpcErrorToHttp(e, c)
		return
	}
	// 注册立即登录
	// 生成token
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(u.Id),
		NickName:    u.NickName,
		AuthorityId: uint(u.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               // 签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, // 过期时间 30天
			Issuer:    "liuzehao",                      // 签名的发行者
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "生成token失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":         u.Id,
		"token":      token,
		"nick_name":  u.NickName,
		"message":    "登录成功",
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
	})
	return
}
