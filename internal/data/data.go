package data

import (
	"context"
	"encoding/json"
	"github.com/asynccnu/classService/internal/biz"
	"github.com/asynccnu/classService/internal/conf"
	"github.com/asynccnu/classService/internal/errcode"
	clog "github.com/asynccnu/classService/internal/log"
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
}

// NewData .
func NewData(c *conf.Data, cli *elastic.Client, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		cli: cli,
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
		clog.LogPrinter.Errorf("es: failed to add class_info[%+v]: %v", classInfo, err)
		return errcode.Err_EsAddClassInfo
	}
	return nil
}
func (d Data) RemoveClassInfo(ctx context.Context, xnm, xqm string) {
	// 创建查询条件，删除除了 year=xnm 和 semester=xqm 之外的所有数据
	query := elastic.NewBoolQuery().
		MustNot( // 这里使用 MustNot 来排除符合条件的数据
			elastic.NewTermQuery("year", xnm),     // 排除 year 字段值为 xnm 的数据
			elastic.NewTermQuery("semester", xqm), // 排除 semester 字段值为 xqm 的数据
		)

	// 执行删除操作
	deleteResponse, err := d.cli.DeleteByQuery().
		Index(indexName). // 指定索引名称
		Query(query).     // 传递查询条件
		Do(ctx)           // 执行删除操作
	if err != nil {
		clog.LogPrinter.Errorf("es: failed to delete class_info[xnm:%v,xqm:%v]: %v", xnm, xqm, err)
		return
	}
	clog.LogPrinter.Infof("Deleted %d documents", deleteResponse.Deleted)
}

func (d Data) SearchClassInfo(ctx context.Context, keyWords string, xnm, xqm string) ([]biz.ClassInfo, error) {
	var classInfos = make([]biz.ClassInfo, 0)
	searchResult, err := d.cli.Search().
		Index(indexName). // 指定索引名称
		Query(
			elastic.NewBoolQuery().
				Should(
					elastic.NewMatchQuery("classname", keyWords), // 匹配 classname
					elastic.NewMatchQuery("teacher", keyWords),   // 匹配 teacher
				).
				MinimumShouldMatch("1"). // 至少匹配一个条件
				Filter(
					elastic.NewTermQuery("year", xnm),     // year 精确匹配 xnm
					elastic.NewTermQuery("semester", xqm), // semester 精确匹配 xqm
				),
		).Do(ctx) // 执行查询

	if err != nil {
		clog.LogPrinter.Errorf("es: failed to search class_info[keywords:%v xnm:%v xqm:%v]: %v", keyWords, xnm, xqm, err)
		return nil, errcode.Err_EsSearchClassInfo
	}
	for _, hit := range searchResult.Hits.Hits {
		var classInfo biz.ClassInfo
		err := json.Unmarshal(hit.Source, &classInfo)
		if err != nil {
			clog.LogPrinter.Errorf("json unmarshal %v failed: %v", hit.Source, err)
			continue
		}
		classInfos = append(classInfos, classInfo)
	}
	return classInfos, nil
}
