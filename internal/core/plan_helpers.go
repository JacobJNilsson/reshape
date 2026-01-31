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
	sort.Slice(plan.LossyDecisions, func(i, j int) bool {
		if plan.LossyDecisions[i].Strategy == plan.LossyDecisions[j].Strategy {
			return plan.LossyDecisions[i].FieldPath < plan.LossyDecisions[j].FieldPath
		}
		return plan.LossyDecisions[i].Strategy < plan.LossyDecisions[j].Strategy
	})
	return plan
}

// ValidateLossyDecisions ensures lossy actions are explicitly acknowledged.
func ValidateLossyDecisions(plan ConversionPlan) error {
	lossyMap := map[string]Strategy{}
	for _, decision := range plan.LossyDecisions {
		key := string(decision.Strategy) + ":" + decision.FieldPath
		lossyMap[key] = decision.Strategy
	}
	for _, rule := range plan.JoinArrays {
		key := string(StrategyJoinArray) + ":" + rule.Path
		if _, ok := lossyMap[key]; !ok {
			return errors.New("join_arrays requires lossy_decisions entry for path: " + rule.Path)
		}
	}
	for _, rule := range plan.TypeCoercions {
		key := string(StrategyCoerceType) + ":" + rule.Path
		if _, ok := lossyMap[key]; !ok {
			return errors.New("type_coercions requires lossy_decisions entry for path: " + rule.Path)
		}
	}
	for _, path := range plan.DropFields {
		key := string(StrategyDropField) + ":" + path
		if _, ok := lossyMap[key]; !ok {
			return errors.New("drop_fields requires lossy_decisions entry for path: " + path)
		}
	}
	return nil
}
