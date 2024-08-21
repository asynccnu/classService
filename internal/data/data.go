package data

import (
	"classService/internal/conf"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/olivere/elastic/v7"
)

const mapping = `
{
  "mappings": {
    "properties": {
      "id": {
        "type": "keyword"
      },
      "created_at": {
        "type": "date",
        "format": "yyyy-MM-dd'T'HH:mm:ss.SSSZ||yyyy-MM-dd'T'HH:mm:ssZ||epoch_millis"
      },
      "updated_at": {
        "type": "date",
        "format": "yyyy-MM-dd'T'HH:mm:ss.SSSZ||yyyy-MM-dd'T'HH:mm:ssZ||epoch_millis"
      },
      "day": {
        "type": "integer"
      },
      "teacher": {
        "type": "text"
      },
      "where": {
        "type": "text"
      },
      "class_when": {
        "type": "text"
      },
      "week_duration": {
        "type": "text"
      },
      "classname": {
        "type": "text"
      },
      "credit": {
        "type": "float"
      },
      "weeks": {
        "type": "integer"
      },
      "semester": {
        "type": "keyword"
      },
      "year": {
        "type": "keyword"
      }
    }
  }
}`

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData)

// Data .
type Data struct {
	// TODO wrapped database client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{}, cleanup, nil
}
func NewEsClient(c *conf.Data, logger log.Logger) (*elastic.Client, error) {
	ctx := context.Background()
	urlOpt := elastic.SetURL(c.Es.Url)
	sniffOpt := elastic.SetSniff(c.Es.Setsniff)
	cli, err := elastic.NewClient(urlOpt, sniffOpt)
	if err != nil {
		panic(fmt.Sprintf("es connect fail: %v", err))
	}
	exist, err := cli.IndexExists(indexName).Do(ctx)
	if err != nil {
		panic(err)
	}
	if !exist {
		createIndex, err := cli.CreateIndex(indexName).BodyString(mapping).Do(ctx)
		if err != nil {
			panic(err)
		}
		if !createIndex.Acknowledged {
			panic("create index failed")
		}
	}
	log.Info("Es create index successfully")
	return cli, nil
}
