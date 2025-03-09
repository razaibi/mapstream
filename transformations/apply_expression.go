package transformations

import (
	"fmt"

	"github.com/blues/jsonata-go"
	"github.com/gofiber/fiber/v2"
)

type ApplyExpression struct{}

func (a ApplyExpression) TransformMap(
	c *fiber.Ctx,
	data map[string]interface{},
	params map[string]interface{},
) (map[string]interface{}, error) {
	expression, ok := params["expression"].(string)
	if !ok {
		return nil, fmt.Errorf("expression is not a string or not provided in params")
	}

	// Compile the expression
	e, err := jsonata.Compile(expression)
	if err != nil {
		return nil, fmt.Errorf("failed to compile expression: %w", err)
	}

	// Evaluate the expression on the provided data
	res, err := e.Eval(data)
	if err != nil {
		return nil, fmt.Errorf("failed to compile expression: %w", err)
	}

	// Perform a type assertion
	resMap := make(map[string]interface{})
	resMap["result"] = res

	// Return the asserted result
	return resMap, nil

}
