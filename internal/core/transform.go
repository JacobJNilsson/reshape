package core

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// TransformData applies the conversion plan to canonical data.
func TransformData(input CanonicalData, plan ConversionPlan) (CanonicalData, []Warning, error) {
	normalizedPlan := NormalizePlan(plan)
	if err := ValidateLossyOperations(normalizedPlan); err != nil {
		return CanonicalData{}, nil, err
	}
	records := deepCopyRecords(input.Values.Records)
	warnings := []Warning{}
	warningSet := map[string]struct{}{}

	for _, path := range normalizedPlan.FlattenFields {
		for index := range records {
			if err := flattenAtPath(records[index], path); err != nil {
				return CanonicalData{}, nil, err
			}
		}
	}

	for _, path := range normalizedPlan.ExplodeArrays {
		var expanded []Record
		for _, record := range records {
			value, exists, err := getValueAtPath(record, path)
			if err != nil {
				return CanonicalData{}, nil, err
			}
			if !exists || value == nil {
				expanded = append(expanded, record)
				continue
			}
			sliceValue, ok := value.([]any)
			if !ok {
				return CanonicalData{}, nil, errors.New("explode target is not an array: " + path)
			}
			if len(sliceValue) == 0 {
				expanded = append(expanded, record)
				continue
			}
			for _, item := range sliceValue {
				copied := deepCopyRecord(record)
				if err := setValueAtPath(copied, path, item); err != nil {
					return CanonicalData{}, nil, err
				}
				expanded = append(expanded, copied)
			}
		}
		records = expanded
	}

	for _, rule := range normalizedPlan.JoinArrays {
		for index := range records {
			value, exists, err := getValueAtPath(records[index], rule.Path)
			if err != nil {
				return CanonicalData{}, nil, err
			}
			if !exists || value == nil {
				continue
			}
			sliceValue, ok := value.([]any)
			if !ok {
				return CanonicalData{}, nil, errors.New("join target is not an array: " + rule.Path)
			}
			joined, err := joinArrayValues(sliceValue, rule.Delimiter)
			if err != nil {
				return CanonicalData{}, nil, err
			}
			if err := setValueAtPath(records[index], rule.Path, joined); err != nil {
				return CanonicalData{}, nil, err
			}
		}
		addWarningOnce(&warnings, warningSet, rule.Path, "joined array into string")
	}

	for _, rule := range normalizedPlan.TypeCoercions {
		for index := range records {
			value, exists, err := getValueAtPath(records[index], rule.Path)
			if err != nil {
				return CanonicalData{}, nil, err
			}
			if !exists || value == nil {
				continue
			}
			coerced, err := coerceValue(value, rule.TargetType)
			if err != nil {
				return CanonicalData{}, nil, err
			}
			if err := setValueAtPath(records[index], rule.Path, coerced); err != nil {
				return CanonicalData{}, nil, err
			}
		}
		addWarningOnce(&warnings, warningSet, rule.Path, "coerced type")
	}

	for _, rule := range normalizedPlan.DefaultValues {
		for index := range records {
			value, exists, err := getValueAtPath(records[index], rule.Path)
			if err != nil {
				return CanonicalData{}, nil, err
			}
			if !exists || value == nil {
				if err := setValueAtPath(records[index], rule.Path, rule.Value); err != nil {
					return CanonicalData{}, nil, err
				}
			}
		}
	}

	for _, path := range normalizedPlan.DropFields {
		for index := range records {
			if err := deleteValueAtPath(records[index], path); err != nil {
				return CanonicalData{}, nil, err
			}
		}
		addWarningOnce(&warnings, warningSet, path, "dropped field")
	}

	output := CanonicalData{Values: DataValues{Records: records}}
	if len(records) > 0 {
		output.Shape = BuildShapeFromRecords(records)
	} else {
		output.Shape = input.Shape
	}
	return output, warnings, nil
}

func joinArrayValues(values []any, delimiter string) (string, error) {
	parts := make([]string, 0, len(values))
	for _, item := range values {
		if item == nil {
			parts = append(parts, "")
			continue
		}
		switch value := item.(type) {
		case string:
			parts = append(parts, value)
		case float64:
			parts = append(parts, strconv.FormatFloat(value, 'f', -1, 64))
		case float32:
			parts = append(parts, strconv.FormatFloat(float64(value), 'f', -1, 32))
		case int, int64, int32, uint, uint64, uint32:
			parts = append(parts, fmt.Sprint(value))
		case bool:
			parts = append(parts, strconv.FormatBool(value))
		default:
			return "", errors.New("join array contains non-scalar value")
		}
	}
	return strings.Join(parts, delimiter), nil
}

func coerceValue(value any, targetType LogicalType) (any, error) {
	switch targetType {
	case LogicalTypeString:
		return fmt.Sprint(value), nil
	case LogicalTypeNumber:
		switch typed := value.(type) {
		case float64:
			return typed, nil
		case float32:
			return float64(typed), nil
		case int:
			return float64(typed), nil
		case int64:
			return float64(typed), nil
		case int32:
			return float64(typed), nil
		case uint:
			return float64(typed), nil
		case uint64:
			return float64(typed), nil
		case uint32:
			return float64(typed), nil
		case bool:
			if typed {
				return float64(1), nil
			}
			return float64(0), nil
		case string:
			parsed, err := strconv.ParseFloat(strings.TrimSpace(typed), 64)
			if err != nil {
				return nil, errors.New("cannot coerce string to number")
			}
			return parsed, nil
		default:
			return nil, errors.New("cannot coerce value to number")
		}
	case LogicalTypeBoolean:
		switch typed := value.(type) {
		case bool:
			return typed, nil
		case string:
			parsed, err := strconv.ParseBool(strings.TrimSpace(typed))
			if err != nil {
				return nil, errors.New("cannot coerce string to boolean")
			}
			return parsed, nil
		case float64:
			return typed != 0, nil
		case float32:
			return typed != 0, nil
		case int:
			return typed != 0, nil
		case int64:
			return typed != 0, nil
		case int32:
			return typed != 0, nil
		case uint:
			return typed != 0, nil
		case uint64:
			return typed != 0, nil
		case uint32:
			return typed != 0, nil
		default:
			return nil, errors.New("cannot coerce value to boolean")
		}
	default:
		return nil, errors.New("unsupported target type for coercion")
	}
}

func addWarningOnce(warnings *[]Warning, warningSet map[string]struct{}, path, message string) {
	key := path + ":" + message
	if _, exists := warningSet[key]; exists {
		return
	}
	warningSet[key] = struct{}{}
	*warnings = append(*warnings, Warning{Path: path, Message: message})
	sort.SliceStable(*warnings, func(i, j int) bool {
		if (*warnings)[i].Path == (*warnings)[j].Path {
			return (*warnings)[i].Message < (*warnings)[j].Message
		}
		return (*warnings)[i].Path < (*warnings)[j].Path
	})
}

func deepCopyRecords(records []Record) []Record {
	copied := make([]Record, len(records))
	for i, record := range records {
		copied[i] = deepCopyRecord(record)
	}
	return copied
}

func deepCopyRecord(record Record) Record {
	copied := Record{}
	for key, value := range record {
		copied[key] = deepCopyValue(value)
	}
	return copied
}

func deepCopyValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		copied := map[string]any{}
		for key, nested := range typed {
			copied[key] = deepCopyValue(nested)
		}
		return copied
	case []any:
		copied := make([]any, len(typed))
		for index, item := range typed {
			copied[index] = deepCopyValue(item)
		}
		return copied
	default:
		return typed
	}
}
