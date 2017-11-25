package main

import (
	"database/sql"
	"log"

	// pq Database driver for Postgresql
	_ "github.com/lib/pq"
)

// Postgresql initialization script for scripture database.

// post indexing step to add index for full text search
// update scriptures set tsv = setweight(to_tsvector(text), 'A');

// Query for full text searching, top 10 hits
// select b.name, chapter, verse, text, ts_rank(tsv, q) as rank, q
// from scriptures s, books b, plainto_tsquery('the word') q
// where s.book_id = b.id and tsv @@ q
// order by rank desc
// limit 10;

// Query for adding highlighting to full text search
// select name as book, chapter, verse, ts_headline(text, q)
// from (
//  select b.name, chapter, verse, text, ts_rank(tsv, q) as rank, q
// from scriptures s, books b, plainto_tsquery('god the father') q
// where s.book_id = b.id and tsv @@ q
// order by rank desc
// limit 10) a
// order by rank desc;

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
			name text NOT NULL UNIQUE,
			abbreviation text NOT NULL UNIQUE
		);

		CREATE TABLE scriptures (
			id serial PRIMARY KEY,
			book_id integer REFERENCES books (id),
			chapter integer NOT NULL,
			verse integer NOT NULL,
			text text NOT NULL,
			tsv tsvector,
			CONSTRAINT scriptures_unq UNIQUE (book_id, chapter, verse)
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
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=scripture sslmode=disable")
	if err != nil {
		log.Fatalf("error connecting to database: %s", err)
	}
	res, err := db.Exec(schema)
	if err != nil {
		log.Fatalf("error initializing database: ", err)
	}
	log.Printf("Created scripture tables: %s", res)
}
