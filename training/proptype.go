package training

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tebeka/snowball"
	util "marboris/nout/utils"
)

func (sentence Sentence) WordsBag(words []string) (bag []float64) {
	for _, word := range words {

		var valueToAppend float64
		if util.Contains(sentence.stem(), word) {
			valueToAppend = 1
		}

		bag = append(bag, valueToAppend)
	}

	return bag
}

func (sentence Sentence) stem() (tokenizeWords []string) {
	locale := GetTagByName(sentence.Locale)

	if locale == "" {
		locale = "english"
	}

	tokens := sentence.tokenize()

	stemmer, err := snowball.New(locale)
	if err != nil {
		fmt.Println("Stemmer error", err)
		return
	}

	for _, tokenizeWord := range tokens {
		word := stemmer.Stem(tokenizeWord)
		tokenizeWords = append(tokenizeWords, word)
	}

	return
}

func (sentence Sentence) tokenize() (tokens []string) {
	tokens = strings.Fields(sentence.Content)

	for i, token := range tokens {
		tokens[i] = strings.ToLower(token)
	}

	tokens = removeStopWords(sentence.Locale, tokens)

	return
}

func (sentence *Sentence) arrange() {
	punctuationRegex := regexp.MustCompile(`[a-zA-Z]( )?(\.|\?|!|¿|¡)`)
	sentence.Content = punctuationRegex.ReplaceAllStringFunc(sentence.Content, func(s string) string {
		punctuation := regexp.MustCompile(`(\.|\?|!)`)
		return punctuation.ReplaceAllString(s, "")
	})

	sentence.Content = strings.ReplaceAll(sentence.Content, "-", " ")
	sentence.Content = strings.TrimSpace(sentence.Content)
}
