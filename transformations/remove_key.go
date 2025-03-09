package transformations

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type RemoveKey struct{}

// Implement the TransformMap method to remove a key from a map
func (r RemoveKey) TransformMap(
	c *fiber.Ctx,
	data map[string]interface{},
	params map[string]interface{},
) (map[string]interface{}, error) {
	key, keyExists := params["key"].(string)
	if !keyExists {
		return data, fmt.Errorf("key parameter is missing or not a string")
	}

	// Check if the key exists in the data map and remove it
	if _, exists := data[key]; exists {
		delete(data, key)
	} else {
		return data, fmt.Errorf("key '%s' does not exist in data", key)
	}

	return data, nil
}
