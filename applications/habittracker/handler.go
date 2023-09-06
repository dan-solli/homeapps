package trackerhttp

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/quii/todo/adapters/todohttp/views"
)

var (
	//go:embed "templates/*"
	trackerTemplates embed.FS
)

type TrackerHandler struct {
	http.Handler

	trackerView *views.ModelView[tracker.Tracker]
	indexView   *views.IndexView
}

func NewTrackerHandler(service *tracker.List, trackerView *views.ModelView[tracker.Tracker], indexView *views.IndexView) (*TrackerHandler, error) {
	router := mux.NewRouter()
	handler := &TrackerHandler{
		Handler:     router,
		list:        service,
		trackerView: trackerView,
		indexView:   indexView,
	}

	staticHandler, err := newStaticHandler()
	if err != nil {
		return nil, fmt.Errorf("problem making static resources handler: %w", err)
	}

	router.HandleFunc("/", handler.index).Methods(http.MethodGet)

	router.HandleFunc("/todos", handler.add).Methods(http.MethodPost)
	router.HandleFunc("/todos", handler.search).Methods(http.MethodGet)
	router.HandleFunc("/todos/sort", handler.reOrder).Methods(http.MethodPost)
	router.HandleFunc("/todos/{ID}/edit", handler.edit).Methods(http.MethodGet)
	router.HandleFunc("/todos/{ID}/toggle", handler.toggle).Methods(http.MethodPost)
	router.HandleFunc("/todos/{ID}", handler.delete).Methods(http.MethodDelete)
	router.HandleFunc("/todos/{ID}", handler.rename).Methods(http.MethodPatch)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticHandler))

	return handler, nil
}
