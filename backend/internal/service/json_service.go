package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

type JsonService struct{}

func NewJsonService() *JsonService { return &JsonService{} }

func (s *JsonService) Validate(input string) error {
	var v any
	return json.Unmarshal([]byte(input), &v)
}

func (s *JsonService) Format(input string, indent int) (string, error) {
	var v any
	if err := json.Unmarshal([]byte(input), &v); err != nil {
		return "", err
	}
	prefix := strings.Repeat(" ", indent)
	b, err := json.MarshalIndent(v, "", prefix)
	return string(b), err
}

func (s *JsonService) Minify(input string) (string, error) {
	var v any
	if err := json.Unmarshal([]byte(input), &v); err != nil {
		return "", err
	}
	b, err := json.Marshal(v)
	return string(b), err
}

func (s *JsonService) Convert(input, target string) (string, error) {
	var v any
	if err := json.Unmarshal([]byte(input), &v); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}
	switch strings.ToLower(target) {
	case "yaml":
		b, err := yaml.Marshal(v)
		return string(b), err
	case "toml":
		return jsonToTOML(v)
	case "xml":
		return jsonToXML(v)
	default:
		return "", errors.New("unsupported target format: " + target)
	}
}

func (s *JsonService) Parse(input, source string) (string, error) {
	switch strings.ToLower(source) {
	case "yaml":
		var v any
		if err := yaml.Unmarshal([]byte(input), &v); err != nil {
			return "", err
		}
		v = convertYAMLToJSON(v)
		b, err := json.MarshalIndent(v, "", "  ")
		return string(b), err
	default:
		return "", errors.New("unsupported source format: " + source)
	}
}

func (s *JsonService) GenerateSchema(input string) (string, error) {
	var v any
	if err := json.Unmarshal([]byte(input), &v); err != nil {
		return "", err
	}
	schema := inferSchema(v)
	b, err := json.MarshalIndent(schema, "", "  ")
	return string(b), err
}

func (s *JsonService) ValidateSchema(schemaStr, dataStr string) ([]string, error) {
	var schema, data any
	if err := json.Unmarshal([]byte(schemaStr), &schema); err != nil {
		return nil, fmt.Errorf("invalid schema: %w", err)
	}
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		return nil, fmt.Errorf("invalid data: %w", err)
	}
	return validateAgainstSchema(schema.(map[string]any), data), nil
}

func (s *JsonService) Diff(a, b string) (any, error) {
	var va, vb any
	if err := json.Unmarshal([]byte(a), &va); err != nil {
		return nil, fmt.Errorf("invalid JSON a: %w", err)
	}
	if err := json.Unmarshal([]byte(b), &vb); err != nil {
		return nil, fmt.Errorf("invalid JSON b: %w", err)
	}
	return diffValues("", va, vb), nil
}

func (s *JsonService) Query(input, path string) (any, error) {
	var v any
	if err := json.Unmarshal([]byte(input), &v); err != nil {
		return nil, err
	}
	return jsonPath(v, path)
}
