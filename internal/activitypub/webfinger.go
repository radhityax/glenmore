package activitypub

import "fmt"

type WebFingerResponse struct {
	Subject string `json:"subject"`
	Aliases []string `json:"aliases"`
	Links []WebFingerLink `json:"links"`
}

type WebFingerLink struct {
	Rel string `json:"rel"`
	Type string `json:"type"`
	Href string `json:"href"`
}

func BuildWebFinger(domain, username string) *WebFingerResponse {
	actorIRI := fmt.Sprintf("https://%s/actor/%s", domain, username)
	return &WebFingerResponse{
		Subject: fmt.Sprintf("acct:%s@%s", username, domain),
		Aliases: []string{actorIRI},
		Links: []WebFingerLink{
			{
				Rel: "self",
				Type: "application/activity+json",
				Href: actorIRI,
			},
		},
	}
}
