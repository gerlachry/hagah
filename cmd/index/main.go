// Main entry point for indexing scripture.
// Currently only supports the Digital Bible Platform http://www.digitalbibleplatform.com
// as a source for scripture.  Other sources such as a local epub document are planned.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gerlachr/scripture/internal/text"
	"github.com/pkg/errors"
)

const (
	dbpBaseURL = "http://dbt.io"
)

func init() {
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

func index(handler text.IndexHandler, books []text.Book, testemant string) error {
	for _, b := range books {
		log.Printf("indexing book: %s", b.BookID)
		textURL := dbpBaseURL + "/text/verse/?v=2&dam_id=ENGESV" + testemant + "2ET&key=" + os.Getenv("DBP_KEY") + "&book_id=" + b.BookID
		log.Printf("verse url: %s", textURL)
		resp, err := http.Get(textURL)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Error fetching verses with url: %s", textURL))
		}

		log.Printf(resp.Status)
		verses := make([]text.Verse, 0)
		err = json.NewDecoder(resp.Body).Decode(&verses)
		if err != nil {
			return errors.Wrap(err, "Error marshelling verses")
		}
		log.Printf("Verse count: %d", len(verses))
		err = handler.Index(verses)
		if err != nil {
			return errors.Wrap(err, "Fatal error indexing verses")
		}
	}
	return nil
}

func main() {
	log.Println("Starting indexing of Scripture")
	var handler text.IndexHandler

	esURL, e := os.LookupEnv("ES_URL")
	esIndex, _ := os.LookupEnv("ES_INDEX")
	pgURL, p := os.LookupEnv("PG_URL")
	var err error
	if e {
		handler = text.NewESHandler(esURL, esIndex)
	} else if p {
		handler, err = text.NewPGHandler(pgURL)
		if err != nil {
			log.Fatalf("error while obtaining pg handler : %s", err)
		}
	}

	booksURL := dbpBaseURL + "/library/book?v=2&dam_id=ENGESVO2ET&key=" + os.Getenv("DBP_KEY")
	log.Println(booksURL)
	resp, err := http.Get(booksURL)
	if err != nil {
		panic("Error fetching booksURL")
	}

	booksOT := make([]text.Book, 0)
	err = json.NewDecoder(resp.Body).Decode(&booksOT)
	if err != nil {
		log.Println("Error decoding Books json response")
		panic(err)
	}

	log.Printf("OT Books count: %v", len(booksOT))
	err = index(handler, booksOT, "O")
	if err != nil {
		log.Println("error idexing")
		panic(err)
	}

	booksURL = dbpBaseURL + "/library/book?v=2&dam_id=ENGESVN2ET&key=" + os.Getenv("DBP_KEY")
	log.Println(booksURL)
	resp, err = http.Get(booksURL)
	if err != nil {
		errors.Wrap(err, "Error fetching booksURL")
	}

	booksNT := make([]text.Book, 0)
	err = json.NewDecoder(resp.Body).Decode(&booksNT)
	if err != nil {
		errors.Wrap(err, "Error decoding Books json response")
	}

	log.Printf("NT Books count: %v", len(booksNT))
	index(handler, booksNT, "N")
}
