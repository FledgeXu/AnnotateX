package model

type UserFilter struct {
	Keyword  string // fuzzy match on username or email
	IsActive string // "true" / "false" / "" (no filtering)
	SortBy   string // column to sort by
	Order    string // "asc" or "desc"
	Limit    int    // number of results to return
	Offset   int    // number of records to skip
}
