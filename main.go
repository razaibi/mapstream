package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"mapstream/transformations"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

type Endpoint struct {
	Name     string                 `json:"name"`
	Url      string                 `json:"url"`
	Method   string                 `json:"method"`
	Pipeline []transformations.Step `json:"pipeline"`
}

type PipelineConfig struct {
	Endpoints []Endpoint `json:"endpoints"`
}

func setPostRoute(app *fiber.App, endpoint Endpoint) {
	app.Post(endpoint.Url, func(c *fiber.Ctx) error {
		fmt.Println("Endpoint path:", c.Path())
		var inputData interface{}

		body := c.Body()
		if len(body) > 0 {
			if err := c.BodyParser(&inputData); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid input: Unable to parse JSON",
				})
			}
		} else {
			inputData = make(map[string]interface{})
		}

		data, err := transformations.ProcessPipeline(c, inputData, endpoint.Pipeline)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(data)
	})
}

func setGetRoute(app *fiber.App, endpoint Endpoint) {
	app.Get(endpoint.Url, func(c *fiber.Ctx) error {
		fmt.Println("Endpoint path:", c.Path())
		inputData := make(map[string]interface{})

		data, err := transformations.ProcessPipeline(c, inputData, endpoint.Pipeline)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(data)
	})
}

func main() {
	config := loadPipelineConfig("pipeline_alpha.json")

	app := fiber.New(
		fiber.Config{
			JSONEncoder: sonic.Marshal,
			JSONDecoder: sonic.Unmarshal,
		},
	)

	for _, endpoint := range config.Endpoints {
		endpoint := endpoint // capture range variable
		//

		switch strings.ToLower(endpoint.Method) {
		case "get":
			setGetRoute(app, endpoint)
		case "post":
			setPostRoute(app, endpoint)
		default:
			setGetRoute(app, endpoint)
		}
	}

	log.Fatal(app.Listen(":8080"))
}

func loadPipelineConfig(path string) PipelineConfig {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Could not read config file: %v", err)
	}

	var config PipelineConfig
	if err := json.Unmarshal(file, &config); err != nil {
		log.Fatalf("Could not read config file: %v", err)
	}

	return config
}
