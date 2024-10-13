package training

import (
	"time"
)

var (
	CapitalTag  = "capital"
	AreaTag     = "area"
	JokesTag    = "jokes"
	CurrencyTag = "currency"
	countries   = SerializeCountries()
	messages    = map[string][]Message{}

	GenresTag = "movies genres"

	MoviesTag = "movies search"

	MoviesAlreadyTag = "already seen movie"

	MoviesDataTag   = "movies search from data"
	userInformation = map[string]Information{}

	MoviesGenres = map[string][]string{
		"en": {
			"Action", "Adventure", "Animation", "Children", "Comedy", "Crime", "Documentary", "Drama", "Fantasy",
			"Film-Noir", "Horror", "Musical", "Mystery", "Romance", "Sci-Fi", "Thriller", "War", "Western",
		},
	}
	movies = SerializeMovies()

	ArticleCountriesm = map[string]func(string) string{}

	rules []Rule

	modulesm = map[string][]Modulem{}

	intents = map[string][]Intent{}

	Locales = []Locale{
		{
			Tag:  "en",
			Name: "english",
		},
	}

	MathTag = "math"

	NameGetterTag = "name getter"

	NameSetterTag = "name setter"

	MathDecimals = map[string]string{
		"en": `(\d+( |-)decimal(s)?)|(number (of )?decimal(s)? (is )?\d+)`,
	}

	names      = SerializeNames()
	decimal    = "\\b\\d+([\\.,]\\d+)?"
	RandomTag  = "random number"
	AdvicesTag = "advices"

	daysOfWeek = map[string]time.Weekday{
		"monday":    time.Monday,
		"tuesday":   time.Tuesday,
		"wednesday": time.Wednesday,
		"thursday":  time.Thursday,
		"friday":    time.Friday,
		"saturday":  time.Saturday,
		"sunday":    time.Sunday,
	}

	RuleTranslations = map[string]RuleTranslation{
		"en": {
			DaysOfWeek: []string{
				"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday",
			},
			Months: []string{
				"january", "february", "march", "april", "may", "june", "july",
				"august", "september", "october", "november", "december",
			},
			RuleToday:         `today|tonight`,
			RuleTomorrow:      `(after )?tomorrow`,
			RuleAfterTomorrow: "after",
			RuleDayOfWeek:     `(next )?(monday|tuesday|wednesday|thursday|friday|saturday|sunday)`,
			RuleNextDayOfWeek: "next",
			RuleNaturalDate:   `january|february|march|april|may|june|july|august|september|october|november|december`,
		},
	}
)

// ----------------------------------------------------------

const (
	jokeURL   = "https://official-joke-api.appspot.com/random_joke"
	day       = time.Hour * 24
	adviceURL = "https://api.adviceslip.com/advice"
)
