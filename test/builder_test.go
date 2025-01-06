package test

import (
	"book_service/pkg/query"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryBuilder_ID(t *testing.T) {
	qb := query.NewQueryBuilder().ID("12345")
	result := qb.Build()

	expected := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{"term": map[string]interface{}{"_id": "12345"}},
				},
			},
		},
	}

	assert.Equal(t, expected, result)
}

func TestQueryBuilder_Title(t *testing.T) {
	qb := query.NewQueryBuilder().Title("test title")
	result := qb.Build()

	expected := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{"match": map[string]interface{}{"title": "test title"}},
				},
			},
		},
	}

	assert.Equal(t, expected, result)
}

func TestQueryBuilder_AuthorName(t *testing.T) {
	qb := query.NewQueryBuilder().AuthorName("John Doe")
	result := qb.Build()

	expected := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{"match": map[string]interface{}{"author_name": "John Doe"}},
				},
			},
		},
	}

	assert.Equal(t, expected, result)
}

func TestQueryBuilder_PriceRange(t *testing.T) {
	qb := query.NewQueryBuilder().PriceRange(10, 20)
	result := qb.Build()

	expected := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"range": map[string]interface{}{
							"price_range": map[string]interface{}{
								"gte": float64(10),
								"lte": float64(20),
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expected, result)
}

func TestQueryBuilder_MatchAll(t *testing.T) {
	qb := query.NewQueryBuilder()
	result := qb.Build()

	expected := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	assert.Equal(t, expected, result)
}

func TestQueryBuilder_Combination(t *testing.T) {
	qb := query.NewQueryBuilder().ID("12345").Title("test title").AuthorName("John Doe")
	result := qb.Build()

	expected := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{"term": map[string]interface{}{"_id": "12345"}},
					{"match": map[string]interface{}{"title": "test title"}},
					{"match": map[string]interface{}{"author_name": "John Doe"}},
				},
			},
		},
	}

	assert.Equal(t, expected, result)
}

func TestQueryBuilder_ComplexCombination(t *testing.T) {
	qb := query.NewQueryBuilder().
		ID("12345").
		Title("test title").
		PriceRange(50, 100).
		AuthorName("Author")

	result := qb.Build()

	expected := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{"term": map[string]interface{}{"_id": "12345"}},
					{"match": map[string]interface{}{"title": "test title"}},
					{"range": map[string]interface{}{
						"price_range": map[string]interface{}{
							"gte": float64(50),
							"lte": float64(100),
						},
					}},
					{"match": map[string]interface{}{"author_name": "Author"}},
				},
			},
		},
	}

	actualMust := result["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"]
	expectedMust := expected["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"]

	assert.ElementsMatch(t, expectedMust, actualMust, "The 'must' array does not match")
}
