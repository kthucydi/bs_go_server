package server

import (
	gmux "github.com/gorilla/mux"
	logging "github.com/kthucydi/bs_go_logrus"
	mw "github.com/kthucydi/bs_go_server/middleware"
	"net/http"
	"os"
	"os/signal"
)

var (
	BackServer = &BackServerType{}
	Logger     = &logging.Log
)

// Run function sets endpoints and runs the server
func (backServer *BackServerType) Run(cfg map[string]string, API APISettings) {

	backServer.Init(cfg, API)
	Logger.Print("from server: Init success")

	Logger.Print("from ListenAndServe:set endpoint success, run")

	go func() {
		if err := backServer.srv.ListenAndServe(); err != nil {
			Logger.Error("from ListenAndServe:", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func (backServer *BackServerType) Init(cfg map[string]string, API APISettings) {
	mw.InitConfig(cfg)
	backServer.mw = mw.Middleware
	backServer.gmux = gmux.NewRouter()
	backServer.api = API
	backServer.Cfg = cfg
	backServer.srv = &http.Server{
		Addr:    ":" + backServer.Cfg["BACKEND_SERVER_PORT"],
		Handler: backServer.gmux,
	}
	backServer.SetEndPoint(API)
}

// SetEndPoint Set endpoint from api config structure
func (backServer *BackServerType) SetEndPoint(API APISettings) {

	// Set PathPrefix for URL
	mux := backServer.gmux.PathPrefix(backServer.Cfg["BACKEND_SERVER_URL_PREFIX"]).Subrouter()

	for name, _ := range backServer.mw {
		Logger.Printf("MW: %s loaded", name)
	}
	// Set endpoints
	for path, methods := range API.GetRouteList() {
		for method, routeConfig := range methods {
			backServer.CreateRoute(mux, path, method, routeConfig)
		}
	}

	// Print Allowed methods in headers
	if value, ok := backServer.Cfg["USE_INNER_CORS"]; ok && value == "true" {
		mux.Use(backServer.mw["innerCORSAllow"])
		mux.Use(gmux.CORSMethodMiddleware(mux))
	}

	if value, ok := backServer.Cfg["USE_INNER_LOGGER"]; ok && value == "true" {
		mux.Use(backServer.mw["innerLogger"])
	}

	// Register common middleware
	CommonMiddleware := API.GetCommonMiddleware()
	if len(CommonMiddleware) > 0 {
		for _, middleware := range CommonMiddleware {
			mux.Use(gmux.MiddlewareFunc(middleware))
		}
	}
}

// CreateRoute creating end register new route
func (backServer *BackServerType) CreateRoute(mux *gmux.Router, path, method string, route Route) {

	// Create new handler
	var finalHandler http.Handler
	finalHandler = http.HandlerFunc(route.Handler)

	// Accept Auth middleware
	if auth := route.Auth; auth != "" {
		if authFunc, ok := backServer.mw[auth]; ok {
			finalHandler = authFunc(finalHandler)
		} else if authFunc, ok := backServer.api.GetAuthMiddleware()[auth]; ok {
			finalHandler = authFunc(finalHandler)
		} else {
			Logger.Fatalf("can not find Auth middleware for %s - %s", path, method)
		}
	}

	// Accept other middlewares from route config
	if len(route.Middlewares) > 0 {
		for i := len(route.Middlewares) - 1; i >= 0; i-- {
			finalHandler = route.Middlewares[i](finalHandler)
		}
	}

	//Handle inner and out common middlewares

	// Register new handler
	mux.Methods(method).Path(path).Handler(finalHandler).Name(route.Name)
}
