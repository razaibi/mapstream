package transformations

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

type ExternalAPICall struct {
}

func (e ExternalAPICall) TransformMap(
	c *fiber.Ctx,
	data map[string]interface{},
	params map[string]interface{},
) (map[string]interface{}, error) {
	client := resty.New()
	url := params["url"].(string)
	method := strings.ToUpper(params["method"].(string))

	requestHeaders := make(map[string]string)
	for key := range c.GetReqHeaders() {
		requestHeaders[key] = c.Get(key)
	}

	if callHeaders, ok := params["headers"].(map[string]interface{}); ok {
		for key, value := range callHeaders {
			requestHeaders[key] = value.(string)
		}
	}

	response, err := client.R().
		SetHeaders(requestHeaders).
		SetBody(data).
		Execute(method, url)

	if err != nil {
		return nil, fmt.Errorf("failed to call external API: %v", err)
	}

	var apiData map[string]interface{}
	if err := json.Unmarshal(response.Body(), &apiData); err != nil {
		return nil, fmt.Errorf("failed to parse response from external API: %v", err)
	}

	for key, value := range apiData {
		data[key] = value
	}

	return data, nil
}
