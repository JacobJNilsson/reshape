package core

import (
	"sort"
)

type fieldStats struct {
	presentCount int
	nullCount    int
	typeCounts   map[LogicalType]int
	repeated     bool
}

// BuildShapeFromRecords infers a shape from canonical records.
func BuildShapeFromRecords(records []Record) DataShape {
	statsByPath := map[string]*fieldStats{}
	for _, record := range records {
		pathsInRecord := map[string]struct{}{}
		collectFieldStats(record, "", statsByPath, pathsInRecord)
		for path, stats := range statsByPath {
			if _, ok := pathsInRecord[path]; ok {
				stats.presentCount++
			}
		}
	}

	fields := make([]FieldDefinition, 0, len(statsByPath))
	for path, stats := range statsByPath {
		fieldType := chooseLogicalType(stats.typeCounts)
		nullable := stats.nullCount > 0 || stats.presentCount < len(records)
		fields = append(fields, FieldDefinition{
			Path:     path,
			Type:     fieldType,
			Nullable: nullable,
			Repeated: stats.repeated,
		})
	}

	sort.Slice(fields, func(i, j int) bool { return fields[i].Path < fields[j].Path })
	return DataShape{Fields: fields}
}

func collectFieldStats(value any, prefix string, statsByPath map[string]*fieldStats, pathsInRecord map[string]struct{}) {
	recordMap, ok := mapFromValue(value)
	if ok {
		if prefix != "" {
			ensureFieldStats(prefix, statsByPath).typeCounts[LogicalTypeObject]++
			pathsInRecord[prefix] = struct{}{}
		}
		for key, nested := range recordMap {
			path := key
			if prefix != "" {
				path = prefix + "." + key
			}
			collectFieldStats(nested, path, statsByPath, pathsInRecord)
		}
		return
	}

	sliceValue, ok := value.([]any)
	if ok {
		stats := ensureFieldStats(prefix, statsByPath)
		stats.repeated = true
		pathsInRecord[prefix] = struct{}{}
		if len(sliceValue) == 0 {
			stats.typeCounts[LogicalTypeArray]++
			return
		}
		for _, item := range sliceValue {
			subPath := prefix
			collectArrayItemStats(item, subPath, statsByPath, pathsInRecord)
		}
		return
	}

	stats := ensureFieldStats(prefix, statsByPath)
	pathsInRecord[prefix] = struct{}{}
	if value == nil {
		stats.nullCount++
		return
	}
	switch value.(type) {
	case string:
		stats.typeCounts[LogicalTypeString]++
	case float64, float32, int, int64, int32, uint, uint64, uint32:
		stats.typeCounts[LogicalTypeNumber]++
	case bool:
		stats.typeCounts[LogicalTypeBoolean]++
	default:
		stats.typeCounts[LogicalTypeString]++
	}
}

func collectArrayItemStats(item any, prefix string, statsByPath map[string]*fieldStats, pathsInRecord map[string]struct{}) {
	if item == nil {
		stats := ensureFieldStats(prefix, statsByPath)
		stats.nullCount++
		pathsInRecord[prefix] = struct{}{}
		return
	}
	if recordMap, ok := mapFromValue(item); ok {
		stats := ensureFieldStats(prefix, statsByPath)
		stats.typeCounts[LogicalTypeObject]++
		pathsInRecord[prefix] = struct{}{}
		for key, nested := range recordMap {
			path := prefix + "." + key
			collectFieldStats(nested, path, statsByPath, pathsInRecord)
		}
		return
	}
	if nestedArray, ok := item.([]any); ok {
		stats := ensureFieldStats(prefix, statsByPath)
		stats.typeCounts[LogicalTypeArray]++
		pathsInRecord[prefix] = struct{}{}
		for _, nested := range nestedArray {
			collectArrayItemStats(nested, prefix, statsByPath, pathsInRecord)
		}
		return
	}
	stats := ensureFieldStats(prefix, statsByPath)
	pathsInRecord[prefix] = struct{}{}
	switch item.(type) {
	case string:
		stats.typeCounts[LogicalTypeString]++
	case float64, float32, int, int64, int32, uint, uint64, uint32:
		stats.typeCounts[LogicalTypeNumber]++
	case bool:
		stats.typeCounts[LogicalTypeBoolean]++
	default:
		stats.typeCounts[LogicalTypeString]++
	}
}

func ensureFieldStats(path string, statsByPath map[string]*fieldStats) *fieldStats {
	stats, ok := statsByPath[path]
	if !ok {
		stats = &fieldStats{typeCounts: map[LogicalType]int{}}
		statsByPath[path] = stats
	}
	return stats
}

func chooseLogicalType(typeCounts map[LogicalType]int) LogicalType {
	if len(typeCounts) == 0 {
		return LogicalTypeString
	}
	if len(typeCounts) == 1 {
		for typeValue := range typeCounts {
			return typeValue
		}
	}
	priority := []LogicalType{LogicalTypeObject, LogicalTypeArray, LogicalTypeString, LogicalTypeNumber, LogicalTypeBoolean}
	for _, typeValue := range priority {
		if typeCounts[typeValue] > 0 {
			return typeValue
		}
	}
	return LogicalTypeString
}
