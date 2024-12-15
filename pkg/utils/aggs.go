package utils

import (
	"book_service/pkg/constants"
	"fmt"
)

type AggregationResult map[string]interface{}

func ParseAggregations(aggregations map[string]interface{}, aggGroup map[string]constants.AggregationConfig) (AggregationResult, error) {
	result := make(AggregationResult)

	if aggregations == nil {
		return result, nil
	}

	for aggName, _ := range aggGroup {
		if aggData, ok := aggregations[aggName].(map[string]interface{}); ok {
			if value, ok := aggData["value"].(float64); ok {
				result[aggName] = int(value)
			} else {
				return nil, fmt.Errorf("invalid format for aggregation %s", aggName)
			}
		}
	}

	return result, nil
}
