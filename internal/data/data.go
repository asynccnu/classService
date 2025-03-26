package data

import (
	"context"
	"encoding/json"
	"github.com/asynccnu/classService/internal/errcode"
	clog "github.com/asynccnu/classService/internal/log"
	"github.com/asynccnu/classService/internal/model"
	"github.com/google/wire"
	"github.com/olivere/elastic/v7"
)

const mapping = `{
  "settings": {
    "number_of_shards": 3,
    "number_of_replicas": 1,
    "analysis": {
      "analyzer": {
        "edge_ngram_analyzer": {
          "tokenizer": "edge_ngram_tk",
          "filter": ["lowercase"]
        }
      },
      "tokenizer": {
        "edge_ngram_tk": {
          "type": "edge_ngram",
          "min_gram": 1,
          "max_gram": 25,
          "token_chars": ["letter", "digit"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id": {
        "type": "keyword"
      },
      "day": {
        "type": "integer"
      },
      "teacher": {
        "type": "text",
        "analyzer": "edge_ngram_analyzer"
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
        "type": "text",
        "analyzer": "edge_ngram_analyzer",
        "fields": {
          "keyword": { "type": "keyword" }
        }
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
}

// NewData .
func NewData(cli *elastic.Client) (*Data, func(), error) {
	cleanup := func() {
		clog.LogPrinter.Info("closing the data resources")
	}
	return &Data{
		cli: cli,
	}, cleanup, nil
}

func (d Data) AddClassInfo(ctx context.Context, classInfo model.ClassInfo) error {
	// 创建文档
	_, err := d.cli.Index().
		Index(indexName).
		Id(classInfo.ID).
		BodyJson(classInfo).
		Do(ctx)
	if err != nil {
		clog.LogPrinter.Errorf("es: failed to add class_info[%+v]: %v", classInfo, err)
		return errcode.Err_EsAddClassInfo
	}
	return nil
}

// 删除除了 year=xnm 和 semester=xqm 之外的所有数据
func (d Data) RemoveClassInfo(ctx context.Context, xnm, xqm string) {
	// 创建查询条件，删除除了 year=xnm 和 semester=xqm 之外的所有数据
	query := elastic.NewBoolQuery().
		Should(
			elastic.NewBoolQuery().MustNot(elastic.NewTermQuery("year", xnm)),     // year != xnm
			elastic.NewBoolQuery().MustNot(elastic.NewTermQuery("semester", xqm)), // semester != xqm
		)

	// 执行删除操作
	deleteResponse, err := d.cli.DeleteByQuery().
		Index(indexName).
		Query(query).
		Slices("auto"). // 自动计算分片数
		Size(1000).     // 每批次删除 1000 条
		Do(ctx)
	if err != nil {
		clog.LogPrinter.Errorf("es: failed to delete class_info[except (xnm:%v,xqm:%v)]: %v", xnm, xqm, err)
		return
	}
	clog.LogPrinter.Infof("Deleted %d documents", deleteResponse.Deleted)
}

func (d Data) SearchClassInfo(ctx context.Context, keyWords string, xnm, xqm string) ([]model.ClassInfo, error) {
	var classInfos = make([]model.ClassInfo, 0)
	searchResult, err := d.cli.Search().
		Index(indexName). // 指定索引名称
		Query(
			elastic.NewBoolQuery().
				Should(
					elastic.NewMatchPhrasePrefixQuery("classname", keyWords).Boost(2),
					elastic.NewMatchPhrasePrefixQuery("teacher", keyWords).Boost(1.5),
				).
				MinimumShouldMatch("1"). // 至少匹配一个条件
				Filter(
					elastic.NewTermQuery("year", xnm),     // year 精确匹配 xnm
					elastic.NewTermQuery("semester", xqm), // semester 精确匹配 xqm
				),
		).Do(ctx)

	if err != nil {
		clog.LogPrinter.Errorf("es: failed to search class_info[keywords:%v xnm:%v xqm:%v]: %v", keyWords, xnm, xqm, err)
		return nil, errcode.Err_EsSearchClassInfo
	}
	for _, hit := range searchResult.Hits.Hits {
		var classInfo model.ClassInfo
		err := json.Unmarshal(hit.Source, &classInfo)
		if err != nil {
			clog.LogPrinter.Errorf("json unmarshal %v failed: %v", hit.Source, err)
			continue
		}
		classInfos = append(classInfos, classInfo)
	}
	return classInfos, nil
}
