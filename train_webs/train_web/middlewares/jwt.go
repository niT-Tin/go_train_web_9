package middlewares

import (
	"errors"
	"net/http"
	"gotrains/train_webs/train_web/global"
	"gotrains/train_webs/train_web/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type JWT struct {
	SingingKey []byte
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.ServerConfig.JWTInfo.SigningKey),
	}
}

func (j *JWT) CreateToken(claims models.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SingingKey)
}

func (j *JWT) ParseToken(tokenStr string) (claims *models.CustomClaims, err error) {
	// 此处传入&models.CustomClaims{} 不能只是传入一个这个类型的指针，而是要传入一个实例化的对象的指针
	token, err := jwt.ParseWithClaims(tokenStr, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SingingKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				err = TokenMalformed
				return
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				err = TokenExpired
				return
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				err = TokenNotValidYet
				return
			} else {
				err = TokenInvalid
				return
			}
		}
	}
	if token == nil {
		err = TokenInvalid
	}
	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "请求未携带token，无权限访问",
			})
			c.Abort()
			return
		}
		j := NewJWT()
		if claims, err := j.ParseToken(token); err != nil {
			if err == TokenExpired {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "授权已过期",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			c.Abort()
			return
		} else {
			c.Set("claims", claims)
			c.Set("userId", claims.ID)
			c.Next()
		}
	}
}
