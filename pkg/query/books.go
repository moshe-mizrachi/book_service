package query

import "book_service/pkg/constants"

type Builder struct {
	id         *string
	title      *string
	authorName *string
	priceMin   *float64
	priceMax   *float64
}

func NewQueryBuilder() *Builder {
	return &Builder{}
}

func (qb *Builder) ID(i string) *Builder {
	if i == "" {
		return qb
	}

	qb.id = &i
	return qb
}

func (qb *Builder) Title(t string) *Builder {
	if t == "" {
		return qb
	}

	qb.title = &t
	return qb
}

func (qb *Builder) AuthorName(a string) *Builder {
	if a == "" {
		return qb
	}

	qb.authorName = &a
	return qb
}

func (qb *Builder) PriceRange(min, max float64) *Builder {
	unseted := min == max && min == 0
	if unseted {
		qb.priceMin = &min
		highestPrice := constants.HighestBookPrice
		qb.priceMax = &highestPrice
		return qb
	}

	qb.priceMin = &min
	qb.priceMax = &max
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

	if qb.priceMin != nil && qb.priceMax != nil {
		mustClauses = append(mustClauses, map[string]interface{}{
			"range": map[string]interface{}{
				"price": map[string]interface{}{
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
