package text

import (
	"context"

	"github.com/pkg/errors"
	"gopkg.in/olivere/elastic.v5"
)

// ESHandler struct for handling Elasticsearch related activities
type ESHandler struct {
	Client  *elastic.Client
	ESIndex string
}

// Index Method for indexing a slice of Verses to Elasticsearch
func (eshandler *ESHandler) Index(verses []Verse) error {
	ctx := context.Background()
	serivce, err := eshandler.Client.BulkProcessor().Name("ScriptureProcessor").Workers(2).BulkActions(1000).Do(ctx)
	if err != nil {
		errors.Wrap(err, "Error initializing BulkProcessor")
	}
	defer serivce.Close()

	for _, v := range verses {
		id := v.GetID()
		r := elastic.NewBulkIndexRequest().Index(eshandler.ESIndex).Type("Verse").Id(id).Doc(v)
		serivce.Add(r)
	}
	return nil
}
