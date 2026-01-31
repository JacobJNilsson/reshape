package core

import (
	"errors"
	"strings"
)

func splitPath(path string) []string {
	cleaned := strings.TrimSpace(path)
	if cleaned == "" {
		return []string{}
	}
	return strings.Split(cleaned, ".")
}

func getValueAtPath(record Record, path string) (any, bool, error) {
	segments := splitPath(path)
	if len(segments) == 0 {
		return nil, false, errors.New("path is empty")
	}
	current := any(record)
	for index, segment := range segments {
		currentMap, ok := mapFromValue(current)
		if !ok {
			return nil, false, errors.New("path segment is not an object: " + segment)
		}
		remainingPath := strings.Join(segments[index:], ".")
		if value, exists := currentMap[remainingPath]; exists {
			return value, true, nil
		}
		value, exists := currentMap[segment]
		if !exists {
			return nil, false, nil
		}
		current = value
	}
	return current, true, nil
}

// ValueAtPath retrieves a value using dot notation paths.
func ValueAtPath(record Record, path string) (any, bool, error) {
	return getValueAtPath(record, path)
}

func setValueAtPath(record Record, path string, value any) error {
	segments := splitPath(path)
	if len(segments) == 0 {
		return errors.New("path is empty")
	}
	current := map[string]any(record)
	for index, segment := range segments {
		remainingPath := strings.Join(segments[index:], ".")
		if _, exists := current[remainingPath]; exists {
			current[remainingPath] = value
			return nil
		}
		if index == len(segments)-1 {
			current[segment] = value
			return nil
		}
		next, ok := current[segment]
		if !ok {
			nextMap := map[string]any{}
			current[segment] = nextMap
			current = nextMap
			continue
		}
		nextMap, ok := mapFromValue(next)
		if !ok {
			return errors.New("path segment is not an object: " + segment)
		}
		current = nextMap
	}
	return nil
}

func deleteValueAtPath(record Record, path string) error {
	segments := splitPath(path)
	if len(segments) == 0 {
		return errors.New("path is empty")
	}
	current := map[string]any(record)
	for index, segment := range segments {
		if index == len(segments)-1 {
			delete(current, segment)
			return nil
		}
		next, ok := current[segment]
		if !ok {
			return nil
		}
		nextMap, ok := mapFromValue(next)
		if !ok {
			return errors.New("path segment is not an object: " + segment)
		}
		current = nextMap
	}
	return nil
}

func flattenAtPath(record Record, path string) error {
	value, exists, err := getValueAtPath(record, path)
	if err != nil {
		return err
	}
	if !exists || value == nil {
		return nil
	}
	objectValue, ok := mapFromValue(value)
	if !ok {
		return errors.New("flatten target is not an object: " + path)
	}
	for key, nestedValue := range objectValue {
		flatKey := path + "." + key
		record[flatKey] = nestedValue
	}
	return deleteValueAtPath(record, path)
}

func mapFromValue(value any) (map[string]any, bool) {
	switch typed := value.(type) {
	case map[string]any:
		return typed, true
	case Record:
		return map[string]any(typed), true
	default:
		return nil, false
	}
}
