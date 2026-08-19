package main

import (
	"bufio"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	Domain   string
	Port     string
	Username string
	Token    string
	DBPath   string
}

type App struct {
	cfg   Config
	db    *sql.DB
	actor *Actor
}

func loadConfig() Config {
	return Config{
		Domain:   env0r("GLENMORE_DOMAIN", "localhost"),
		Port:     env0r("GLENMORE_PORT", "8080"),
		Username: env0r("GLENMORE_USERNAME", "rad"),
		Token:    os.Getenv("GLENMORE_TOKEN"),
		DBPath:   env0r("GLENMORE_DB", "glenmore.db"),
	}
}

func env0r(x, y string) string {
	if v := os.Getenv(x); v != "" {
		return v
	}
	return y
}

func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		if os.Getenv(k) == "" {
			os.Setenv(k, v)
		}
	}
}

func main() {
	loadDotEnv(".env")
	cfg := loadConfig()
	if cfg.Token == "" {
		log.Println("Warning: flow_token is empty")
	}
	db, err := opendb(cfg.DBPath)
	if err != nil {
		log.Fatalf("opendb: %v", err)
	}
	defer db.Close()

	actor, err := loca(db, cfg.Username)
	if err != nil {
		log.Fatalf("init actor: %v", err)
	}

	app := &App{cfg: cfg, db: db, actor: actor}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	log.Printf("glenmorepub is running on :%s (domain=%s actor=%s", cfg.Port,
		cfg.Domain, app.ActorURI())
	mux.HandleFunc("GET /.well-known/webfinger", app.hdlwf)
	mux.HandleFunc("GET /.well-known/host-meta", app.hdlhm)
	mux.HandleFunc("GET /users/{username}", app.hdla)
	mux.HandleFunc("GET /users/{username}/outbox", app.hdlo)
	mux.HandleFunc("GET /users/{username}/posts/{id}", app.hdln)
	mux.HandleFunc("GET /users/{username}/followers", app.hdlfollowers)
	mux.HandleFunc("GET /users/{username}/following", app.hdlfollowing)

	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
