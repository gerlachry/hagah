package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	// pq Database driver for Postgresql
	_ "github.com/lib/pq"
)

// Postgresql initialization script for scripture database.

func main() {
	log.Println("Starting db init")
	var schema = `
		CREATE TABLE users (
			id serial PRIMARY KEY,
			first_name text,
			last_name text,
			username text,
			email text,
			active boolean,
			create_dt timestamp,
			modified_dt timestamp
		);

		CREATE TABLE books (
			id serial PRIMARY KEY,
			name text,
			abbreviation text
		);

		CREATE TABLE scriptures (
			id serial PRIMARY KEY,
			book_id integer REFERENCES books (id),
			chapter integer,
			verse integer,
			text text
		);

		CREATE TABLE comments (
			id serial PRIMARY KEY,
			user_id integer REFERENCES users (id),
			scripture_id integer REFERENCES scriptures (id),
			comment text,
			created_dt timestamp,
			modified_dt timestamp,
			active boolean
		);	
	`
	db, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=scripture sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db.MustExec(schema)
	log.Println("Created scripture tables")
}
