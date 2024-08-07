package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	gmux "github.com/gorilla/mux"
	logging "github.com/kthucydi/bs_go_logrus"
	mw "github.com/kthucydi/bs_go_server/middleware"
)

var (
	BackServer = &BackServerType{}
	Logger     = &logging.Log
)

// Run function sets endpoints and runs the server
func (backServer *BackServerType) Run(cfg map[string]string, API APISettings) {

	backServer.Init(cfg, API)
	Logger.Infof("server: set endpoint success, try runing at %s port", backServer.Cfg["BACKEND_SERVER_PORT"])

	go func() {
		if err := backServer.srv.ListenAndServe(); err != http.ErrServerClosed {
			Logger.Fatalf("(exit) error server: ListenAndServe: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	Logger.Info("HTTP server gracefull Shutdown by interrupt signal")
	if err := backServer.srv.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout:
		Logger.Infof("HTTP server Shutdown: %v", err)
	}
}

// Run function sets endpoints and runs the server
func (backServer *BackServerType) RunGracefullShutdown(cfg map[string]string, API APISettings) {

	backServer.Init(cfg, API)
	Logger.Infof("server: set endpoint success, try to runing at %s port", backServer.Cfg["BACKEND_SERVER_PORT"])

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		Logger.Info("HTTP server gracefull Shutdown by interrupt signal")
		if err := backServer.srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			Logger.Infof("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := backServer.srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		Logger.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
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
	var mux *gmux.Router

	if backServer.Cfg["BACKEND_SERVER_URL_PREFIX"] != "/" {
		// Set PathPrefix for URL
		mux = backServer.gmux.PathPrefix(backServer.Cfg["BACKEND_SERVER_URL_PREFIX"]).Subrouter()
	} else {
		mux = backServer.gmux.PathPrefix("").Subrouter()
	}
	// Set endpoints
	for path, methods := range API.GetRouteList() {
		for method, routeConfig := range methods {
			backServer.CreateRoute(mux, path, method, routeConfig)
		}
	}

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

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
