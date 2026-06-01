package activitypub

const (
	ContextActivityStreams = "https://www.w3.org/ns/activitystreams"
	ContextW3IDSecurity = "https://w3id.org/security/v1"
)

type PublicKey struct {
	ID string `json:"id"`
	Owner string `json:"owner"`
	PublicKeyPem string `json:"publicKeyPem"`
}

type Actor struct {
	Context []string `json:"@context"`
	ID string `json:"id"`
	Type string `json:"type"`
	PreferredUsername string `json:"prefferedUsername"`
	Name string `json:"name,omitempty"`
	Summary string `json:"summary,omitempty"`
	Inbox string `json:"inbox"`
	Outbox string `json:"outbox"`
	Following string `json:"following,omitempty"`
	Followers string `json:"followers,omitempty"`
	PublicKey PublicKey `json:"publicKey"`
}

func NewPerson(id, username string) *Actor {
	return &Actor{
		Context: []string{
			ContextActivityStreams,
			ContextW3IDSecurity,
		},
		ID: id,
		Type: "Person",
		PreferredUsername: username,
		Inbox: id + "/inbox",
		Outbox: id + "/outbox",
		Following: id + "/following",
		Followers: id + "/followers",
	}
}

type Note struct {
	Context string `json:"@context"`
	ID string `json:"id"`
	Type string `json:"type"`
	AttributedTo string `json:"attributedTo"`
	Content string `json:"content"`
	Summary string `json:"summary,omitempty"`
	Published string `json:"published,omitempty"`
	To string `json:"to,omitempty"`
	CC string `json:"cc,omitempty"`
}

func NewNote(id, attributedTo, content string) *Note {
	return &Note {
		Context: ContextActivityStreams,
		ID: id,
		Type: "Note",
		AttributedTo: attributedTo,
		Content: content,
		To: "https://www.w3.org/ns/activitystreams#Public",
	}
}

type Activity struct {
	Context string `json:"@context"`
	ID string `json:"id"`
	Type string `json:"type"`
	Actor string `json:"actor"`
	Object interface{} `json:"object"`
	To string `json:"to,omitempty"`
	CC string `json:"cc,omitempty"`
}

func NewActivity(activityType, id, actor string, object interface{}) *Activity {
	return &Activity{
		Context: ContextActivityStreams,
		ID: id,
		Type: activityType,
		Actor: actor,
		Object: object,
		To: "https://www.w3.org/ns/activitystreams#Public",
	}
}

type Follow struct {
	Context string `json:"@context"`
	ID string `json:"id"`
	Type string `json:"type"`
	Actor string `json:"actor"`
	Object string `json:"object"`
}

func NewFollow(id, actor, target string) *Follow {
	return &Follow{
		Context: ContextActivityStreams,
		ID: id,
		Type: "Follow",
		Actor: actor,
		Object: target,
	}
}

type OrderedCollection struct {
	Context     string `json:"@context"`
	ID          string `json:"id"`
	Type        string `json:"type"`
	TotalItems  int    `json:"totalItems"`
	First       string `json:"first,omitempty"`
	Last        string `json:"last,omitempty"`
}

type OrderedCollectionPage struct {
	Context      string        `json:"@context"`
	ID           string        `json:"id"`
	Type         string        `json:"type"`
	PartOf       string        `json:"partOf"`
	OrderedItems []interface{} `json:"orderedItems"`
	Next         string        `json:"next,omitempty"`
	Prev         string        `json:"prev,omitempty"`
}
