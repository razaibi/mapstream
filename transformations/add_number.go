package transformations

import "github.com/gofiber/fiber/v2"

type AddNumber struct{}

func (a AddNumber) TransformMap(
	c *fiber.Ctx,
	data map[string]interface{},
	params map[string]interface{},
) (map[string]interface{}, error) {
	key := params["key"].(string)
	addValue, _ := params["value"].(float64)
	if currentValue, exists := data[key].(float64); exists {
		data[key] = currentValue + addValue
	}

	return data, nil
}
