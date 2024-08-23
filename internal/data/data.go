package data

import (
	"classService/internal/biz"
	"classService/internal/conf"
	"classService/internal/errcode"
	"classService/internal/logPrinter"
	"context"
	"encoding/json"
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
const indexName = "class_info"

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewEsClient)

// Data .
type Data struct {
	cli *elastic.Client
	log logPrinter.LogerPrinter
}

// NewData .
func NewData(c *conf.Data, cli *elastic.Client, logPinter logPrinter.LogerPrinter, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		cli: cli,
		log: logPinter,
	}, cleanup, nil
}

func (d Data) AddClassInfo(ctx context.Context, classInfo biz.ClassInfo) error {
	// 创建文档
	_, err := d.cli.Index().
		Index(indexName).
		Id(classInfo.ID).
		BodyJson(classInfo).
		Do(ctx)
	if err != nil {
		d.log.FuncError(d.cli.Index().
			Index(indexName).
			BodyJson(classInfo).
			Do, err)
		return errcode.Err_EsAddClassInfo
	}
	return nil
}
func (d Data) SearchClassInfo(ctx context.Context, keyWords string) ([]biz.ClassInfo, error) {
	var classInfos = make([]biz.ClassInfo, 0)
	searchResult, err := d.cli.Search().
		Index(indexName). // 指定索引名称
		Query(elastic.NewBoolQuery().
			Should(
				elastic.NewMatchQuery("classname", keyWords),
				elastic.NewMatchQuery("teacher", keyWords),
			).
			MinimumShouldMatch("1"), // 至少匹配一个条件
		).Do(ctx) // 执行查询
	if err != nil {
		d.log.FuncError(d.cli.Search().
			Index(indexName).
			Query(elastic.NewBoolQuery().
				Should(
					elastic.NewMatchQuery("classname", keyWords),
					elastic.NewMatchQuery("teacher", keyWords),
				).
				MinimumShouldMatch("1"), // 至少匹配一个条件
			).Do, err)
		return nil, errcode.Err_EsSearchClassInfo
	}
	for _, hit := range searchResult.Hits.Hits {
		var classInfo biz.ClassInfo
		err := json.Unmarshal(hit.Source, &classInfo)
		if err != nil {
			d.log.FuncError(json.Unmarshal, err)
			continue
		}
		classInfos = append(classInfos, classInfo)
	}
	return classInfos, nil
}
