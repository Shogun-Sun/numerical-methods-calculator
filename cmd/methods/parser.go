package methods

import (
	"fmt"
	"math"

	"github.com/Knetic/govaluate"
)

func MakeFunctionFromString(exprStr string) (func(float64) float64, error) {

	functions := map[string]govaluate.ExpressionFunction{
		"ln": func(args ...interface{}) (interface{}, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("функция ln требует 1 аргумент")
			}
			return math.Log(args[0].(float64)), nil
		},
		"sin": func(args ...interface{}) (interface{}, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("функция sin требует 1 аргумент")
			}
			return math.Sin(args[0].(float64)), nil
		},
		"cos": func(args ...interface{}) (interface{}, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("функция cos требует 1 аргумент")
			}
			return math.Cos(args[0].(float64)), nil
		},
		"pow": func(args ...interface{}) (interface{}, error) {
			if len(args) < 2 {
				return nil, fmt.Errorf("функция pow требует 2 аргумента, например pow(x, 2)")
			}
			return math.Pow(args[0].(float64), args[1].(float64)), nil
		},
	}

	expression, err := govaluate.NewEvaluableExpressionWithFunctions(exprStr, functions)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга выражения: %v", err)
	}

	return func(x float64) float64 {
		parameters := make(map[string]interface{}, 1)
		parameters["x"] = x

		result, err := expression.Evaluate(parameters)
		if err != nil {
			return math.NaN()
		}

		return result.(float64)
	}, nil
}
