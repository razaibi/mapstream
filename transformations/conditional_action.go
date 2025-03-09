package transformations

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type ConditionalTransformation struct{}

func (c ConditionalTransformation) Transform(
	ctx *fiber.Ctx,
	data interface{},
	params map[string]interface{},
) (interface{}, error) {
	key, keyOk := params["key"].(string)
	checkValue, checkValueOk := params["value"]
	action, actionOk := params["actions"].(string)

	if !keyOk || !checkValueOk || !actionOk {
		return data, fmt.Errorf("missing required condition parameters: 'key', 'value', or 'action'")
	}

	actionParams, actionParamsOk := params["actionParams"].(map[string]interface{})
	if !actionParamsOk {
		return data, fmt.Errorf("actionParams should be a map[string]interface{}")
	}

	actionTransformation, found := transformations[action]
	if !found {
		return data, fmt.Errorf("action transformation '%s' not found", action)
	}

	switch d := data.(type) {
	case map[string]interface{}:
		return applyTransformationIfConditionMet(
			ctx,
			d,
			key,
			checkValue,
			actionTransformation,
			actionParams,
		)

	case []map[string]interface{}:
		for i, item := range d {
			transformed, err := applyTransformationIfConditionMet(
				ctx,
				item,
				key,
				checkValue,
				actionTransformation,
				actionParams,
			)
			if err != nil {
				return nil, err
			}
			d[i] = transformed.(map[string]interface{})
		}
		return d, nil
	default:
		return nil, fmt.Errorf("unsupported data type for transformation; expecting map or slice of maps")
	}
}

func applyTransformationIfConditionMet(
	ctx *fiber.Ctx,
	data map[string]interface{},
	key string,
	checkValue interface{},
	actionTransformation FlexibleTransformation,
	actionParams map[string]interface{},
) (interface{}, error) {
	if value, exists := data[key]; exists && value == checkValue {
		return actionTransformation.Transform(ctx, data, actionParams)
	}
	return data, nil
}
