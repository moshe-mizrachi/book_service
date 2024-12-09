package query

import (
	"encoding/json"
	"fmt"
)

type Builder struct {
	id         *string
	title      *string
	authorName *string
	priceMin   *float64
	priceMax   *float64
	username   *string
}

func NewQueryBuilder() *Builder {
	return &Builder{}
}

func (qb *Builder) ID(i string) *Builder {
	qb.id = &i
	return qb
}

func (qb *Builder) Title(t string) *Builder {
	qb.title = &t
	return qb
}

func (qb *Builder) AuthorName(a string) *Builder {
	qb.authorName = &a
	return qb
}

func (qb *Builder) PriceRange(min, max float64) *Builder {
	qb.priceMin = &min
	qb.priceMax = &max
	return qb
}

func (qb *Builder) Username(u string) *Builder {
	qb.username = &u
	return qb
}

func (qb *Builder) Build() map[string]interface{} {
	var mustClauses []map[string]interface{}

	if qb.id != nil {
		mustClauses = append(mustClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"_id": *qb.id,
			},
		})
	}

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
	// Here are some examples that now include the id field:

	// Example 1: ID only (exact match)
	q1 := NewQueryBuilder().ID("12345").Build()
	data1, _ := json.MarshalIndent(q1, "", "  ")
	fmt.Println("Example 1:")
	fmt.Println(string(data1))

	// Example 2: ID and Title
	q2 := NewQueryBuilder().ID("12345").Title("some bullshit").Build()
	data2, _ := json.MarshalIndent(q2, "", "  ")
	fmt.Println("\nExample 2:")
	fmt.Println(string(data2))

	// Example 3: ID, Title, and PriceRange
	q3 := NewQueryBuilder().ID("abc123").Title("another title").PriceRange(10, 20).Build()
	data3, _ := json.MarshalIndent(q3, "", "  ")
	fmt.Println("\nExample 3:")
	fmt.Println(string(data3))

	// Example 4: Title and Username
	q4 := NewQueryBuilder().Title("some bullshit").Username("johnd").Build()
	data4, _ := json.MarshalIndent(q4, "", "  ")
	fmt.Println("\nExample 4:")
	fmt.Println(string(data4))

	// Example 5: AuthorName only
	q5 := NewQueryBuilder().AuthorName("Jane Smith").Build()
	data5, _ := json.MarshalIndent(q5, "", "  ")
	fmt.Println("\nExample 5:")
	fmt.Println(string(data5))

	// Example 6: ID, AuthorName, Username
	q6 := NewQueryBuilder().ID("xyz789").AuthorName("Jane Smith").Username("janes").Build()
	data6, _ := json.MarshalIndent(q6, "", "  ")
	fmt.Println("\nExample 6:")
	fmt.Println(string(data6))

	// Example 7: No fields (match_all)
	q7 := NewQueryBuilder().Build()
	data7, _ := json.MarshalIndent(q7, "", "  ")
	fmt.Println("\nExample 7:")
	fmt.Println(string(data7))
}
