package transformations

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type SelectAttributes struct{}

func (s SelectAttributes) Transform(
	ctx *fiber.Ctx,
	data interface{},
	params map[string]interface{},
) (interface{}, error) {
	// Retrieve and validate the columns
	columns, columnsExist := params["columns"].([]interface{})
	if !columnsExist {
		return data, fmt.Errorf("columns paramter is missing or not a slice")
	}

	columnKeys := make(map[string]struct{})
	for _, col := range columns {
		if colStr, ok := col.(string); ok {
			columnKeys[colStr] = struct{}{}
		} else {
			return data, fmt.Errorf("columns parameter contains non-string values")
		}
	}

	switch d := data.(type) {
	case []map[string]interface{}:
		return applySelectAttributes(d, columnKeys, params)
	case map[string]interface{}:
		return selectFromMap(d, columnKeys), nil
	default:
		return data, fmt.Errorf("data is of unsupported type, must be a map or slice of maps")
	}
}

// Helper to apply selection for a slice of maps
func applySelectAttributes(
	dataSlice []map[string]interface{},
	columnKeys map[string]struct{},
	params map[string]interface{},
) (interface{}, error) {
	offset, limit := parseOffsetLimit(params, len(dataSlice))

	//Adjust slice based on on offset and limit
	if offset > len(dataSlice) {
		offset = len(dataSlice)
	}
	if limit > len(dataSlice) {
		limit = len(dataSlice)
	}

	selectedDataSlice := dataSlice[offset:limit]

	// Filter each map to include only specified columns
	for i, row := range selectedDataSlice {
		selectedDataSlice[i] = selectFromMap(row, columnKeys)
	}

	return selectedDataSlice, nil
}

// Helper to select columns from a single map
func selectFromMap(
	data map[string]interface{},
	columnKeys map[string]struct{},
) map[string]interface{} {
	filteredRow := make(map[string]interface{})
	for col := range columnKeys {
		if val, exists := data[col]; exists {
			filteredRow[col] = val
		}
	}
	return filteredRow
}

// Extracts and verifies offset and limit parameters
func parseOffsetLimit(
	params map[string]interface{},
	length int,
) (int, int) {
	// Default offset is 0
	offset, ok := params["offset"].(int)
	if !ok || offset < 0 {
		offset = 0
	}

	// Default limit is the length of the data
	limit, ok := params["limit"].(int)
	if !ok || limit < 0 {
		limit = length
	}

	return offset, offset + limit
}
