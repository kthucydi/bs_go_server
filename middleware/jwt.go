package middleware

import (
	jwtmw "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
)

var cfg map[string]string

var Jwt = jwtmw.New(jwtmw.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg["JWT_SECRET"]), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

func InitConfig(config map[string]string) {
	cfg = config
}
