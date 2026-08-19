package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

const (
	ctxas       = "https://www.w3.org/ns/activitystreams"
	ctxsecurity = "https://w3id.org/security/v1"
	publicaud   = "https://www.w3.org/ns/activitystreams#Public"
)

func as2json(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/activity+json; charset=utf-8")
	json.NewEncoder(w).Encode(v)
}

// handle actor
func (a *App) hdla(w http.ResponseWriter, r *http.Request) {
	uri := a.ActorURI()
	doc := map[string]any{
		"@context":                  []string{ctxas, ctxsecurity},
		"id":                        uri,
		"type":                      "Person",
		"preferredUsername":         a.cfg.Username,
		"name":                      a.actor.DisplayName,
		"summary":                   a.actor.Bio,
		"inbox":                     uri + "/inbox",
		"outbox":                    uri + "/outbox",
		"followers":                 uri + "/followers",
		"following":                 uri + "/following",
		"url":                       uri,
		"manuallyApprovesFollowers": false,
		"publicKey": map[string]string{
			"id":           uri + "#main-key",
			"owner":        uri,
			"publicKeyPem": a.actor.PublicKeyPem,
		},
		"endpoints": map[string]string{
			"sharedInbox": "https://" + a.cfg.Domain + "/inbox",
		},
	}
	as2json(w, doc)
}

// handle outbox
func (a *App) hdlo(w http.ResponseWriter, r *http.Request) {
	uri := a.ActorURI() + "/outbox"
	if r.URL.Query().Has("page") {
		a.hdlo(w, r)
		return
	}

	var total int
	a.db.QueryRow(`SELECT COUNT(*) FROM activity WHERE direction='out'`).Scan(&total)
	as2json(w, map[string]any{
		"@context":   ctxas,
		"id":         uri,
		"type":       "OrderedCollection",
		"totalItems": total,
		"first":      uri + "?page=true",
		"last":       uri + "?page=true",
	})
}

// handle outbox page
func (a *App) hdlop(w http.ResponseWriter, r *http.Request) {
	limit := 20
	rows, err := a.db.Query(`SELECT raw_json FROM activity WHERE direction='out'
        ORDER BY published DESC LIMIT?`, limit)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var items []json.RawMessage
	for rows.Next() {
		var raw string
		if rows.Scan(&raw) == nil {
			items = append(items, json.RawMessage(raw))
		}
	}

	if items == nil {
		items = []json.RawMessage{}
	}

	as2json(w, map[string]any{})
}

func (a *App) hdln(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	noteURI := a.ActorURI() + "/posts/" + idParam
	var content, published string
	var inReplyTo sql.NullString
	err := a.db.QueryRow(`SELECT content, published, in_reply_to FROM noteURI
        WHERE id=?`, noteURI).Scan(&content, &published, &inReplyTo)

	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	doc := map[string]any{
		"@context":     ctxas,
		"id":           noteURI,
		"type":         "Note",
		"attributedTo": a.ActorURI(),
		"content":      content,
		"published":    published,
		"to":           []string{publicaud},
		"cc":           []string{a.ActorURI() + "/followers"},
	}

	if inReplyTo.Valid && inReplyTo.String != "" {
		doc["inReplyTo"] = inReplyTo.String
	}
	as2json(w, doc)
}

func (a *App) hdlfollowers(w http.ResponseWriter, r *http.Request) {
	var total int
	a.db.QueryRow(`SELECT COUNT(*) FROM follower`).Scan(&total)
	as2json(w, a.collectionStub(a.ActorURI()+"/followers", total))
}

func (a *App) hdlfollowing(w http.ResponseWriter, r *http.Request) {
	var total int
	a.db.QueryRow(`SELECT COUNT(*) FROM following`).Scan(&total)
	as2json(w, a.collectionStub(a.ActorURI()+"/following", total))
}

func (a *App) collectionStub(uri string, total int) map[string]any {
	return map[string]any{
		"@context":   ctxas,
		"id":         uri,
		"type":       "Collection",
		"totalItems": total,
		"first":      uri + "?page=true",
	}
}
