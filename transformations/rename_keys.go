package transformations

import (
	"github.com/gofiber/fiber/v2"
)

type RenameKeys struct{}

// Implement the TransformMap method to remove a key from a map
func (r RenameKeys) TransformMap(
	c *fiber.Ctx,
	data map[string]interface{},
	params map[string]interface{},
) (map[string]interface{}, error) {
	oldKey := params["old_key"].(string)
	newKey := params["new_key"].(string)
	if value, exists := data[oldKey]; exists {
		data[newKey] = value
		delete(data, oldKey)
	}
	return data, nil
}
