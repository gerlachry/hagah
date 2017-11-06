package text

import (
	"context"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
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
		errors.Wrap(err, "Error initializing BulkProcessor")
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
	Connection *sqlx.DB
}

// NewPGHandler Function for creating a new PGHandler with a Postgres connection.
func NewPGHandler(url string) (pgHandler *PGHandler, err error) {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		errors.Wrap(err, "Error connecting to Postgres")
	}
	return &PGHandler{Connection: db}, nil
}

// Index Method for indexing a slice of Verses to Postgresql.
func (pgHandler *PGHandler) Index(verses []Verse) error {
	log.Println("Indexing verses to Postgres")

	return nil
}

// Close Method for closing down Postgresql resources.
func (pgHandler *PGHandler) Close() error {
	return nil
}
