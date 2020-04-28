package main

import (
	legacyLog "log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/samedguener/ImageService/handlers"
	"github.com/samedguener/ImageService/middleware"
	"github.com/samedguener/ImageService/utils"
	"github.com/sirupsen/logrus"

	ghandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	utils.InitEnvironmentVariables()
	utils.InitGCPCloudStorageBucket()
	initLogger()

	root := mux.NewRouter()

	ghandlers.LoggingHandler(os.Stdout, root)

	/*** Un-Secured ***/
	api := root.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/health", handlers.Health.Get).Methods(http.MethodGet)
	api.HandleFunc("/images/endpoint", handlers.Images.GetImageAccessEndpoint).Methods(http.MethodGet)

	/*** Secured ***/
	securedAPI := root.PathPrefix("/api/v1").Subrouter()
	securedAPI.Use(middleware.VerifyToken)
	securedAPI.HandleFunc("/images/endpoint", handlers.Images.GetImageAccessEndpoint).Methods(http.MethodGet)
	securedAPI.HandleFunc("/images", handlers.Images.Post).Methods(http.MethodPost)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
		logrus.Printf("Defaulting to port %s", port)
	}

	logrus.Printf("Listening on port %s", port)

	http.Handle("/", &Server{root})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logrus.Fatal(err)
	}
}

// Server ...
type Server struct {
	r *mux.Router
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	auth := req.Header.Get("Authorization")
	if len(auth) > 0 {
		rw.Header().Set("Authorization", auth)
	}

	rw.Header().Set("Cache-Control", "no-cache")

	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", strings.Join([]string{http.MethodPost, http.MethodGet, http.MethodOptions, http.MethodPut, http.MethodDelete}, ", "))
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}

	if req.Method == http.MethodOptions {
		return
	}

	s.r.ServeHTTP(rw, req)
}

func initLogger() {
	defer func() {
		if p := recover(); p != nil {
			legacyLog.Fatal("Error while init logger")
		}
	}()

	var skipPath string
	//var skipPkg string
	_, file, _, ok := runtime.Caller(1)
	if ok {
		path := strings.Split(file, "/")
		skipPath = strings.Join(path[:len(path)-1], "/")
	}

	// Log as JSON instead of the default ASCII formatter.
	if utils.Environment.Value == "dev" {
		beautifier := func(frame *runtime.Frame) (function string, file string) {
			fx := strings.SplitAfter(frame.Func.Name(), ".")
			fn := strings.SplitAfter(frame.File, skipPath)

			return fx[len(fx)-1] + "()", fn[len(fn)-1] + ":" + strconv.Itoa(frame.Line)
		}

		logrus.SetOutput(os.Stdout)
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, PadLevelText: true, FullTimestamp: true, CallerPrettyfier: beautifier})
	} else {
		logrus.SetOutput(os.Stdout)
		logrus.SetLevel(logrus.WarnLevel)
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	logrus.SetReportCaller(true)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.DebugLevel)
}
