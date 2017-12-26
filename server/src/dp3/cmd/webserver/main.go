package main

import (
	"dp3/pkg/api"
	"flag"
	"go.uber.org/zap"
	"goji.io"
	"goji.io/pat"
	"log"
	"net/http"
)

var logger *zap.Logger

// TODO(nick - 12/21/17) - this is a simple logger for debugging testing
// It needs replacing with something we can use in production
func requestLogger(h http.Handler) http.Handler {
	zap.L().Info("Request logger installed")
	wrapper := func(w http.ResponseWriter, r *http.Request) {
		zap.L().Info("Request", zap.String("url", r.URL.String()))
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(wrapper)
}

func main() {

	entry := flag.String("entry", "../client/build/index.html", "the entrypoint to serve.")
	static := flag.String("static", "../client/build/static", "the directory to serve static files from.")
	port := flag.String("port", ":8080", "the `port` to listen on.")
	debugLogging := flag.Bool("debug_logging", false, "log messages at the debug level.")
	flag.Parse()

	// Set up logger for the system
	var err error
	if *debugLogging {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	// api routes
	api := api.Mux()

	// Base routes
	root := goji.NewMux()
	root.Handle(pat.New("/api/*"), api)
	root.Handle(pat.Get("/static/*"),
		http.StripPrefix("/static", http.FileServer(http.Dir(*static))))
	root.HandleFunc(pat.Get("/*"), IndexHandler(entry))

	// And request logging
	root.Use(requestLogger)

	zap.L().Info("Starting the server listening", zap.String("port", *port))
	http.ListenAndServe(*port, root)
}

// IndexHandler serves up our index.html
func IndexHandler(entrypoint *string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, *entrypoint)
	}
}