package training

type Country struct {
	Name     map[string]string `json:"name"`
	Capital  string            `json:"capital"`
	Code     string            `json:"code"`
	Area     float64           `json:"area"`
	Currency string            `json:"currency"`
}

type Joke struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

type Message struct {
	Tag      string   `json:"tag"`
	Messages []string `json:"messages"`
}

type Information struct {
	Name           string   `json:"name"`
	MovieGenres    []string `json:"movie_genres"`
	MovieBlacklist []string `json:"movie_blacklist"`
}

type Movie struct {
	Name   string
	Genres []string
	Rating float64
}

type Modulem struct {
	Tag       string
	Patterns  []string
	Responses []string
	Replacer  func(string, string, string, string) (string, string)
	Context   string
}

type Intent struct {
	Tag       string   `json:"tag"`
	Patterns  []string `json:"patterns"`
	Responses []string `json:"responses"`
	Context   string   `json:"context"`
}

type Sentence struct {
	Locale  string
	Content string
}

type Document struct {
	Sentence Sentence
	Tag      string
}

type Locale struct {
	Tag  string
	Name string
}
