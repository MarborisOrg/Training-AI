package training

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/soudy/mathcat"
	matrix "marboris/nout/matrix"
	util "marboris/nout/utils"
)

func GetUserInformation(token string) Information {
	return userInformation[token]
}

func RegisterModules(locale string, _modules []Modulem) {
	modulesm[locale] = append(modulesm[locale], _modules...)
}

func SerializeCountries() (countries []Country) {
	err := json.Unmarshal(util.ReadFile(util.GetResDir("datasets", "countries.json")), &countries)
	if err != nil {
		fmt.Println(err)
	}

	return countries
}

func FindCountry(locale, sentence string) Country {
	for _, country := range countries {
		name, exists := country.Name[locale]

		if !exists {
			continue
		}

		if !strings.Contains(strings.ToLower(sentence), strings.ToLower(name)) {
			continue
		}

		return country
	}

	return Country{}
}

func AreaReplacer(locale, entry, response, _ string) (string, string) {
	country := FindCountry(locale, entry)

	if country.Currency == "" {
		responseTag := "no country"
		return responseTag, GetMessageu(locale, responseTag)
	}

	return AreaTag, fmt.Sprintf(response, ArticleCountriesm[locale](country.Name[locale]), country.Area)
}

func CapitalReplacer(locale, entry, response, _ string) (string, string) {
	country := FindCountry(locale, entry)

	if country.Currency == "" {
		responseTag := "no country"
		return responseTag, GetMessageu(locale, responseTag)
	}

	articleFunction, exists := ArticleCountriesm[locale]
	countryName := country.Name[locale]
	if exists {
		countryName = articleFunction(countryName)
	}

	return CapitalTag, fmt.Sprintf(response, countryName, country.Capital)
}

func CurrencyReplacer(locale, entry, response, _ string) (string, string) {
	country := FindCountry(locale, entry)

	if country.Currency == "" {
		responseTag := "no country"
		return responseTag, GetMessageu(locale, responseTag)
	}

	return CurrencyTag, fmt.Sprintf(response, ArticleCountriesm[locale](country.Name[locale]), country.Currency)
}

func JokesReplacer(locale, entry, response, _ string) (string, string) {
	resp, err := http.Get(jokeURL)
	if err != nil {
		responseTag := "no jokes"
		return responseTag, GetMessageu(locale, responseTag)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		responseTag := "no jokes"
		return responseTag, GetMessageu(locale, responseTag)
	}

	joke := &Joke{}

	err = json.Unmarshal(body, joke)
	if err != nil {
		responseTag := "no jokes"
		return responseTag, GetMessageu(locale, responseTag)
	}

	jokeStr := joke.Setup + " " + joke.Punchline

	return JokesTag, fmt.Sprintf(response, jokeStr)
}

func FindMathOperation(entry string) string {
	mathRegex := regexp.MustCompile(
		`((\()?(((\d+|pi)(\^\d+|!|.)?)|sqrt|cos|sin|tan|acos|asin|atan|log|ln|abs)( )?[+*\/\-x]?( )?(\))?[+*\/\-]?)+`,
	)

	operation := mathRegex.FindString(entry)

	operation = strings.Replace(operation, "x", "*", -1)
	return strings.TrimSpace(operation)
}

func FindNumberOfDecimals(locale, entry string) int {
	decimalsRegex := regexp.MustCompile(
		MathDecimals[locale],
	)
	numberRegex := regexp.MustCompile(`\d+`)

	decimals := numberRegex.FindString(decimalsRegex.FindString(entry))
	decimalsInt, _ := strconv.Atoi(decimals)

	return decimalsInt
}

func MathReplacer(locale, entry, response, _ string) (string, string) {
	operation := FindMathOperation(entry)

	if operation == "" {
		responseTag := "don't understand"
		return responseTag, GetMessageu(locale, responseTag)
	}

	res, err := mathcat.Eval(operation)
	if err != nil {
		responseTag := "math not valid"
		return responseTag, GetMessageu(locale, responseTag)
	}

	decimals := FindNumberOfDecimals(locale, entry)
	if decimals == 0 {
		decimals = 6
	}

	result := res.FloatString(decimals)

	trailingZerosRegex := regexp.MustCompile(`\.?0+$`)
	result = trailingZerosRegex.ReplaceAllString(result, "")

	return MathTag, fmt.Sprintf(response, result)
}

func NameGetterReplacer(locale, _, response, token string) (string, string) {
	name := GetUserInformation(token).Name

	if strings.TrimSpace(name) == "" {
		responseTag := "don't know name"
		return responseTag, GetMessageu(locale, responseTag)
	}

	return NameGetterTag, fmt.Sprintf(response, name)
}

func SerializeNames() (names []string) {
	namesFile := string(util.ReadFile(util.GetResDir("datasets", "names.txt")))

	names = append(names, strings.Split(strings.TrimSuffix(namesFile, "\n"), "\n")...)
	return
}

func FindName(sentence string) string {
	for _, name := range names {
		if !strings.Contains(strings.ToLower(" "+sentence+" "), " "+name+" ") {
			continue
		}

		return name
	}

	return ""
}

func NameSetterReplacer(locale, entry, response, token string) (string, string) {
	name := FindName(entry)

	if name == "" {
		responseTag := "no name"
		return responseTag, GetMessageu(locale, responseTag)
	}

	name = strings.Title(name)

	ChangeUserInformation(token, func(information Information) Information {
		information.Name = name
		return information
	})

	return NameSetterTag, fmt.Sprintf(response, name)
}

func GetMessageu(locale, tag string) string {
	for _, message := range messages[locale] {

		if message.Tag != tag {
			continue
		}

		if len(message.Messages) == 1 {
			return message.Messages[0]
		}

		rand.NewSource(time.Now().UnixNano()) // Seed
		return message.Messages[rand.Intn(len(message.Messages))]
	}

	return ""
}

func FindRangeLimits(local, entry string) ([]int, error) {
	decimalsRegex := regexp.MustCompile(decimal)
	limitStrArr := decimalsRegex.FindAllString(entry, 2)
	limitArr := make([]int, 0)

	if limitStrArr == nil {
		return make([]int, 0), errors.New("no range")
	}

	if len(limitStrArr) != 2 {
		return nil, errors.New("need 2 numbers, a lower and upper limit")
	}

	for _, v := range limitStrArr {
		num, err := strconv.Atoi(v)
		if err != nil {
			return nil, errors.New("non integer range")
		}
		limitArr = append(limitArr, num)
	}

	sort.Ints(limitArr)
	return limitArr, nil
}

func RandomNumberReplacer(locale, entry, response, _ string) (string, string) {
	limitArr, err := FindRangeLimits(locale, entry)
	if err != nil {
		if limitArr != nil {
			return RandomTag, fmt.Sprintf(response, strconv.Itoa(rand.Intn(100)))
		}

		responseTag := "no random range"
		return responseTag, GetMessageu(locale, responseTag)
	}

	min := limitArr[0]
	max := limitArr[1]
	randNum := rand.Intn((max - min)) + min
	return RandomTag, fmt.Sprintf(response, strconv.Itoa(randNum))
}

func ChangeUserInformation(token string, changer func(Information) Information) {
	userInformation[token] = changer(userInformation[token])
}

func GenresReplacer(locale, entry, response, token string) (string, string) {
	genres := FindMoviesGenres(locale, entry)

	if len(genres) == 0 {
		responseTag := "no genres"
		return responseTag, GetMessageu(locale, responseTag)
	}

	ChangeUserInformation(token, func(information Information) Information {
		for _, genre := range genres {

			if util.Contains(information.MovieGenres, genre) {
				continue
			}

			information.MovieGenres = append(information.MovieGenres, genre)
		}
		return information
	})

	return GenresTag, response
}

func LevenshteinDistance(first, second string) int {
	if first == "" {
		return len(second)
	}
	if second == "" {
		return len(first)
	}

	if first[0] == second[0] {
		return LevenshteinDistance(first[1:], second[1:])
	}

	a := LevenshteinDistance(first[1:], second[1:])
	if b := LevenshteinDistance(first, second[1:]); a > b {
		a = b
	}

	if c := LevenshteinDistance(first[1:], second); a > c {
		a = c
	}

	return a + 1
}

func LevenshteinContains(sentence, matching string, rate int) bool {
	words := strings.Split(sentence, " ")
	for _, word := range words {
		if LevenshteinDistance(word, matching) <= rate {
			return true
		}
	}

	return false
}

func SerializeMovies() (movies []Movie) {
	path := util.GetResDir("datasets", "movies.csv")
	bytes, err := os.Open(path)
	if err != nil {
		bytes, _ = os.Open("../" + path)
	}

	reader := csv.NewReader(bufio.NewReader(bytes))
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		rating, _ := strconv.ParseFloat(line[3], 64)

		movies = append(movies, Movie{
			Name:   line[1],
			Genres: strings.Split(line[2], "|"),
			Rating: rating,
		})
	}

	return
}

func FindMoviesGenres(locale, content string) (output []string) {
	for i, genre := range MoviesGenres[locale] {
		if LevenshteinContains(strings.ToUpper(content), strings.ToUpper(genre), 2) {
			output = append(output, MoviesGenres["en"][i])
		}
	}

	return
}

func MovieSearchReplacer(locale, entry, response, token string) (string, string) {
	genres := FindMoviesGenres(locale, entry)

	if len(genres) == 0 {
		responseTag := "no genres"
		return responseTag, GetMessageu(locale, responseTag)
	}

	movie := SearchMovie(genres[0], token)

	return MoviesTag, fmt.Sprintf(response, movie.Name, movie.Rating)
}

func MovieSearchFromInformationReplacer(locale, _, response, token string) (string, string) {
	genres := GetUserInformation(token).MovieGenres
	if len(genres) == 0 {
		responseTag := "no genres saved"
		return responseTag, GetMessageu(locale, responseTag)
	}

	movie := SearchMovie(genres[rand.Intn(len(genres))], token)
	genresJoined := strings.Join(genres, ", ")
	return MoviesDataTag, fmt.Sprintf(response, genresJoined, movie.Name, movie.Rating)
}

func SearchMovie(genre, userToken string) (output Movie) {
	for _, movie := range movies {
		userMovieBlacklist := GetUserInformation(userToken).MovieBlacklist

		if !util.Contains(movie.Genres, genre) || util.Contains(userMovieBlacklist, movie.Name) {
			continue
		}

		if reflect.DeepEqual(output, Movie{}) || movie.Rating > output.Rating {
			output = movie
		}
	}

	ChangeUserInformation(userToken, func(information Information) Information {
		information.MovieBlacklist = append(information.MovieBlacklist, output.Name)
		return information
	})

	return
}

func ArticleCountries(name string) string {
	if name == "United States" {
		return "the " + name
	}

	return name
}

func AdvicesReplacer(locale, entry, response, _ string) (string, string) {
	resp, err := http.Get(adviceURL)
	if err != nil {
		responseTag := "no advices"
		return responseTag, GetMessageu(locale, responseTag)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		responseTag := "no advices"
		return responseTag, GetMessageu(locale, responseTag)
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	advice := result["slip"].(map[string]interface{})["advice"]

	return AdvicesTag, fmt.Sprintf(response, advice)
}

func CacheIntents(locale string, _intents []Intent) {
	intents[locale] = _intents
}

func SerializeIntents(locale string) (_intents []Intent) {
	err := json.Unmarshal(util.ReadFile(util.GetResDir("locales", "intents.json", locale)), &_intents)
	if err != nil {
		panic(err)
	}

	CacheIntents(locale, _intents)

	return _intents
}

func GetModules(locale string) []Modulem {
	return modulesm[locale]
}

func SerializeModulesIntents(locale string) []Intent {
	registeredModules := GetModules(locale)
	intents := make([]Intent, len(registeredModules))

	for k, module := range registeredModules {
		intents[k] = Intent{
			Tag:       module.Tag,
			Patterns:  module.Patterns,
			Responses: module.Responses,
			Context:   "",
		}
	}

	return intents
}

func removeStopWords(locale string, words []string) []string {
	if len(words) <= 4 {
		return words
	}
	stopWords := string(util.ReadFile(util.GetResDir("locales", "stopwords.txt", locale)))
	var wordsToRemove []string
	for _, stopWord := range strings.Split(stopWords, "\n") {
		for _, word := range words {
			if !strings.Contains(stopWord, word) {
				continue
			}
			wordsToRemove = append(wordsToRemove, word)
		}
	}
	return util.Difference(words, wordsToRemove)
}

func GetTagByName(name string) string {
	for _, locale := range Locales {
		if locale.Name != name {
			continue
		}

		return locale.Tag
	}

	return ""
}

func Organize(locale string) (words, classes []string, documents []Document) {
	intents := append(
		SerializeIntents(locale),
		SerializeModulesIntents(locale)...,
	)

	for _, intent := range intents {
		for _, pattern := range intent.Patterns {

			patternSentence := Sentence{locale, pattern}
			patternSentence.arrange()

			for _, word := range patternSentence.stem() {
				if !util.Contains(words, word) {
					words = append(words, word)
				}
			}

			documents = append(documents, Document{
				patternSentence,
				intent.Tag,
			})
		}

		classes = append(classes, intent.Tag)
	}

	sort.Strings(words)
	sort.Strings(classes)

	return words, classes, documents
}

func TrainData(locale string) (inputs, outputs [][]float64) {
	words, classes, documents := Organize(locale)

	for _, document := range documents {
		outputRow := make([]float64, len(classes))
		bag := document.Sentence.WordsBag(words)

		outputRow[util.Index(classes, document.Tag)] = 1

		inputs = append(inputs, bag)
		outputs = append(outputs, outputRow)
	}

	return inputs, outputs
}

func CreateNeuralNetwork(locale string, rate float64, hiddensNodes int) (neuralNetwork matrix.Network) {
	tempDir := os.TempDir()
	saveFile := filepath.Join(tempDir, "Marboris-Training.json")

	inputs, outputs := TrainData(locale)
	neuralNetwork = matrix.CreateNetwork(locale, rate, inputs, outputs, hiddensNodes)
	neuralNetwork.Train(200)

	neuralNetwork.Save(saveFile)

	return
}
