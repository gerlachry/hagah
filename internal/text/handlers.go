package text

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/olivere/elastic.v5"
	// pq Database driver for Postgresql
	_ "github.com/lib/pq"
)

// IndexHandler interface for indexing documents
type IndexHandler interface {
	//Init(url string) error
	Index(verses []Verse) error
	Close() error
}

// ESHandler struct for handling Elasticsearch related activities
type ESHandler struct {
	Client  *elastic.Client
	ESIndex string
}

// NewESHandler Method for creating a new ESHandler
func NewESHandler(url string, index string) *ESHandler {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetBasicAuth(os.Getenv("ES_USER"), os.Getenv("ES_PWD")),
		elastic.SetSniff(false))
	if err != nil {
		// Handle error
		panic(err)
	}

	return &ESHandler{Client: client, ESIndex: index}

}

// Index Method for indexing a slice of Verses to Elasticsearch.
func (esHandler *ESHandler) Index(verses []Verse) error {
	ctx := context.Background()
	serivce, err := esHandler.Client.BulkProcessor().Name("ScriptureProcessor").Workers(2).BulkActions(1000).Do(ctx)
	if err != nil {
		return errors.Wrap(err, "Error initializing BulkProcessor")
	}
	defer serivce.Close()

	for _, v := range verses {
		id := v.GetID()
		r := elastic.NewBulkIndexRequest().Index(esHandler.ESIndex).Type("Verse").Id(id).Doc(v)
		serivce.Add(r)
	}
	return nil
}

// Close Method for closing down the Elasticsearch resources.
func (esHandler *ESHandler) Close() error {
	log.Println("Stopping Elasticsearch client")
	esHandler.Client.Stop()
	return nil
}

// PGHandler handler for writing scriptures to Postgresql.
type PGHandler struct {
	Connection *sql.DB
	Books      map[string]int
}

// NewPGHandler Function for creating a new PGHandler with a Postgres connection.
func NewPGHandler(url string) (*PGHandler, error) {
	log.Printf("Postgres url: %s", url)
	db, err := sql.Open("postgres", url)
	if err != nil {
		panic("Error connecting to Postgres")
	}
	err = db.Ping()
	if err != nil {
		errors.Wrap(err, "Error pinging postgres db")
	} else {
		log.Println("DB all good")
	}

	// Populate list of books and their IDs for later lookup
	var book Books
	rows, err := db.Query("select id, name, abbreviation from books order by id asc")
	if err != nil {
		return nil, errors.Wrap(err, "error selecting books")
	}
	defer rows.Close()

	bookMap := make(map[string]int)
	for rows.Next() {
		err = rows.Scan(&book.ID, &book.Name, &book.Abbreviation)
		if err != nil {
			log.Printf("Error fetching existing books: %s", err)
			return nil, errors.Wrap(err, "error fetching book results")
		}
		bookMap[book.Abbreviation] = book.ID
	}
	return &PGHandler{Connection: db, Books: bookMap}, nil
}

// Index Method for indexing a slice of Verses to Postgresql.
func (pgHandler *PGHandler) Index(verses []Verse) error {
	log.Println("Indexing verses to Postgres")
	log.Printf("Connetion: %v", pgHandler.Connection)

	tx, err := pgHandler.Connection.Begin()
	if err != nil {
		return errors.Wrap(err, "error obtaining a transaction lock")
	}

	for _, v := range verses {
		// Insert book if needed
		bookID, ok := pgHandler.Books[v.BookID]
		if !ok {
			log.Printf("Creating book: %s", v.BookID)
			err := tx.QueryRow("Insert into books (name, abbreviation) values ($1, $2) returning id", v.BookName, v.BookID).Scan(&bookID)
			if err != nil {
				log.Printf("error creating book: %s", err)
				return errors.Wrap(err, "error creating book")
			}
			pgHandler.Books[v.BookID] = bookID
		}

		// Insert verse
		stmt, err := tx.Prepare("Insert into scriptures (book_id, chapter, verse, text) values ($1, $2, $3, $4) ON CONFLICT (book_id, chapter, verse) DO NOTHING")
		_, err = stmt.Exec(bookID, v.ChapterID, v.Verse, v.VerseText)
		if err != nil {
			log.Printf("error upserting verse: %s; err: %s", v, err)
			return errors.Wrap(err, "error upserting verse")
		}
	}
	tx.Commit()
	log.Printf("Inserted %d records", len(verses))
	return nil
}

// Close Method for closing down Postgresql resources.
func (pgHandler *PGHandler) Close() error {
	return nil
}
