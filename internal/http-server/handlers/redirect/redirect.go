package redirect

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "url-service/internal/lib/api/response"
	"url-service/internal/storage"
)

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.With(slog.String("request_id", middleware.GetReqID(r.Context())))

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Info("alias is empty")
			render.JSON(w, r, resp.ErrorResponse("alias is empty"))
			return
		}

		url, err := urlGetter.GetURL(alias)

		if errors.Is(err, storage.ErrURLNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			log.Info("url not found")
			render.JSON(w, r, resp.ErrorResponse("URL not found"))
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("failed to get url", slog.String("error", err.Error()))
			render.JSON(w, r, resp.ErrorResponse("Internal server error"))
			return
		}

		log.Info("url found", slog.String("url", url))
		http.Redirect(w, r, url, http.StatusFound)
	}
}
