package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/fiatjaf/eventstore"
	"github.com/fiatjaf/eventstore/postgresql"
	"github.com/fiatjaf/relayer/v2"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/nbd-wtf/go-nostr"
)

type Config struct {
	POSTGRESQL_DATABASE string
	RELAYER_URL         string
}

var config Config

type Relay struct {
	//PostgresDatabase string `envconfig:"POSTGRESQL_DATABASE"`

	storage *postgresql.PostgresBackend
}

func (r *Relay) Name() string {
	return "BasicRelay"
}

func (r *Relay) Storage(ctx context.Context) eventstore.Store {
	return r.storage
}

func (r *Relay) Init() error {
	err := envconfig.Process("", &config)
	if err != nil {
		return fmt.Errorf("couldn't process envconfig: %w", err)
	}
	fmt.Printf("config: %s\n", config)
	// let config = Config{

	// }

	// every hour, delete all very old events
	// go func() {
	// 	db := r.Storage(context.TODO()).(*postgresql.PostgresBackend)

	// 	for {
	// 		time.Sleep(60 * time.Minute)
	// 		db.DB.Exec(`DELETE FROM event WHERE created_at < $1`, time.Now().AddDate(0, -3, 0).Unix()) // 3 months
	// 	}
	// }()

	return nil
}

func (r *Relay) AcceptEvent(ctx context.Context, evt *nostr.Event) bool {
	// block events that are too large
	jsonb, _ := json.Marshal(evt)
	if len(jsonb) > 10000 {
		return false
	}

	return true
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	r := Relay{}
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("failed to read from env: %v", err)
		return
	}
	fmt.Printf("Database URL: %s\n", config.POSTGRESQL_DATABASE)
	fmt.Printf("Relay URL: %s\n", config.RELAYER_URL)

	r.storage = &postgresql.PostgresBackend{DatabaseURL: config.POSTGRESQL_DATABASE}
	server, err := relayer.NewServer(&r)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}
	if err := server.Start(config.RELAYER_URL, 8008); err != nil {
		log.Fatalf("server terminated: %v", err)
	}

	// You can now use these variables in your relayer logic
	// Example: connectToRelayer(relayerURL, relayerSecret)
}
