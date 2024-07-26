package middleware

import (
	gmux "github.com/gorilla/mux"
	logging "github.com/kthucydi/bs_go_logrus"
)

var (
	Logger     = &logging.Log
	Middleware = make(map[string]gmux.MiddlewareFunc)
)

func init() {
	Middleware["innerAuthJWT"] = Jwt.Handler
	Middleware["innerLogger"] = MiddlewareSmartLogger()
	Middleware["innerCORSAllow"] = MiddlewareCORSAllow()
}
