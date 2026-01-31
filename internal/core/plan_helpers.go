package core

import (
	"errors"
	"sort"
)

// NormalizePlan sorts plan slices for deterministic application.
func NormalizePlan(plan ConversionPlan) ConversionPlan {
	sort.Strings(plan.FlattenFields)
	sort.Strings(plan.ExplodeArrays)
	sort.Strings(plan.DropFields)
	sort.Slice(plan.JoinArrays, func(i, j int) bool { return plan.JoinArrays[i].Path < plan.JoinArrays[j].Path })
	sort.Slice(plan.TypeCoercions, func(i, j int) bool { return plan.TypeCoercions[i].Path < plan.TypeCoercions[j].Path })
	sort.Slice(plan.DefaultValues, func(i, j int) bool { return plan.DefaultValues[i].Path < plan.DefaultValues[j].Path })
	sort.Slice(plan.LossyOperations, func(i, j int) bool {
		if plan.LossyOperations[i].Operation == plan.LossyOperations[j].Operation {
			return plan.LossyOperations[i].Path < plan.LossyOperations[j].Path
		}
		return plan.LossyOperations[i].Operation < plan.LossyOperations[j].Operation
	})
	return plan
}

// ValidateLossyOperations ensures lossy actions are explicitly acknowledged.
func ValidateLossyOperations(plan ConversionPlan) error {
	lossyMap := map[string]LossyOperationType{}
	for _, op := range plan.LossyOperations {
		key := string(op.Operation) + ":" + op.Path
		lossyMap[key] = op.Operation
	}
	for _, rule := range plan.JoinArrays {
		key := string(LossyOperationJoinArray) + ":" + rule.Path
		if _, ok := lossyMap[key]; !ok {
			return errors.New("join_arrays requires lossy_operations entry for path: " + rule.Path)
		}
	}
	for _, rule := range plan.TypeCoercions {
		key := string(LossyOperationCoerceType) + ":" + rule.Path
		if _, ok := lossyMap[key]; !ok {
			return errors.New("type_coercions requires lossy_operations entry for path: " + rule.Path)
		}
	}
	for _, path := range plan.DropFields {
		key := string(LossyOperationDropField) + ":" + path
		if _, ok := lossyMap[key]; !ok {
			return errors.New("drop_fields requires lossy_operations entry for path: " + path)
		}
	}
	return nil
}
