package errcode

import (
	v1 "github.com/asynccnu/be-api/gen/proto/classService/v1"
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	Err_EsAddClassInfo    = errors.New(450, v1.ErrorReason_ES_AddClassFailed.String(), "创建classInfo失败")
	Err_EsSearchClassInfo = errors.New(451, v1.ErrorReason_ES_SearchClassFailed.String(), "查询classInfo失败")
)
