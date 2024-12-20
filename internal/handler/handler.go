package handler

import (
	"disapp/internal/config"
	"disapp/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Dependences struct {
	AssetsFS       http.FileSystem
	Msgs           storage.Messages
	Config         config.Config
	ProjectBaseDir string
}

type handlerFunc func(http.ResponseWriter, *http.Request) error

func RegisterRoutes(router chi.Router, deps Dependences) {
	messageHandler := messageHandler{
		msgs:   deps.Msgs,
		domain: deps.Config.Address,
	}

	router.Handle("/assets/*", http.StripPrefix("/assets", http.FileServer(deps.AssetsFS)))
	router.Get("/", handler(homeHandler{}.handleIndex))
	router.Get("/m/{uuid}", handler(messageHandler.handleMessageView))
	router.Post("/api/messages", handler(messageHandler.createMessage))
}

func handler(h handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			handleError(w, r, err)
		}
	}
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	_ = w
	_ = r
	slog.Error("error during request", slog.String("err", err.Error()))
}
