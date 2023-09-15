package redirect

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/storage"
)

type UrlGetter interface {
	GetUrl(alias string) (string, error)
}

func New(log *slog.Logger, getter UrlGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("Alias is empty")

			render.JSON(w, r, response.Error("Alias is empty"))
			return
		}

		url, err := getter.GetUrl(alias)

		if errors.Is(err, storage.ErrorUrlNotFound) {
			log.Info("Url not found")

			render.JSON(w, r, response.Error("Url not found"))
			return
		}

		if err != nil {
			log.Info("Failed to get url")

			render.JSON(w, r, response.Error("Failed to get url"))
			return
		}

		log.Info("Got url", url)

		http.Redirect(w, r, url, http.StatusFound)
	}
}
