package create

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
)

type Request struct {
	Url string `json:"url" validate:"required,url"`
}

type Response struct {
	response.Response
	Alias string `json:"alias"`
}

type UrlSaver interface {
	CreateUrl(urlValue string) (string, error)
}

func New(log *slog.Logger, saver UrlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.create"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if errors.Is(err, io.EOF) {
			log.Error("Request body id empty")

			render.JSON(w, r, response.Error("Empty request"))

			return
		}

		if err != nil {
			log.Error("Failed to decode request body", err)

			render.JSON(w, r, response.Error("Failed to decode request body"))

			return
		}

		log.Info("Request body decoded", slog.Any("req", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("Invalid request", err)
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			render.JSON(w, r, response.Error(validateErr.Error()))

			return
		}

		alias, err := saver.CreateUrl(req.Url)

		if err != nil {
			log.Error("Failed to add URL", err)

			render.JSON(w, r, response.Error("Failed to add URL"))

			return
		}

		log.Error("Url created", alias)

		render.JSON(w, r, Response{Alias: alias})

		return
	}
}
