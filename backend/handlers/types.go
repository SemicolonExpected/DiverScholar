package handlers

// Used to store individual papers
type Paper struct {
	URL   string
	Title string
}

// Used to track a search
type Search struct {
	Hash  int64
	URL   string
	Pages []string
}

// Used to cache full authors
type Author struct {
	AuthorLink string // Key
	FullName   string
	FirstName  string
	Metadata   string
}

// Used to cache gendered name results
type Name struct {
	FirstName  string // Key
	Gender     string
	Confidence float64
}
