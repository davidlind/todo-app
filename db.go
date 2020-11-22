package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type todo struct {
	id      int
	content string
}

type db interface {
	create()
	read()
	update()
	delete()
}

type pg struct {
	db *pgxpool.Pool
}

func (p *pg) create(content string) int {
	id := 0
	err := p.db.QueryRow(context.Background(), "insert into webb(content) values ($1) returning id", content).Scan(&id)
	if err != nil {
		log.Println(err)
	}

	return id
}

func (p *pg) read(id int) (todo, error) {
	t := todo{}
	err := p.db.QueryRow(context.Background(), "select * from webb where id=$1", id).Scan(&t.id, &t.content)
	if err != nil {
		return t, err
	}
	return t, nil
}

func (p *pg) update(id int, content string) error {
	_, err := p.db.Exec(context.Background(), "UPDATE webb SET content=$1 WHERE id=$2", content, id)
	if err != nil {
		return err
	}

	return nil
}

func (p *pg) delete(id int) error {
	_, err := p.db.Exec(context.Background(), "DELETE FROM webb WHERE id=$1", id)
	if err != nil {
		return err
	}

	return nil
}

func newDB(hostname string, port int, database string, username string, password string) *pg {
	log.Printf("connecting to %s:%d", hostname, port)
	conn, err := pgxpool.Connect(context.Background(), fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", username, password, hostname, port, database))
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	db := &pg{db: conn}

	return db
}
