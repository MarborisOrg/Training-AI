package utils

import (
	"time"
)

type Rule func(string, string) time.Time

type RuleTranslation struct {
	DaysOfWeek        []string
	Months            []string
	RuleToday         string
	RuleTomorrow      string
	RuleAfterTomorrow string
	RuleDayOfWeek     string
	RuleNextDayOfWeek string
	RuleNaturalDate   string
}