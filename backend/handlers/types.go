package handlers

type User struct {
	Key string `json:"key,omitempty" ,datastore:"key"`
}

type Paper struct {
	Key     int64    `json:"key,omitempty" ,datastore:"key"`
	URL     string   `json:"url,omitempty" ,datastore:"url,noindex"`
	Title   string   `json:"title,omitempty" ,datastore:"title,noindex"`
	Authors []Author `json:"authors,omitempty" ,datastore:"authors,noindex"`
}

type Search struct {
	Key    int64   `json:"key,omitempty" ,datastore:"key"`
	URL    string  `json:"url,omitempty" ,datastore:"url,noindex"`
	Papers []Paper `json:"pages,omitempty" ,datastore:"papers,noindex"`
	Ranked []int   `json:"ranked,omitempty" ,datastore:"ranked,noindex"`
}

type Author struct {
	Key        int64  `json:"key,omitempty" ,datastore:"key"`
	AuthorLink string `json:"author_link,omitempty" ,datastore:"author_link,noindex"`
	FullName   string `json:"full_name,omitempty" ,datastore:"full_name,noindex"`
	FirstName  string `json:"first_name,omitempty" ,datastore:"first_name,noindex"`
	Score      Name   `json:"score" ,datastore:"score,noindex"`
}

type Name struct {
	FirstName  string  `json:"name,omitempty" ,datastore:"key"`
	Gender     string  `json:"gender,omitempty" ,datastore:"gender,noindex"`
	Confidence float64 `json:"probability,omitempty" ,datastore:"conf,noindex"`
}
