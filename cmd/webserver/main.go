package main

import (
	"flag"
	"log"
	"net/http"
	"path"

	"github.com/go-openapi/loads"
	"github.com/markbates/pop"
	"go.uber.org/zap"
	"goji.io"
	"goji.io/pat"

	"github.com/transcom/mymove/pkg/gen/restapi"
	"github.com/transcom/mymove/pkg/gen/restapi/operations"
	issueop "github.com/transcom/mymove/pkg/gen/restapi/operations/issues"
	"github.com/transcom/mymove/pkg/handlers"
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

	build := flag.String("build", "build", "the directory to serve static files from.")
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, configures the database, presenetly.")
	port := flag.String("port", ":8080", "the `port` to listen on.")
	swagger := flag.String("swagger", "swagger.yaml", "The location of the swagger API definition")
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

	//DB connection
	pop.AddLookupPaths(*config)
	dbConnection, err := pop.Connect(*env)
	if err != nil {
		log.Panic(err)
	}

	// initialize api pkg with dbConnection created above
	handlers.Init(dbConnection)

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewMymoveAPI(swaggerSpec)

	api.IssuesCreateIssueHandler = issueop.CreateIssueHandlerFunc(handlers.CreateIssueHandler)
	api.IssuesIndexIssuesHandler = issueop.IndexIssuesHandlerFunc(handlers.IndexIssuesHandler)

	// Serves files out of build folder
	clientHandler := http.FileServer(http.Dir(*build))

	// Base routes
	root := goji.NewMux()
	root.Handle(pat.Get("/api/v1/swagger.yaml"), fileHandler(*swagger))
	root.Handle(pat.Get("/api/v1/docs"), fileHandler(path.Join(*build, "swagger-ui", "index.html")))
	root.Handle(pat.New("/api/*"), api.Serve(nil)) // Serve(nil) returns an http.Handler for the swagger api
	root.Handle(pat.Get("/static/*"), clientHandler)
	root.Handle(pat.Get("/swagger-ui/*"), clientHandler)
	root.Handle(pat.Get("/favicon.ico"), clientHandler)
	root.HandleFunc(pat.Get("/*"), fileHandler(path.Join(*build, "index.html")))

	// And request logging
	root.Use(requestLogger)

	zap.L().Info("Starting the server listening", zap.String("port", *port))
	log.Fatal(http.ListenAndServe(*port, root))
}

// fileHandler serves up a single file
func fileHandler(entrypoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}
}
