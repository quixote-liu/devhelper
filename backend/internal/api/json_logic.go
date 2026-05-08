package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"devhelper/internal/utils"

	"gopkg.in/yaml.v2"
)

func validateJSON(input string) error {
	var v any
	return json.Unmarshal([]byte(input), &v)
}

func formatJSON(input string, indent int) (string, error) {
	var v any
	if err := json.Unmarshal([]byte(input), &v); err != nil {
		return "", err
	}
	prefix := strings.Repeat(" ", indent)
	b, err := json.MarshalIndent(v, "", prefix)
	return string(b), err
}

func minifyJSON(input string) (string, error) {
	var v any
	if err := json.Unmarshal([]byte(input), &v); err != nil {
		return "", err
	}
	b, err := json.Marshal(v)
	return string(b), err
}

func convertJSON(input, target string) (string, error) {
	var v any
	if err := json.Unmarshal([]byte(input), &v); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}
	switch strings.ToLower(target) {
	case "yaml":
		b, err := yaml.Marshal(v)
		return string(b), err
	case "toml":
		return utils.JsonToTOML(v)
	case "xml":
		return utils.JsonToXML(v)
	default:
		return "", errors.New("unsupported target format: " + target)
	}
}

func parseJSON(input, source string) (string, error) {
	switch strings.ToLower(source) {
	case "yaml":
		var v any
		if err := yaml.Unmarshal([]byte(input), &v); err != nil {
			return "", err
		}
		v = utils.ConvertYAMLToJSON(v)
		b, err := json.MarshalIndent(v, "", "  ")
		return string(b), err
	default:
		return "", errors.New("unsupported source format: " + source)
	}
}

func generateSchema(input string) (string, error) {
	var v any
	if err := json.Unmarshal([]byte(input), &v); err != nil {
		return "", err
	}
	schema := utils.InferSchema(v)
	b, err := json.MarshalIndent(schema, "", "  ")
	return string(b), err
}

func validateSchema(schemaStr, dataStr string) ([]string, error) {
	var schema, data any
	if err := json.Unmarshal([]byte(schemaStr), &schema); err != nil {
		return nil, fmt.Errorf("invalid schema: %w", err)
	}
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		return nil, fmt.Errorf("invalid data: %w", err)
	}
	return utils.ValidateAgainstSchema(schema.(map[string]any), data), nil
}

func diffJSON(a, b string) (any, error) {
	var va, vb any
	if err := json.Unmarshal([]byte(a), &va); err != nil {
		return nil, fmt.Errorf("invalid JSON a: %w", err)
	}
	if err := json.Unmarshal([]byte(b), &vb); err != nil {
		return nil, fmt.Errorf("invalid JSON b: %w", err)
	}
	return utils.DiffValues("", va, vb), nil
}

func queryJSON(input, path string) (any, error) {
	var v any
	if err := json.Unmarshal([]byte(input), &v); err != nil {
		return nil, err
	}
	return utils.JsonPath(v, path)
}
