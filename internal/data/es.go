package data

import (
	"context"
	"fmt"
	"github.com/asynccnu/classService/internal/conf"
	clog "github.com/asynccnu/classService/internal/log"
	"github.com/olivere/elastic/v7"
)

func NewEsClient(c *conf.Data) (*elastic.Client, error) {
	ctx := context.Background()

	// 配置 Elasticsearch 的 URL 和嗅探选项
	urlOpt := elastic.SetURL(c.Es.Url)
	sniffOpt := elastic.SetSniff(c.Es.Setsniff)

	// 配置基本认证，使用用户名和密码
	authOpt := elastic.SetBasicAuth(c.Es.Username, c.Es.Password)

	// 创建 Elasticsearch 客户端
	cli, err := elastic.NewClient(urlOpt, sniffOpt, authOpt)
	if err != nil {
		panic(fmt.Sprintf("es connect fail: %v", err))
	}

	clog.LogPrinter.Info("connect to elasticsearch successfully")

	// 检查索引是否存在
	exist, err := cli.IndexExists(indexName).Do(ctx)
	if err != nil {
		panic(err)
	}

	// 如果索引存在，先删除索引
	if exist {
		deleteIndex, err := cli.DeleteIndex(indexName).Do(ctx)
		if err != nil {
			panic(fmt.Sprintf("failed to delete existing index: %v", err))
		}
		if !deleteIndex.Acknowledged {
			panic("delete index failed")
		}
		clog.LogPrinter.Info("Existing index deleted successfully")
	}

	// 创建新的索引
	createIndex, err := cli.CreateIndex(indexName).BodyString(mapping).Do(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to create index: %v", err))
	}
	if !createIndex.Acknowledged {
		panic("create index failed")
	}
	clog.LogPrinter.Info("Es create index successfully")

	return cli, nil
}
