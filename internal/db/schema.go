package db

const Schema = `
CREATE TABLE IF NOT EXISTS accounts (
	id	TEXT PRIMARY KEY,
	username TEXT NOT NULL UNIQUE,
	display_name TEXT,
	summary TEXT,
	domain TEXT NOT NULL DEFAULT '',
	inbox_url TEXT NOT NULL,
	outbox_url TEXT NOT NULL,
	following_url TEXT,
	followers_url TEXT,
	public_key TEXT NOT NULL,
	private_key TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS posts (
	id TEXT PRIMARY KEY,
	author_id TEXT REFERENCES accounts(id),
	content TEXT NOT NULL,
	content_warning TEXT DEFAULT '',
	visibility TEXT DEFAULT 'public',
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS follows (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	follower_id TEXT REFERENCES accounts(id),
	target_id TEXT REFERENCES accounts(id),
	approved INTEGER DEFAULT 0,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(follower_id, target_id)
);

CREATE TABLE IF NOT EXISTS inbox (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	account_id TEXT REFERENCES accounts(id),
	activity_id TEXT UNIQUE,
	activity_json TEXT NOT NULL,
	processed INTEGER DEFAULT 0,
	received_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`
