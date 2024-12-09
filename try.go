package main

import (
	"encoding/json"
	"fmt"
)

type QueryBuilder struct {
	title      *string
	authorName *string
	priceMin   *float64
	priceMax   *float64
	username   *string
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

func (qb *QueryBuilder) Title(t string) *QueryBuilder {
	qb.title = &t
	return qb
}

func (qb *QueryBuilder) AuthorName(a string) *QueryBuilder {
	qb.authorName = &a
	return qb
}

func (qb *QueryBuilder) PriceRange(min, max float64) *QueryBuilder {
	qb.priceMin = &min
	qb.priceMax = &max
	return qb
}

func (qb *QueryBuilder) Username(u string) *QueryBuilder {
	qb.username = &u
	return qb
}

func (qb *QueryBuilder) Build() map[string]interface{} {
	var mustClauses []map[string]interface{}

	if qb.title != nil {
		mustClauses = append(mustClauses, map[string]interface{}{
			"match": map[string]interface{}{
				"title": *qb.title,
			},
		})
	}

	if qb.authorName != nil {
		mustClauses = append(mustClauses, map[string]interface{}{
			"match": map[string]interface{}{
				"author_name": *qb.authorName,
			},
		})
	}

	if qb.username != nil {
		mustClauses = append(mustClauses, map[string]interface{}{
			"match": map[string]interface{}{
				"username": *qb.username,
			},
		})
	}

	if qb.priceMin != nil && qb.priceMax != nil {
		mustClauses = append(mustClauses, map[string]interface{}{
			"range": map[string]interface{}{
				"price_range": map[string]interface{}{
					"gte": *qb.priceMin,
					"lte": *qb.priceMax,
				},
			},
		})
	}

	if len(mustClauses) > 0 {
		return map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": mustClauses,
				},
			},
		}
	}

	return map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
}

func main() {
	// Example 1: Title only
	q1 := NewQueryBuilder().Title("some bullshit").Build()
	data1, _ := json.MarshalIndent(q1, "", "  ")
	fmt.Println("Example 1:")
	fmt.Println(string(data1))

	// Example 2: Title and PriceRange
	q2 := NewQueryBuilder().Title("some bullshit").PriceRange(10, 20).Build()
	data2, _ := json.MarshalIndent(q2, "", "  ")
	fmt.Println("\nExample 2:")
	fmt.Println(string(data2))

	// Example 3: AuthorName only
	q3 := NewQueryBuilder().AuthorName("John Doe").Build()
	data3, _ := json.MarshalIndent(q3, "", "  ")
	fmt.Println("\nExample 3:")
	fmt.Println(string(data3))

	// Example 4: Username only
	q4 := NewQueryBuilder().Username("johnd").Build()
	data4, _ := json.MarshalIndent(q4, "", "  ")
	fmt.Println("\nExample 4:")
	fmt.Println(string(data4))

	// Example 5: Title, AuthorName, and Username
	q5 := NewQueryBuilder().Title("another title").AuthorName("Jane Smith").Username("janes").Build()
	data5, _ := json.MarshalIndent(q5, "", "  ")
	fmt.Println("\nExample 5:")
	fmt.Println(string(data5))

	// Example 6: Title, PriceRange, Username
	q6 := NewQueryBuilder().Title("search query").PriceRange(5, 15).Username("user123").Build()
	data6, _ := json.MarshalIndent(q6, "", "  ")
	fmt.Println("\nExample 6:")
	fmt.Println(string(data6))

	// Example 7: No fields (match_all)
	q7 := NewQueryBuilder().Build()
	data7, _ := json.MarshalIndent(q7, "", "  ")
	fmt.Println("\nExample 7:")
	fmt.Println(string(data7))
}
