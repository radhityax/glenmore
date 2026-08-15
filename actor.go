package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
	"fmt"
	"time"
)

type Actor struct {
	Username      string
	DisplayName   string
	Bio           string
	PublicKeyPem  string
	PrivateKeyPem string
	PrivateKey    *rsa.PrivateKey
}

func (a *App) ActorURI() string {
	return fmt.Sprintf("https://%s/users/%s", a.cfg.Domain, a.cfg.Username)
}

// load or create actor: grab actor id=1 from db; if it is not exist, generate
// new keypair
func loca(db *sql.DB, username string) (*Actor, error) {
	var a Actor
	err := db.QueryRow(`SELECT username, display_name, bio, public_key_pem,
        private_key_pem FROM actor WHERE id=1`).Scan(&a.Username,
		&a.DisplayName, &a.Bio, &a.PublicKeyPem, &a.PrivateKeyPem)
	if err == nil {
		block, _ := pem.Decode([]byte(a.PrivateKeyPem))
		if block == nil {
			return nil, errors.New("private key PEM is not valid on db")
		}
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			k2, err2 := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err2 != nil {
				return nil, fmt.Errorf("parse private key: %w", err2)
			}
			key = k2.(*rsa.PrivateKey)
		}
		a.PrivateKey = key
		return &a, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	privPEM := string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}))

	pubBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return nil, err
	}

	pubPEM := string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}))

	a = Actor{
		Username:      username,
		PublicKeyPem:  pubPEM,
		PrivateKeyPem: privPEM,
		PrivateKey:    key,
	}

	_, err = db.Exec(`INSERT INTO actor(id, username, display_name, bio,
        public_key_pem, private_key_pem, created_at)
        VALUES (1, ?, '', '', ?, ?, ?)`, username, pubPEM, privPEM,
		time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	return &a, nil
}
