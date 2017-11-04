// Main entry point for indexing scripture.
// Currently only supports the Digital Bible Platform http://www.digitalbibleplatform.com
// as a source for scripture.  Other sources such as a local epub document are planned.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gerlachr/scripture/internal/text"
	"github.com/pkg/errors"
	"gopkg.in/olivere/elastic.v5"
)

const (
	dbpBaseURL = "http://dbt.io"
)

// esURL the url to use for connecting to Elasticsearch if using for the database
var esURL string

// esIndex the Elasticsearch index to use when indexing scripture
var esIndex string

// Postgresql connection string URL if using Postgres as the database.
var pgURL string

func init() {
	flag.StringVar(&esURL, "esURL", "http://localhost:9200", "Elasticsearch host and port.  Defaults to localhost:9200")
	flag.StringVar(&esIndex, "esIndex", "scripture", "Name of the Elastisearch index to use")
	flag.StringVar(&pgURL, "pgURL", "localhost", "Postgresql connection URL string")

	_, e := os.LookupEnv("ES_URL")
	if e {
		_, b := os.LookupEnv("ES_USER")
		if !b {
			panic("Missing required ES_USER environment varialbe")
		}
		_, b = os.LookupEnv("ES_PWD")
		if !b {
			panic("Missing required ES_PWD environment varialbe")
		}
	}

	_, p := os.LookupEnv("PG_URL")
	if !p && !e {
		panic("ES_URL or PG_URL required")
	}

	_, b := os.LookupEnv("DBP_KEY")
	if !b {
		panic("Missing required DBP_KEY (Digital Bible Platform API Key) environment variable")
	}

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.Lmicroseconds)
}

func main() {
	log.Printf("Starting indexing of Scripture to host: %s; index: %s", esURL, esIndex)
	//ctx := context.Background()
	client, err := elastic.NewClient(
		elastic.SetURL(esURL),
		elastic.SetBasicAuth(os.Getenv("ES_USER"), os.Getenv("ES_PWD")),
		elastic.SetSniff(false))
	if err != nil {
		// Handle error
		panic(err)
	}
	defer client.Stop()

	esversion, err := client.ElasticsearchVersion(esURL)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)
	handler := text.ESHandler{Client: client, ESIndex: esIndex}

	booksURL := dbpBaseURL + "/library/book?v=2&dam_id=ENGESVO2ET&key=" + os.Getenv("DBP_KEY")
	resp, err := http.Get(booksURL)
	if err != nil {
		errors.Wrap(err, "Error fetching booksURL")
	}

	books := make([]text.Book, 0)
	err = json.NewDecoder(resp.Body).Decode(&books)
	if err != nil {
		errors.Wrap(err, "Error decoding Books json response")
	}

	for _, b := range books {
		log.Printf("indexing book: %s", b.BookID)
		textURL := dbpBaseURL + "/text/verse/?v=2&dam_id=ENGESVO2ET&key=" + os.Getenv("DBP_KEY") + "&book_id=" + b.BookID
		log.Printf("verse url: %s", textURL)
		resp, err = http.Get(textURL)
		if err != nil {
			errors.Wrap(err, fmt.Sprintf("Error fetching verses with url: %s", textURL))
		}
		log.Printf(resp.Status)
		verses := make([]text.Verse, 0)
		err = json.NewDecoder(resp.Body).Decode(&verses)
		if err != nil {
			errors.Wrap(err, "Error marshelling verses")
		}
		log.Printf("Verse count: %d", len(verses))
		handler.Index(verses)
	}

}
