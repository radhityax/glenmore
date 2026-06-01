package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"fmt"

	"github.com/go-chi/chi/v5"
	"glenmore/internal/activitypub"
	"glenmore/internal/db"
)

type Handler struct {
	DB *sql.DB
	Cfg *Config
}

type Config struct {
	Host string
	Port int
}

/* Webfinger: GET /.well-known/webfinger?resource=acct:user@domain */
func (h *Handler) WebFinger(w http.ResponseWriter, r *http.Request) {
	resource := r.URL.Query().Get("resource")
	if resource == "" {
		http.Error(w, "missing resource", http.StatusBadRequest)
	}

	var username, domain string
	if n, _ := fmt.Scanf(resource, "acct:%s@%s", &username, &domain); n != 2 {
		http.Error(w, "invalid resource format", http.StatusBadRequest)
		return
	}

	acc, err := db.GetAccountByUsername(h.DB, username)
	if err != nil {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	wf := activitypub.BuildWebFinger(domain, acc.Username)
	w.Header().Set("Content-Type", "application/jrd+json")
	json.NewEncoder(w).Encode(wf)
}

/* Actor: GET /actor/{username} */

func (h *Handler) Actor(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	acc, err := db.GetAccountByUsername(h.DB, username)
	if err != nil {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	actor := activitypub.NewPerson(acc.ID, acc.Username)
	actor.Name = acc.DisplayName
	actor.Summary = acc.Summary
	actor.PublicKey = activitypub.PublicKey {
		ID: acc.ID + "#main-key",
		Owner: acc.ID,
		PublicKeyPem: acc.PublicKey,
	}

	w.Header().Set("Content-Type", "application/activity+json")
	json.NewEncoder(w).Encode(actor)
}

/* Inbox: POST /actor/{username}/inbox */
func (h *Handler) Inbox(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

/* Outbox: GET /actor/{username}/outbox */
func (h *Handler) Outbox(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/activity+json")
	json.NewEncoder(w).Encode(activitypub.OrderedCollection{
		Context: activitypub.ContextActivityStreams,
		ID: r.URL.String(),
		Type: "OrderedCollection",
		TotalItems: 0,
	})
}

/* Followers: GET /actor/{username}/followers */
func (h *Handler) Followers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/activity+json")
	json.NewEncoder(w).Encode(activitypub.OrderedCollection{
		Context: activitypub.ContextActivityStreams,
		ID: r.URL.String(),
		Type: "OrderedCollection",
		TotalItems: 0,
	})
}


func (h *Handler) Following(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/activity+json")
	json.NewEncoder(w).Encode(activitypub.OrderedCollection{
		Context: activitypub.ContextActivityStreams,
		ID: r.URL.String(),
		Type: "OrderedCollection",
		TotalItems: 0,
	})
}

/* Register: POST /api/register */
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

/* Login: POST /api/login */
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

/* Timeline: GET /api/timeline */
func (h *Handler) Timeline(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
