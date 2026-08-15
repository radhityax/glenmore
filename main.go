package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
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
		Domain:   env0r("FLOW_DOMAIN", "localhost"),
		Port:     env0r("FLOW_PORT", "8080"),
		Username: env0r("FLOW_USERNAME", "rad"),
		Token:    os.Getenv("FLOW_TOKEN"),
		DBPath:   env0r("FLOW_DB", "flow.db"),
	}
}

func env0r(x, y string) string {
	if v := os.Getenv(x); v != "" {
		return v
	}
	return y
}

func main() {
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

	log.Printf("flowpub is running on :%s (domain=%s actor=%s", cfg.Port,
		cfg.Domain, app.ActorURI())
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
