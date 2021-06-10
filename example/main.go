package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/cbroglie/mustache"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/maverickames/omeh"
)

// App Main App controller
type App struct {
	LogFile *os.File

	appName  string
	debug    bool
	router   *chi.Mux
	em       *omeh.ErrManager
	log      *log.Logger
	logError *log.Logger
	weblog   *log.Logger
}

func main() {
	app := App{}
	app.debug = true
	app.appName = "Example App"
	app.LogFile = os.Stdout
	// Setup Logging
	app.log = log.New(app.LogFile, "["+app.appName+"] Logger: ", log.Ldate|log.Ltime)
	app.logError = log.New(app.LogFile, "["+app.appName+"] ErrorHandler: ", log.Ldate|log.Ltime)
	app.weblog = log.New(app.LogFile, "["+app.appName+"] WebLog: ", log.Ldate|log.Ltime)

	// Setup Error Handling for routes
	app.em = omeh.New(app.debug)

	// Setup the default processor for error
	app.em.SetDefaultHandler(app.em.ReturnError)
	// Or pass a custom error creation function
	//app.em.SetDefaultHandler(ReturnError)

	// Setup the Error Logger
	// this is your operertunity to process errors however you wants
	app.em.SetErrorLogHandler(app.logErrorHandler)

	// ProcessError you can use the default or build a custom one
	// This example shows how to return a application layer response
	app.em.SetHTTPErrorHandler(func(h omeh.ErrorHandlerdesc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ErrResponse := h(w, r)
			if ErrResponse != nil {
				app.em.LogError(ErrResponse)
				renderdView, err := mustache.RenderFile("errorResponse.html.mustache", ErrResponse)
				if err != nil {
					jsonData, err := json.Marshal(ErrResponse)
					if err != nil {
						app.logError.Println(err)
						return
					}
					w.WriteHeader(ErrResponse.HTTPStatusCode)
					w.Write([]byte(jsonData))
					return
				}
				w.Write([]byte(renderdView))
			}
		})
	})

	// Create Router
	app.router = chi.NewRouter()

	// Add Middleware
	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: app.weblog, NoColor: true})
	app.router.Use(middleware.Logger)
	app.router.Use(render.SetContentType(render.ContentTypeJSON))

	// Set NotFound route
	app.router.NotFound(app.em.HandleHTTPErrors(app.Render404notFound))

	// Set basic route
	app.router.Get("/", app.em.HandleHTTPErrors(app.root))
	app.router.Get("/broken", app.em.HandleHTTPErrors(app.brokenRoute))

	// Log an error in the application
	err := errors.New("Server broke doing stuff.")
	if err != nil {
		app.em.LogError(app.em.ErrorHandler(err, omeh.StatusInternalServerError, "Addition detail if desired."))
	}

	// Start webserver with router
	app.log.Println("Server runing from " + "localhost:5000")
	err = http.ListenAndServe("localhost:5000", app.router)
	if err != nil {
		app.weblog.Fatal(err)
	}
}

// logError Handle all errors received from Router
func (app *App) logErrorHandler(ErrResponse *omeh.ErrResponse) {
	if app.debug {
		app.logError.Printf(
			"\n  -- Function: %s\n  -- SourceFile: %s\n  -- LineNumber: %d\n  -- ErrorDetails: %v\n  -- RequestDetail: %s\n  -- ErrorCode: %d\n",
			runtime.FuncForPC(ErrResponse.FuncPC).Name(),
			ErrResponse.FuncFN,
			ErrResponse.FuncLine,
			ErrResponse.Err,
			ErrResponse.RequestDetail,
			ErrResponse.AppCode,
		)
	} else {
		app.logError.Printf(
			"\n  -- ErrorDetails: %v\n  -- RequestDetail: %s\n  -- ErrorCode: %d\n",
			ErrResponse.Err,
			ErrResponse.RequestDetail,
			ErrResponse.AppCode,
		)
	}
}

func ReturnError(err error, errResp *omeh.ErrResponse, reqDetial string) *omeh.ErrResponse {

	return omeh.StatusNotFound
}
