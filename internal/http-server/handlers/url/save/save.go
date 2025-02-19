package save

import (
	"errors"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"log/slog"
	"net/http"
	resp "url-service/internal/lib/api/response"
	"url-service/internal/lib/random"
	"url-service/internal/storage"
)

type Request struct {
	URL string `json:"url" validate:"required,url"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=URLSaver
type URLSaver interface {
	SaveURL(url, alias string) error
}

const aliasLength = 10

func New(logger *slog.Logger, saver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.With(slog.String("request_id", middleware.GetReqID(r.Context())))
		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("failed to decode request", slog.String("error", err.Error()))
			render.JSON(w, r, resp.ErrorResponse("failed to decode request"))
			return
		}

		log.Info("request decoded", slog.Any("request", req))

		if err = validator.New().Struct(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			log.Error("Invalid request", slog.String("error", err.Error()))
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		alias := random.NewRandomString(aliasLength)

		err = saver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			w.WriteHeader(http.StatusBadRequest)
			log.Info("URL already exists", slog.String("url", req.URL))
			render.JSON(w, r, resp.ErrorResponse("URL already exists"))
			return
		}

		if errors.Is(err, storage.ErrAliasExists) {
			w.WriteHeader(http.StatusInternalServerError)
			log.Info("Alias already exists", slog.String("alias", alias))
			render.JSON(w, r, resp.ErrorResponse("Error generating alias, try again"))
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("failed to save url", slog.String("error", err.Error()))
			render.JSON(w, r, resp.ErrorResponse("failed to save url"))
			return
		}

		w.WriteHeader(http.StatusCreated)
		log.Info("saved url", slog.String("url", req.URL))
		render.JSON(w, r, Response{
			Response: *resp.OkResponse(),
			Alias:    alias,
		})
	}
}
