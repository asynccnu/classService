package data

import (
	"context"
	"fmt"
	"github.com/asynccnu/classService/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/olivere/elastic/v7"
)

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
