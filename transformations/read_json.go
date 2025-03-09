package transformations

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

type ReadJSONFile struct{}

func (a ReadJSONFile) TransformMap(
	c *fiber.Ctx,
	data map[string]interface{},
	params map[string]interface{},
) (map[string]interface{}, error) {
	location, ok := params["location"].(string)
	if !ok || location == "" {
		return nil, fmt.Errorf("invalid of missing location value in params")
	}

	// Read the file contents
	fileData, err := os.ReadFile(location)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	//Unmarshal the JSON data into a map
	var result map[string]interface{}
	err = json.Unmarshal(fileData, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data: %w", err)
	}

	return result, nil
}
