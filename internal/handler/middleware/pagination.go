package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type QueryParams struct {
	Filter     bson.M
	Projection bson.M
	Sort       bson.D
	Page       int64
	Limit      int64
}

type contextKey string

const QueryParamsKey = contextKey("queryParams")

func ParseAdvancedQuery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		filter := bson.M{}
		projection := bson.M{}
		sort := bson.D{}
		page := int64(1)
		limit := int64(10)

		// 1. Filtering: Support rating[between], branch[in], catalog[and]
		for key, vals := range q {
			value := vals[0]
			if strings.Contains(key, "[") && strings.Contains(key, "]") {
				// Parse key with operator, e.g. rating[between]
				field := key[:strings.Index(key, "[")]
				op := key[strings.Index(key, "[")+1 : strings.Index(key, "]")]

				switch op {
				case "between":
					// Expect value like "3,5"
					parts := strings.Split(value, ",")
					if len(parts) == 2 {
						min, errMin := strconv.ParseFloat(parts[0], 64)
						max, errMax := strconv.ParseFloat(parts[1], 64)
						if errMin == nil && errMax == nil {
							filter[field] = bson.M{
								"$gte": min,
								"$lte": max,
							}
						}
					}
				case "in":
					// e.g. branch[in]=tech,gadget
					inVals := strings.Split(value, ",")
					filter[field] = bson.M{"$in": inVals}
				case "or":
					// e.g. branch[in]=tech,gadget
					inVals := strings.Split(value, ",")
					filter[field] = bson.M{"$or": inVals}
				case "and":
					// e.g. catalog[and]=ifb,cb  -> $all
					allVals := strings.Split(value, ",")
					filter[field] = bson.M{"$all": allVals}
				}
			}
		}

		// 2. Projection: ?fields=name,category
		if fields := q.Get("fields"); fields != "" {
			for _, f := range strings.Split(fields, ",") {
				projection[f] = 1
			}
		}

		// 3. Full-text search: ?search=ad_name
		if search := q.Get("search"); search != "" {
			// Basic $text search, assumes text index on relevant fields exists
			filter["$text"] = bson.M{"$search": search}
		}

		// 4. Sorting: ?sort=-name,creation_date
		if sortStr := q.Get("sort"); sortStr != "" {
			for _, field := range strings.Split(sortStr, ",") {
				order := int32(1)
				if strings.HasPrefix(field, "-") {
					order = -1
					field = field[1:]
				}
				sort = append(sort, bson.E{Key: field, Value: order})
			}
		}

		// 5. Pagination: ?page=2&limit=10
		if p := q.Get("page"); p != "" {
			if parsed, err := strconv.ParseInt(p, 10, 64); err == nil && parsed > 0 {
				page = parsed
			}
		}

		if l := q.Get("limit"); l != "" {
			if parsed, err := strconv.ParseInt(l, 10, 64); err == nil && parsed > 0 {
				limit = parsed
			}
		}

		// Save in context for next handlers
		qp := QueryParams{
			Filter:     filter,
			Projection: projection,
			Sort:       sort,
			Page:       page,
			Limit:      limit,
		}

		ctx := context.WithValue(r.Context(), QueryParamsKey, qp)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
