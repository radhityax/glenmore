package main

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

const schemaSQL = `
    CREATE TABLE IF NOT EXISTS actor (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    username TEXT NOT NULL,
    display_name TEXT,
    bio TEXT,
    public_key_pem TEXT NOT NULL,
    private_key_pem TEXT NOT NULL,
    created_at TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS note (
    id TEXT PRIMARY KEY,
    content TEXT NOT NULL,
    published TEXT NOT NULL,
    in_reply_to TEXT,
    local INTEGER NOT NULL DEFAULT 1
    );

    CREATE TABLE IF NOT EXISTS activity (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    actor TEXT NOT NULL,
    object_uri TEXT,
    raw_json TEXT NOT NULL,
    published TEXT NOT NULL,
    direction TEXT NOT NULL
    );
    
    CREATE TABLE IF NOT EXISTS follower (
    actor_uri TEXT PRIMARY KEY,
    inbox TEXT NOT NULL,
    followed_at TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS following (
    actor_uri TEXT PRIMARY KEY,
    inbox TEXT NOT NULL,
    followed_at TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS liked (
    object_uri TEXT PRIMARY KEY,
    activity_uri TEXT NOT NULL,
    liked_at TEXT NOT NULL
    );
`

func opendb(path string) (*sql.DB, error) {
	dbfile := "file:" + path + "?_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)"
	db, err := sql.Open("sqlite", dbfile)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	if _, err := db.Exec(schemaSQL); err != nil {
		return nil, err
	}
	return db, nil
}
