package formats

import (
	"encoding/json"
	"errors"

	"reshape/internal/core"
)

// ParseJSON converts JSON bytes into canonical data.
func ParseJSON(input []byte) (core.CanonicalData, error) {
	var decoded any
	if err := json.Unmarshal(input, &decoded); err != nil {
		return core.CanonicalData{}, err
	}
	records, err := jsonToRecords(decoded)
	if err != nil {
		return core.CanonicalData{}, err
	}
	shape := core.BuildShapeFromRecords(records)
	return core.CanonicalData{
		Shape:  shape,
		Values: core.DataValues{Records: records},
	}, nil
}

// RenderJSON converts canonical data into JSON bytes.
func RenderJSON(data core.CanonicalData) ([]byte, error) {
	if len(data.Values.Records) == 0 {
		return json.Marshal([]any{})
	}
	if len(data.Values.Records) == 1 {
		return json.Marshal(data.Values.Records[0])
	}
	return json.Marshal(data.Values.Records)
}

func jsonToRecords(value any) ([]core.Record, error) {
	switch typed := value.(type) {
	case map[string]any:
		return []core.Record{core.Record(typed)}, nil
	case []any:
		records := make([]core.Record, 0, len(typed))
		for _, item := range typed {
			mapValue, ok := item.(map[string]any)
			if !ok {
				return nil, errors.New("json array contains non-object value")
			}
			records = append(records, core.Record(mapValue))
		}
		return records, nil
	default:
		return nil, errors.New("json input must be an object or array of objects")
	}
}
