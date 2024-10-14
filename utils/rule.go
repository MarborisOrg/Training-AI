package utils

import (
	"time"
	"regexp"
	"strings"
	"fmt"
	"strconv"
)

func RegisterRule(rule Rule) {
	rules = append(rules, rule)
}

func RuleToday(locale, sentence string) (result time.Time) {
	todayRegex := regexp.MustCompile(RuleTranslations[locale].RuleToday)
	today := todayRegex.FindString(sentence)

	if today == "" {
		return time.Time{}
	}

	return time.Now()
}

func RuleTomorrow(locale, sentence string) (result time.Time) {
	tomorrowRegex := regexp.MustCompile(RuleTranslations[locale].RuleTomorrow)
	date := tomorrowRegex.FindString(sentence)

	if date == "" {
		return time.Time{}
	}

	result = time.Now().Add(day)

	if strings.Contains(date, RuleTranslations[locale].RuleAfterTomorrow) {
		return result.Add(day)
	}

	return
}

func RuleDayOfWeek(locale, sentence string) time.Time {
	dayOfWeekRegex := regexp.MustCompile(RuleTranslations[locale].RuleDayOfWeek)
	date := dayOfWeekRegex.FindString(sentence)

	if date == "" {
		return time.Time{}
	}

	var foundDayOfWeek int

	for _, dayOfWeek := range daysOfWeek {

		stringDayOfWeek := strings.ToLower(dayOfWeek.String())

		if strings.Contains(date, stringDayOfWeek) {
			foundDayOfWeek = int(dayOfWeek)
		}
	}

	currentDay := int(time.Now().Weekday())

	calculatedDate := foundDayOfWeek - currentDay

	if calculatedDate <= 0 {
		calculatedDate += 7
	}

	if strings.Contains(date, RuleTranslations[locale].RuleNextDayOfWeek) {
		calculatedDate += 7
	}

	return time.Now().Add(day * time.Duration(calculatedDate))
}

func RuleNaturalDate(locale, sentence string) time.Time {
	naturalMonthRegex := regexp.MustCompile(
		RuleTranslations[locale].RuleNaturalDate,
	)
	naturalDayRegex := regexp.MustCompile(`\d{2}|\d`)

	month := naturalMonthRegex.FindString(sentence)
	day := naturalDayRegex.FindString(sentence)

	if locale != "en" {
		monthIndex := Index(RuleTranslations[locale].Months, month)
		month = RuleTranslations["en"].Months[monthIndex]
	}

	parsedMonth, _ := time.Parse("January", month)
	parsedDay, _ := strconv.Atoi(day)

	if day == "" && month == "" {
		return time.Time{}
	}

	if day == "" {
		calculatedMonth := parsedMonth.Month() - time.Now().Month()

		if calculatedMonth <= 0 {
			calculatedMonth += 12
		}

		return time.Now().AddDate(0, int(calculatedMonth), -time.Now().Day()+1)
	}

	parsedDate := fmt.Sprintf("%d-%02d-%02d", time.Now().Year(), parsedMonth.Month(), parsedDay)
	date, err := time.Parse("2006-01-02", parsedDate)
	if err != nil {
		return time.Time{}
	}

	if time.Now().After(date) {
		date = date.AddDate(1, 0, 0)
	}

	return date
}

func RuleDate(locale, sentence string) time.Time {
	dateRegex := regexp.MustCompile(`(\d{2}|\d)/(\d{2}|\d)`)
	date := dateRegex.FindString(sentence)

	if date == "" {
		return time.Time{}
	}

	parsedDate, err := time.Parse("01/02", date)
	if err != nil {
		return time.Time{}
	}

	parsedDate = parsedDate.AddDate(time.Now().Year(), 0, 0)

	if time.Now().After(parsedDate) {
		parsedDate = parsedDate.AddDate(1, 0, 0)
	}

	return parsedDate
}