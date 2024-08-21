package data

import (
	"classService/internal/biz"
	"classService/internal/errcode"
	"classService/internal/logPrinter"
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
)

const indexName = "class_info"

type Es struct {
	cli *elastic.Client
	log logPrinter.LogerPrinter
}

func (es Es) AddClassInfo(ctx context.Context, classInfo biz.ClassInfo) error {
	// 创建文档
	_, err := es.cli.Index().
		Index(indexName).
		BodyJson(classInfo).
		Do(ctx)
	if err != nil {
		es.log.FuncError(es.cli.Index().
			Index(indexName).
			BodyJson(classInfo).
			Do, err)
		return errcode.Err_EsAddClassInfo
	}
	return nil
}
func (es Es) SearchClassInfo(ctx context.Context, keyWords string) ([]biz.ClassInfo, error) {
	var classInfos = make([]biz.ClassInfo, 0)
	searchResult, err := es.cli.Search().
		Index(indexName). // 指定索引名称
		Query(elastic.NewBoolQuery().
			Should(
				elastic.NewMatchQuery("classname", keyWords),
				elastic.NewMatchQuery("teacher", keyWords),
			).
			MinimumShouldMatch("1"), // 至少匹配一个条件
		).Do(ctx) // 执行查询
	if err != nil {
		es.log.FuncError(es.cli.Search().
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
			es.log.FuncError(json.Unmarshal, err)
			continue
		}
		classInfos = append(classInfos, classInfo)
	}
	return classInfos, nil
}
