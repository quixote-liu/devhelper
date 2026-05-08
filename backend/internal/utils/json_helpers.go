package utils

import (
	"fmt"
	"strings"
)

// InferSchema generates a basic JSON Schema from a value.
func InferSchema(v any) map[string]any {
	schema := map[string]any{"$schema": "http://json-schema.org/draft-07/schema#"}
	fillSchema(schema, v)
	return schema
}

func fillSchema(s map[string]any, v any) {
	switch val := v.(type) {
	case map[string]any:
		s["type"] = "object"
		props := map[string]any{}
		required := []string{}
		for k, child := range val {
			childSchema := map[string]any{}
			fillSchema(childSchema, child)
			props[k] = childSchema
			required = append(required, k)
		}
		s["properties"] = props
		if len(required) > 0 {
			s["required"] = required
		}
	case []any:
		s["type"] = "array"
		if len(val) > 0 {
			items := map[string]any{}
			fillSchema(items, val[0])
			s["items"] = items
		}
	case string:
		s["type"] = "string"
	case float64:
		if val == float64(int64(val)) {
			s["type"] = "integer"
		} else {
			s["type"] = "number"
		}
	case bool:
		s["type"] = "boolean"
	case nil:
		s["type"] = "null"
	}
}

// ValidateAgainstSchema performs basic JSON Schema validation.
func ValidateAgainstSchema(schema map[string]any, data any) []string {
	var errs []string
	validateValue(schema, data, "$", &errs)
	return errs
}

func validateValue(schema map[string]any, data any, path string, errs *[]string) {
	schemaType, _ := schema["type"].(string)
	switch schemaType {
	case "object":
		obj, ok := data.(map[string]any)
		if !ok {
			*errs = append(*errs, fmt.Sprintf("%s: expected object", path))
			return
		}
		if required, ok := schema["required"].([]any); ok {
			for _, r := range required {
				key := fmt.Sprint(r)
				if _, exists := obj[key]; !exists {
					*errs = append(*errs, fmt.Sprintf("%s.%s: required field missing", path, key))
				}
			}
		}
		if props, ok := schema["properties"].(map[string]any); ok {
			for k, propSchema := range props {
				if ps, ok := propSchema.(map[string]any); ok {
					if v, exists := obj[k]; exists {
						validateValue(ps, v, path+"."+k, errs)
					}
				}
			}
		}
	case "array":
		arr, ok := data.([]any)
		if !ok {
			*errs = append(*errs, fmt.Sprintf("%s: expected array", path))
			return
		}
		if items, ok := schema["items"].(map[string]any); ok {
			for i, item := range arr {
				validateValue(items, item, fmt.Sprintf("%s[%d]", path, i), errs)
			}
		}
	case "string":
		if _, ok := data.(string); !ok {
			*errs = append(*errs, fmt.Sprintf("%s: expected string", path))
		}
	case "number", "integer":
		if _, ok := data.(float64); !ok {
			*errs = append(*errs, fmt.Sprintf("%s: expected number", path))
		}
	case "boolean":
		if _, ok := data.(bool); !ok {
			*errs = append(*errs, fmt.Sprintf("%s: expected boolean", path))
		}
	}
}

// DiffEntry represents a single diff entry between two JSON values.
type DiffEntry struct {
	Path string `json:"path"`
	Type string `json:"type"` // "added", "removed", "changed"
	Old  any    `json:"old,omitempty"`
	New  any    `json:"new,omitempty"`
}

// DiffValues computes a structural diff between two JSON values.
func DiffValues(path string, a, b any) []DiffEntry {
	var diffs []DiffEntry
	aMap, aIsMap := a.(map[string]any)
	bMap, bIsMap := b.(map[string]any)

	if aIsMap && bIsMap {
		for k, av := range aMap {
			p := path + "." + k
			if bv, ok := bMap[k]; ok {
				diffs = append(diffs, DiffValues(p, av, bv)...)
			} else {
				diffs = append(diffs, DiffEntry{Path: p, Type: "removed", Old: av})
			}
		}
		for k, bv := range bMap {
			if _, ok := aMap[k]; !ok {
				diffs = append(diffs, DiffEntry{Path: path + "." + k, Type: "added", New: bv})
			}
		}
		return diffs
	}

	if fmt.Sprint(a) != fmt.Sprint(b) {
		diffs = append(diffs, DiffEntry{Path: path, Type: "changed", Old: a, New: b})
	}
	return diffs
}

// JsonPath evaluates a simple dot-notation path like $.store.book[0].title
func JsonPath(v any, path string) (any, error) {
	path = strings.TrimPrefix(path, "$")
	parts := tokenizePath(path)
	return traversePath(v, parts)
}

func tokenizePath(path string) []string {
	var parts []string
	for _, p := range strings.Split(path, ".") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if idx := strings.Index(p, "["); idx != -1 {
			key := p[:idx]
			rest := p[idx:]
			if key != "" {
				parts = append(parts, key)
			}
			for rest != "" {
				end := strings.Index(rest, "]")
				if end == -1 {
					break
				}
				parts = append(parts, rest[1:end])
				rest = rest[end+1:]
			}
		} else {
			parts = append(parts, p)
		}
	}
	return parts
}

func traversePath(v any, parts []string) (any, error) {
	if len(parts) == 0 {
		return v, nil
	}
	key := parts[0]
	switch val := v.(type) {
	case map[string]any:
		child, ok := val[key]
		if !ok {
			return nil, fmt.Errorf("key %q not found", key)
		}
		return traversePath(child, parts[1:])
	case []any:
		var idx int
		if _, err := fmt.Sscanf(key, "%d", &idx); err != nil {
			return nil, fmt.Errorf("expected array index, got %q", key)
		}
		if idx < 0 || idx >= len(val) {
			return nil, fmt.Errorf("index %d out of range", idx)
		}
		return traversePath(val[idx], parts[1:])
	default:
		return nil, fmt.Errorf("cannot traverse into %T with key %q", v, key)
	}
}

// ConvertYAMLToJSON converts yaml.v2 map types to JSON-compatible map[string]any.
func ConvertYAMLToJSON(v any) any {
	switch val := v.(type) {
	case map[any]any:
		m := map[string]any{}
		for k, child := range val {
			m[fmt.Sprint(k)] = ConvertYAMLToJSON(child)
		}
		return m
	case []any:
		for i, item := range val {
			val[i] = ConvertYAMLToJSON(item)
		}
		return val
	default:
		return v
	}
}

// JsonToTOML produces a simple TOML representation.
func JsonToTOML(v any) (string, error) {
	obj, ok := v.(map[string]any)
	if !ok {
		return "", fmt.Errorf("TOML conversion requires a JSON object at the root")
	}
	var sb strings.Builder
	writeTOMLSection(&sb, obj, "")
	return sb.String(), nil
}

func writeTOMLSection(sb *strings.Builder, obj map[string]any, prefix string) {
	for k, val := range obj {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch child := val.(type) {
		case map[string]any:
			sb.WriteString(fmt.Sprintf("[%s]\n", key))
			writeTOMLSection(sb, child, "")
		case string:
			sb.WriteString(fmt.Sprintf("%s = %q\n", k, child))
		case float64:
			sb.WriteString(fmt.Sprintf("%s = %v\n", k, child))
		case bool:
			sb.WriteString(fmt.Sprintf("%s = %v\n", k, child))
		case nil:
			// skip null values
		default:
			sb.WriteString(fmt.Sprintf("# %s = (complex value)\n", k))
		}
	}
}

// JsonToXML produces a simple XML representation.
func JsonToXML(v any) (string, error) {
	var sb strings.Builder
	sb.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	sb.WriteString("<root>\n")
	writeXMLValue(&sb, v, 1)
	sb.WriteString("</root>")
	return sb.String(), nil
}

func writeXMLValue(sb *strings.Builder, v any, depth int) {
	indent := strings.Repeat("  ", depth)
	switch val := v.(type) {
	case map[string]any:
		for k, child := range val {
			sb.WriteString(fmt.Sprintf("%s<%s>", indent, k))
			switch c := child.(type) {
			case map[string]any, []any:
				sb.WriteString("\n")
				writeXMLValue(sb, c, depth+1)
				sb.WriteString(fmt.Sprintf("%s</%s>\n", indent, k))
			default:
				sb.WriteString(fmt.Sprintf("%v</%s>\n", c, k))
			}
		}
	case []any:
		for _, item := range val {
			sb.WriteString(fmt.Sprintf("%s<item>", indent))
			switch c := item.(type) {
			case map[string]any, []any:
				sb.WriteString("\n")
				writeXMLValue(sb, c, depth+1)
				sb.WriteString(fmt.Sprintf("%s</item>\n", indent))
			default:
				sb.WriteString(fmt.Sprintf("%v</item>\n", c))
			}
		}
	}
}
