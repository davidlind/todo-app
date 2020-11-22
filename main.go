package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"
)

func main() {
	db := newDB("localhost", 5432, "webb", "admin", "admin")
	db.create("#1")
	db.create("#2")
	t, err := db.read(3)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t.content)
}

func runMigrations(hostname string, port int, username string, password string, database string) {
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", username, password, hostname, port, database))
	if err != nil {
		log.Fatal(err)
	}

	migrator, err := migrate.NewMigrator(context.Background(), conn, "version")
	if err != nil {
		log.Fatal(err)
	}

	err = migrator.LoadMigrations("migrations/")
	if err != nil {
		log.Fatalf("failed loading migrations: %v", err)
	}

	err = migrator.Migrate(context.Background())
	if err != nil {
		log.Fatalf("failed running migrations: %v", err)
	}
}
