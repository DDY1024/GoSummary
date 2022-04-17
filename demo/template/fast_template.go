package template

import (
	"sync"

	"github.com/Knetic/govaluate"
	"github.com/valyala/fasttemplate"
)

// fasttempalte 通常结合着 govaluate 一起使用，实现一个简单的规则引擎

const (
	defaultStartTag = "{{"
	defaultEndTag   = "}}"
)

var (
	templateCache sync.Map
	exprCache     sync.Map
)

// GetTemplate 生成模板，startTag 和 endTag 采用官方建议的 {{ 和 }} 即可，不再另外指定
func GetTemplate(template string) (*fasttemplate.Template, error) {
	if val, ok := templateCache.Load(template); ok {
		return val.(*fasttemplate.Template), nil
	}

	fastTemplate, err := fasttemplate.NewTemplate(template, defaultStartTag, defaultEndTag)
	if err != nil {
		return nil, err
	}
	defer templateCache.Store(template, fastTemplate)

	return fastTemplate, nil
}

func GetExpression(expr string) (*govaluate.EvaluableExpression, error) {
	if val, ok := exprCache.Load(expr); ok {
		return val.(*govaluate.EvaluableExpression), nil
	}

	evalExpr, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		return nil, err
	}
	defer exprCache.Store(expr, evalExpr)

	return evalExpr, nil
}
