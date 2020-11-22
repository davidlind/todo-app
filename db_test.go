package main

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var postgresUser = "admin"
var postgresPassword = "admin"
var postgressDb = "webb"

func TestInsert(t *testing.T) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"), //wait.ForLog("database system is ready to accept connections"),
		Env:          map[string]string{"POSTGRES_USER": postgresUser, "POSTGRES_PASSWORD": postgresPassword, "POSTGRES_DB": postgressDb},
	}

	pgc, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		t.Error(err)
	}

	defer pgc.Terminate(context.Background())

	port, err := pgc.MappedPort(context.Background(), "5432")
	if err != nil {
		t.Error(err)
	}

	runMigrations("localhost", port.Int(), postgresUser, postgresPassword, postgressDb)

	count := 0
	expected := 2

	pg := newDB("localhost", port.Int(), "webb", "admin", "admin")
	pg.create("new todo")
	pg.create("another new todo")

	err = pg.db.QueryRow(context.Background(), "SELECT COUNT(*) FROM webb").Scan(&count)
	if err != nil {
		t.Error(err)
	}

	if count != expected {
		t.Errorf("Expected %d but got %d", expected, count)
	}
}

func TestDelete(t *testing.T) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"), //wait.ForLog("database system is ready to accept connections"),
		Env:          map[string]string{"POSTGRES_USER": postgresUser, "POSTGRES_PASSWORD": postgresPassword, "POSTGRES_DB": postgressDb},
	}

	pgc, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		t.Error(err)
	}

	defer pgc.Terminate(context.Background())

	port, err := pgc.MappedPort(context.Background(), "5432")
	if err != nil {
		t.Error(err)
	}

	runMigrations("localhost", port.Int(), postgresUser, postgresPassword, postgressDb)

	count := 0
	expected := 0

	pg := newDB("localhost", port.Int(), "webb", "admin", "admin")
	id := pg.create("new todo")
	pg.delete(id)

	err = pg.db.QueryRow(context.Background(), "SELECT COUNT(*) FROM webb").Scan(&count)
	if err != nil {
		t.Error(err)
	}

	if count != expected {
		t.Errorf("Expected %d but got %d", expected, count)
	}
}

func TestRead(t *testing.T) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"), //wait.ForLog("database system is ready to accept connections"),
		Env:          map[string]string{"POSTGRES_USER": postgresUser, "POSTGRES_PASSWORD": postgresPassword, "POSTGRES_DB": postgressDb},
	}

	pgc, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		t.Error(err)
	}

	defer pgc.Terminate(context.Background())

	port, err := pgc.MappedPort(context.Background(), "5432")
	if err != nil {
		t.Error(err)
	}

	runMigrations("localhost", port.Int(), postgresUser, postgresPassword, postgressDb)

	expected := "new todo"

	pg := newDB("localhost", port.Int(), "webb", "admin", "admin")
	id := pg.create("new todo")
	got, err := pg.read(id)
	if err != nil {
		t.Error(err)
	}

	if got.content != expected {
		t.Errorf("Expected \"%s\" but got \"%s\"", expected, got.content)
	}
}
