package types

// NewQueryPaginationParams creates a new instance of QueryPaginationParams.
func NewQueryPaginationParams(page uint64, limit uint64) QueryPaginationParams {
	return QueryPaginationParams{Page: page, Limit: limit}
}
