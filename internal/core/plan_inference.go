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
	lossyDecisions := map[string]LossyDecision{}

	for _, record := range data.Values.Records {
		collectPlanSuggestions(record, "", flattenSet, explodeSet, joinRules, lossyDecisions)
	}

	flattenFields := setToSortedSlice(flattenSet)
	explodeArrays := setToSortedSlice(explodeSet)
	joinArrays := make([]JoinArrayRule, 0, len(joinRules))
	for _, rule := range joinRules {
		joinArrays = append(joinArrays, rule)
	}
	sort.Slice(joinArrays, func(i, j int) bool { return joinArrays[i].Path < joinArrays[j].Path })
	decisions := make([]LossyDecision, 0, len(lossyDecisions))
	for _, decision := range lossyDecisions {
		decisions = append(decisions, decision)
	}
	sort.Slice(decisions, func(i, j int) bool {
		if decisions[i].Strategy == decisions[j].Strategy {
			return decisions[i].FieldPath < decisions[j].FieldPath
		}
		return decisions[i].Strategy < decisions[j].Strategy
	})

	return ConversionPlan{
		FlattenFields:   flattenFields,
		ExplodeArrays:   explodeArrays,
		JoinArrays:      joinArrays,
		LossyDecisions:  decisions,
	}
}

func collectPlanSuggestions(value any, prefix string, flattenSet map[string]struct{}, explodeSet map[string]struct{}, joinRules map[string]JoinArrayRule, lossyDecisions map[string]LossyDecision) {
	if recordMap, ok := mapFromValue(value); ok {
		if prefix != "" {
			flattenSet[prefix] = struct{}{}
		}
		for key, nested := range recordMap {
			path := key
			if prefix != "" {
				path = prefix + "." + key
			}
			collectPlanSuggestions(nested, path, flattenSet, explodeSet, joinRules, lossyDecisions)
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
			lossyDecisions[prefix] = LossyDecision{
				FieldPath: prefix,
				Reason:    LossReasonFormatLimit,
				Strategy:  StrategyJoinArray,
			}
		default:
			joinRules[prefix] = JoinArrayRule{Path: prefix, Delimiter: ","}
			lossyDecisions[prefix] = LossyDecision{
				FieldPath: prefix,
				Reason:    LossReasonFormatLimit,
				Strategy:  StrategyJoinArray,
			}
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
