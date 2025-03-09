package transformations

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

var transformations = map[string]FlexibleTransformation{
	"read_json":          MapTransformation{TransformMap: ReadJSONFile{}.TransformMap},
	"rename_keys":        MapTransformation{TransformMap: RenameKeys{}.TransformMap},
	"remove_key":         MapTransformation{TransformMap: RemoveKey{}.TransformMap},
	"add_number":         MapTransformation{TransformMap: AddNumber{}.TransformMap},
	"external_api_call":  MapTransformation{TransformMap: ExternalAPICall{}.TransformMap},
	"apply_expression":   MapTransformation{TransformMap: ApplyExpression{}.TransformMap},
	"conditional_action": ConditionalTransformation{},
	"select":             SelectAttributes{},
}

type Step struct {
	Step   string                 `json:"step"`
	Params map[string]interface{} `json:"params"`
}

type MapTransformation struct {
	TransformMap func(
		*fiber.Ctx,
		map[string]interface{},
		map[string]interface{},
	) (map[string]interface{}, error)
}

func (mt MapTransformation) Transform(
	c *fiber.Ctx,
	data interface{},
	params map[string]interface{},
) (interface{}, error) {
	switch d := data.(type) {
	case map[string]interface{}:
		return mt.TransformMap(c, d, params)

	case []map[string]interface{}:
		for i, item := range d {
			transformed, err := mt.TransformMap(c, item, params)
			if err != nil {
				return nil, err
			}
			d[i] = transformed
		}
		return d, nil
	default:
		return nil, fmt.Errorf("unsupported data type for transformation")
	}
}

type FlexibleTransformation interface {
	Transform(
		*fiber.Ctx,
		interface{},
		map[string]interface{},
	) (interface{}, error)
}

func ProcessPipeline(
	c *fiber.Ctx,
	data interface{},
	steps []Step,
) (interface{}, error) {
	for _, step := range steps {
		if transformation, found := transformations[step.Step]; found {
			var err error
			data, err = transformation.Transform(c, data, step.Params)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("unknown step: %s", step.Step)
		}
	}
	return data, nil
}
