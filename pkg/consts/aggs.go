package consts

type AggregationConfig struct {
	Field string
	Type  string
}

var AggregationConfigs = map[string]map[string]AggregationConfig{
	"BookStats": {
		"distinct_authors": {
			Field: "author_name.keyword",
			Type:  "cardinality",
		},
		"total_books": {
			Field: "_id",
			Type:  "value_count",
		},
	},
}
