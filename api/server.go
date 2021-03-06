package api

import (
	"net/http"

	"github.com/getsentry/raven-go"
)

var mux = http.NewServeMux()

// API defines socket on which we listen for commands
type API struct {
	HTTP *http.Server
}

// Registrator defines interface that help us to register http handler
type Registrator interface {
	GetName() string
	RegisterAPIHandles() map[string]func(w http.ResponseWriter, r *http.Request)
	DefaultAPIHandles() map[string]func(w http.ResponseWriter, r *http.Request)
}

// Register all available http handles
func (a *API) Register(r Registrator) *API {
	for p, f := range r.RegisterAPIHandles() {
		mux.HandleFunc("/"+r.GetName()+"/"+p, f)
	}

	for p, f := range r.DefaultAPIHandles() {
		mux.HandleFunc("/"+r.GetName()+"/"+p, f)
	}

	return a
}

// RegisterHandle append handle before server is started
func (a *API) RegisterHandle(e string, h func(w http.ResponseWriter, r *http.Request)) {
	mux.HandleFunc("/"+e, h)
}

// Start api server
func (a *API) Start() {
	a.HTTP.Handler = mux
	if err := a.HTTP.ListenAndServe(); err != nil {
		raven.CaptureError(err, nil)
		panic(err)
	}
}

// New creates new socket
func New() *API {
	api := &API{
		&http.Server{
			Addr: ":11401",
		},
	}

	return api
}
