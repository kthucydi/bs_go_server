
Модуль bs_go_server - внутренний модуль, созданный чтобы упростить конфигурацию и запуск сервера приложения. 

# Настройка

Все настройки сервер берет из map[string]string и интерфейса, которые передаются ему при запуске
Для работы должны быть установлены следующие параметры:
```
JWT_TOKEN_DURATION=			# продолжительность действия токена в часах
JWT_SECRET=					# секретная строка для jwt
BACKEND_SERVER_PORT=		# порт на котором будет работать сервер
BACKEND_SERVER_URL_PREFIX=	# префикс пути для api for example: /api/v1
USE_INNER_CORS=				# Использовать встроенный middleware CORS
USE_INNER_LOGGER=			# использовать встроенный middleware логгер
```

На вход кроме параметров сервер принимает интерфейс ApiConfig:
```go
type APIConfig interface {
	GetRouteList() map[string]Methods
	GetCommonMiddleware() map[string]gmux.MiddlewareFunc
	GetAuthMiddleware() map[string]gmux.MiddlewareFunc
}
```
где GetRouteList() выдает нам структуру со всеми роутами,\
а GetCommonMiddleware() и GetAuthMiddleware() соответственно добавляет общие для всех middleware и middleware для аутентификации \
gmux - это горилла mux


# Пример
```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	gmux "github.com/gorilla/mux"
	// подключаем модуль
	"github.com/kthucydi/bs_go_server"
)

type API struct {
	Routes map[string]server.Methods
	CommonMiddleware map[string]gmux.MiddlewareFunc
	AuthMiddleware map[string]gmux.MiddlewareFunc
}

func (a *API) GetRouteList() (map[string]server.Methods) {
	return a.Routes
}

func (a *ApiRoutes) GetCommonMiddleware() map[string]gmux.MiddlewareFunc {
	return a.CommonMiddleware
}

func (a *ApiRoutes) GetAuthMiddleware() map[string]gmux.MiddlewareFunc {
	return a.AuthMiddleware
}

func main() {
	// for example, better, if use environment variable
	config := map[string]string{
		"JWT_TOKEN_DURATION":		"72",		
		"JWT_SECRET":				"secret"
		"BACKEND_SERVER_PORT":		"5001"
		"BACKEND_SERVER_URL_PREFIX":"/api/v1"
		"USE_INNER_CORS":			"true"
		"USE_INNER_LOGGER":			"true"
	}

	api := API{
		Routes: make(map[string]server.Methods){
		"/login":	Methods{
			"GET":     Route{Name: "login", Auth: "", Handler: listener},
			}
		},
		CommonMiddleware: make(map[string]gmux.MiddlewareFunc)
		AuthMiddleware make(map[string]gmux.MiddlewareFunc)
	}

	// запускаем сервер
	server.BackServer.Run(config, api)
	
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	sendJson(w, http.StatusOK, "you are login now")
	log.Println("someone login")
}

func sendJson(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{
		Message: msg,
	})
}

```
