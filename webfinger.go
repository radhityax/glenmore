package main

import (
	"encoding/json"
	"net/http"
)

type jrdLink struct {
	Rel  string `json:"rel"`
	Type string `json:"type,omitempty"`
	Href string `json:"href,omitempty"`
}

type jrd struct {
	Subject string    `json:"subject"`
	Aliases []string  `json:"aliases,omitempty"`
	Links   []jrdLink `json:"links"`
}

// handle web finger
func (a *App) hdlwf(w http.ResponseWriter, r *http.Request) {
	resource := r.URL.Query().Get("resource")
	if resource == "" {
		http.Error(w, "missing resource", http.StatusBadRequest)
		return
	}

	actorURI := a.ActorURI()
	expectedAcct := "acct:" + a.cfg.Username + "@" + a.cfg.Domain

	match := resource == expectedAcct || resource == actorURI
	if !match {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	resp := jrd{
		Subject: expectedAcct,
		Aliases: []string{actorURI},
		Links: []jrdLink{
			{
				Rel:  "self",
				Type: "application/activity+json",
				Href: actorURI,
			},
		},
	}

	w.Header().Set("Content-Type", "application/jrd+json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}

// handle host meta
func (a *App) hdlhm(w http.ResponseWriter, r *http.Request) {
	tmpl := "https://" + a.cfg.Domain + "/.well-known/webfinger?resource={uri}"
	body := `<?xml version="1.0" encoding="UTF-8"?>
    <XRD xmlns="https://docs.oasis-open.org/ns/xri/xrd-1.0">
    <Link rel="lrdd" template="` + tmpl + `" type="application/jrd+json"/>
    </XRD>`

	w.Header().Set("Content-Type", "application/xrd+xml; charset=utf-8")
	w.Write([]byte(body))
}
