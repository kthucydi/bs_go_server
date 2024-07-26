package server

import (
	gmux "github.com/gorilla/mux"
	"net/http"
)

type Methods map[string]Route

type Route struct {
	Auth        string
	Name        string
	Handler     http.HandlerFunc
	Middlewares []gmux.MiddlewareFunc
}

type APISettings interface {
	GetRouteList() map[string]Methods
	GetCommonMiddleware() map[string]gmux.MiddlewareFunc
	GetAuthMiddleware() map[string]gmux.MiddlewareFunc
}

type BackServerType struct {
	gmux *gmux.Router
	srv  *http.Server
	Cfg  map[string]string
	api  APISettings
	mw   map[string]gmux.MiddlewareFunc
}
