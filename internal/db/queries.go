package db

import (
	"database/sql"
	"time"
)

type Account struct {
	ID string `json:"id"`
	Username string `json:"username"`
	DisplayName string `json:"display_name,omitempty"`
	Summary string `json:"summary, omitempty"`
	Domain string `json:"domain"`
	InboxURL string `json:"inbox_url"`
	OutboxURL string `json:"outbox_url"`
	FollowingURL string `json:"following_url,omitempty"`
	FollowersURL string `json:"followers_url,omitempty"`
	PublicKey string `json:"public_key"`
	PrivateKey string `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateAccount(db *sql.DB, acc *Account) error {
	_, err := db.Exec(`
	INSERT INTO accounts (id, username, display_name, summary, domain,
	inbox_url, outbox_url, following_url, followers_url, public_key,
	private_key)
	VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
	acc.ID, acc.Username, acc.DisplayName, acc.Summary, acc.Domain,
	acc.InboxURL, acc.OutboxURL, acc.FollowingURL, acc.FollowersURL,
	acc.PublicKey, acc.PrivateKey,
)
return err
}

func GetAccountByUsername(db *sql.DB, username string) (*Account, error) {
	acc := &Account{}
	err := db.QueryRow(`
	SELECT id, username, display_name, summary, domain,
	inbox_url, outbox_url, following_url, followers_url,
	public_key, private_key, created_at FROM accounts WHERE username = ?`,
	username).Scan(&acc.ID, &acc.Username, &acc.DisplayName, &acc.Summary,
	&acc.Domain, &acc.InboxURL, &acc.OutboxURL, &acc.FollowingURL, &acc.FollowersURL,
	&acc.PublicKey, &acc.PrivateKey, &acc.CreatedAt)
	if err != nil {
		return nil, err
	}
	return acc, nil
}
func GetAccountByIRI(db *sql.DB, iri string) (*Account, error) {
	acc := &Account{}
	err := db.QueryRow(`
	SELECT id, username, display_name, summary, domain,
	inbox_url, outbox_url, following_url, followers_url,
	public_key, private_key, created_at
	FROM accounts WHERE id = ?`, iri).
	Scan(&acc.ID, &acc.Username, &acc.DisplayName, &acc.Summary, &acc.Domain,
	&acc.InboxURL, &acc.OutboxURL, &acc.FollowingURL, &acc.FollowersURL,
	&acc.PublicKey, &acc.PrivateKey, &acc.CreatedAt)
	if err != nil {
		return nil, err
	}
	return acc, nil
}
