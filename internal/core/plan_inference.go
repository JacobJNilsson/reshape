package core

import (
	"sort"
)

// InferConversionPlan suggests a plan for the target format.
func InferConversionPlan(data CanonicalData, targetFormat string) ConversionPlan {
	if targetFormat != "csv" {
		return ConversionPlan{}
	}
	flattenSet := map[string]struct{}{}
	explodeSet := map[string]struct{}{}
	joinRules := map[string]JoinArrayRule{}
	lossyOps := map[string]LossyOperation{}

	for _, record := range data.Values.Records {
		collectPlanSuggestions(record, "", flattenSet, explodeSet, joinRules, lossyOps)
	}

	flattenFields := setToSortedSlice(flattenSet)
	explodeArrays := setToSortedSlice(explodeSet)
	joinArrays := make([]JoinArrayRule, 0, len(joinRules))
	for _, rule := range joinRules {
		joinArrays = append(joinArrays, rule)
	}
	sort.Slice(joinArrays, func(i, j int) bool { return joinArrays[i].Path < joinArrays[j].Path })
	lossyOperations := make([]LossyOperation, 0, len(lossyOps))
	for _, op := range lossyOps {
		lossyOperations = append(lossyOperations, op)
	}
	sort.Slice(lossyOperations, func(i, j int) bool {
		if lossyOperations[i].Operation == lossyOperations[j].Operation {
			return lossyOperations[i].Path < lossyOperations[j].Path
		}
		return lossyOperations[i].Operation < lossyOperations[j].Operation
	})

	return ConversionPlan{
		FlattenFields:   flattenFields,
		ExplodeArrays:   explodeArrays,
		JoinArrays:      joinArrays,
		LossyOperations: lossyOperations,
	}
}

func collectPlanSuggestions(value any, prefix string, flattenSet map[string]struct{}, explodeSet map[string]struct{}, joinRules map[string]JoinArrayRule, lossyOps map[string]LossyOperation) {
	if recordMap, ok := mapFromValue(value); ok {
		if prefix != "" {
			flattenSet[prefix] = struct{}{}
		}
		for key, nested := range recordMap {
			path := key
			if prefix != "" {
				path = prefix + "." + key
			}
			collectPlanSuggestions(nested, path, flattenSet, explodeSet, joinRules, lossyOps)
		}
		return
	}
	if sliceValue, ok := value.([]any); ok {
		if prefix == "" {
			return
		}
		arrayType := inferArrayType(sliceValue)
		switch arrayType {
		case LogicalTypeObject:
			explodeSet[prefix] = struct{}{}
		case LogicalTypeString, LogicalTypeNumber, LogicalTypeBoolean:
			joinRules[prefix] = JoinArrayRule{Path: prefix, Delimiter: ","}
			lossyOps[prefix] = LossyOperation{Path: prefix, Operation: LossyOperationJoinArray, Reason: "CSV does not support arrays"}
		default:
			joinRules[prefix] = JoinArrayRule{Path: prefix, Delimiter: ","}
			lossyOps[prefix] = LossyOperation{Path: prefix, Operation: LossyOperationJoinArray, Reason: "CSV does not support arrays"}
		}
		return
	}
}

func inferArrayType(values []any) LogicalType {
	for _, item := range values {
		if item == nil {
			continue
		}
		switch item.(type) {
		case map[string]any:
			return LogicalTypeObject
		case []any:
			return LogicalTypeArray
		case string:
			return LogicalTypeString
		case float64, float32, int, int64, int32, uint, uint64, uint32:
			return LogicalTypeNumber
		case bool:
			return LogicalTypeBoolean
		}
	}
	return LogicalTypeArray
}

func setToSortedSlice(items map[string]struct{}) []string {
	result := make([]string, 0, len(items))
	for item := range items {
		result = append(result, item)
	}
	sort.Strings(result)
	return result
}
