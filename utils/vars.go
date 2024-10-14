package utils

import (
	"time"
)

var (
	rules []Rule

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

	daysOfWeek = map[string]time.Weekday{
		"monday":    time.Monday,
		"tuesday":   time.Tuesday,
		"wednesday": time.Wednesday,
		"thursday":  time.Thursday,
		"friday":    time.Friday,
		"saturday":  time.Saturday,
		"sunday":    time.Sunday,
	}
)

const (
	day = time.Hour * 24
)
