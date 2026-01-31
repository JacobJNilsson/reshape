package formats

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"sort"

	"reshape/internal/core"
)

// ParseCSV converts CSV bytes into canonical data.
func ParseCSV(input []byte) (core.CanonicalData, error) {
	reader := csv.NewReader(bytes.NewReader(input))
	rows, err := reader.ReadAll()
	if err != nil {
		return core.CanonicalData{}, err
	}
	if len(rows) == 0 {
		return core.CanonicalData{}, errors.New("csv input is empty")
	}
	headers := rows[0]
	records := make([]core.Record, 0, len(rows)-1)
	for _, row := range rows[1:] {
		if len(row) != len(headers) {
			return core.CanonicalData{}, errors.New("csv row has different column count than headers")
		}
		record := core.Record{}
		for index, header := range headers {
			value := row[index]
			if value == "" {
				record[header] = nil
				continue
			}
			record[header] = value
		}
		records = append(records, record)
	}

	schema := core.BuildSchemaFromRecords(records)
	return core.CanonicalData{Schema: schema, Records: records}, nil
}

// RenderCSV converts canonical data into CSV bytes.
func RenderCSV(data core.CanonicalData) ([]byte, error) {
	headers := schemaHeaders(data)

	buffer := &bytes.Buffer{}
	writer := csv.NewWriter(buffer)
	if err := writer.Write(headers); err != nil {
		return nil, err
	}
	for _, record := range data.Records {
		row := make([]string, len(headers))
		for index, header := range headers {
			value, exists, err := core.ValueAtPath(record, header)
			if err != nil {
				return nil, err
			}
			if !exists || value == nil {
				row[index] = ""
				continue
			}
			if _, ok := value.(map[string]any); ok {
				return nil, fmt.Errorf("csv output requires scalar at path %s", header)
			}
			if _, ok := value.([]any); ok {
				return nil, fmt.Errorf("csv output requires scalar at path %s", header)
			}
			scalar, err := formatScalar(value)
			if err != nil {
				return nil, err
			}
			row[index] = scalar
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func schemaHeaders(data core.CanonicalData) []string {
	headers := []string{}
	for _, field := range data.Schema.Fields {
		headers = append(headers, field.Path)
	}
	if len(headers) == 0 {
		headers = recordHeaders(data.Records)
	}
	sort.Strings(headers)
	return headers
}

func recordHeaders(records []core.Record) []string {
	set := map[string]struct{}{}
	for _, record := range records {
		for key := range record {
			set[key] = struct{}{}
		}
	}
	result := make([]string, 0, len(set))
	for key := range set {
		result = append(result, key)
	}
	return result
}

func formatScalar(value any) (string, error) {
	switch typed := value.(type) {
	case string:
		return typed, nil
	case float64, float32, int, int64, int32, uint, uint64, uint32, bool:
		return fmt.Sprint(typed), nil
	default:
		return "", errors.New("csv output contains unsupported value type")
	}
}
