package utils

import (
	"book_service/pkg/consts"
	"fmt"

	"github.com/samber/lo"
)

type AggregationResult map[string]interface{}

func ParseAggregations(aggregations map[string]interface{}, aggGroup map[string]consts.AggregationConfig) (AggregationResult, error) {
	result := make(AggregationResult)

	if aggregations == nil {
		return result, nil
	}

	result = lo.MapValues(aggGroup, func(_ consts.AggregationConfig, aggName string) interface{} {
		if aggData, ok := aggregations[aggName].(map[string]interface{}); ok {
			if value, ok := aggData["value"].(float64); ok {
				return int(value)
			}
		}
		return nil
	})

	for aggName, value := range result {
		if value == nil {
			return nil, fmt.Errorf("invalid format for aggregation %s", aggName)
		}
	}

	return result, nil
}
