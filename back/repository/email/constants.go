package email

const (
	SearchTypeMatchAll    = "matchall"
	SearchTypeMatch       = "match"
	SearchTypeMatchPhrase = "matchphrase"
	SearchTypeTerm        = "term"
	SearchTypeQueryString = "querystring"
	SearchTypePrefix      = "prefix"
	SearchTypeWildcard    = "wildcard"
	SearchTypeFuzzy       = "fuzzy"
	SearchTypeDateRange   = "daterange"
)

// Default values
const (
	DefaultMaxResults = 20
	DefaultFrom       = 0
)
