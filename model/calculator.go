package model

import (
	"github.com/zgwit/iot-master/calc"
)

//Calculator 计算器
type Calculator struct {
	Expression string `json:"expression"`
	As         string `json:"as"`

	expr *calc.Expression
	//ctx  *Context
}

//Init 初始化（编译）
func (c *Calculator) Init() (err error) {
	c.expr, err = calc.NewExpression(c.Expression)
	return
}

//Evaluate 计算
func (c *Calculator) Evaluate(ctx calc.Context) (float64, error) {
	return c.expr.Evaluate(ctx)
}
