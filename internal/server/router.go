package server

import (
	"database/sql"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(db *sql.DB, cfg Config) http.Handler {
	h := &Handler{DB: db, Cfg: &cfg}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(Logging)
	r.Use(middleware.Recoverer)
	r.Use(ActivityPubContentType)
	rl := NewRateLimiter(10, 20)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("glenmore is running"))
	})

	r.Group(func(r chi.Router) {
		r.Use(RateLimit(rl))
		r.Get("/.well-known/webfinger", h.WebFinger)
		r.Get("/actor/{username}", h.Actor)
		r.Post("/actor/{username}/inbox", h.Inbox)
		r.Get("/actor/{username}/outbox", h.Outbox)
		r.Get("/actor/{username}/followers", h.Followers)
		r.Get("/actor/{username}/following", h.Following)
	})

	r.Group(func(r chi.Router) {
		r.Post("/api/register", h.Register)
		r.Post("/api/login", h.Login)
		r.Get("/api/timeline", h.Timeline)
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("glenmore — minimal ActivityPub server"))
	})

	return r
}
