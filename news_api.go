package news_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type NewsAPIDAO interface {
	GetTopHeadlines(queryParams map[string]interface{}) (NewsResp, error)
	GetEveryThing(queryParams map[string]interface{}) (NewsResp, error)
	GetSources(queryParams map[string]interface{}) (SourcesResp, error)
}

type Articles struct {
	Source      interface{} `json:"source,omitempty"`
	Author      string      `json:"author,omitempty"`
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description,omitempty"`
	Url         string      `json:"url,omitempty"`
	UrlToImage  string      `json:"urlToImage,omitempty"`
	PublishedAt string      `json:"publishedAt,omitempty"`
	Content     string      `json:"content,omitempty"`
}

type Sources struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Url         string `json:"url,omitempty"`
	Category    string `json:"category,omitempty"`
	Language    string `json:"language,omitempty"`
	Country     string `json:"country,omitempty"`
}

type NewsResp struct {
	Status       string     `json:"status,omitempty"`
	TotalResults int        `json:"totalResults,omitempty"`
	Articles     []Articles `json:"articles,omitempty"`
	Code         string     `json:"code,omitempty"`
	Message      string     `json:"message,omitempty"`
}

type SourcesResp struct {
	Status  string    `json:"status,omitempty"`
	Sources []Sources `json:"sources,omitempty"`
	Code    string    `json:"code,omitempty"`
	Message string    `json:"message,omitempty"`
}

type newsAPI struct {
	apikey string
}

func InitializeNewsAPI(apikey string) (NewsAPIDAO, error) {
	trimStr := strings.TrimSpace(apikey)
	if len(trimStr) == 0 {
		return nil, errors.New("invalid api key")
	}
	return &newsAPI{
		apikey: apikey,
	}, nil
}

const (
	maxpageSize      = int64(100)
	defaultPageSize  = int64(100)
	defaultPage      = int64(1)
	defaultSortBy    = "publishedAt"
	maxQstringLength = int64(500)
	everyThingUrl    = "https://newsapi.org/v2/everything?"
	topHeadlinesUrl  = "https://newsapi.org/v2/top-headlines?"
	sourceUrl        = "https://newsapi.org/v2/top-headlines/sources?"
)

var (
	allowedSortBys    = []string{"publishedAt", "popularity", "relevancy"}
	allowedLanguage   = []string{"ar", "de", "en", "es", "fr", "he", "it", "nl", "no", "pt", "ru", "sv", "ud", "zh"}
	allowedCountries  = []string{"ae", "ar", "at", "au", "be", "bg", "br", "ca", "ch", "cn", "co", "cu", "cz", "de", "eg", "fr", "gb", "gr", "hk", "hu", "id", "ie", "il", "in", "it", "jp", "kr", "lt", "lv", "ma", "mx", "my", "ng", "nl", "no", "nz", "ph", "pl", "pt", "ro", "rs", "ru", "sa", "se", "sg", "si", "sk", "th", "tr", "tw", "ua", "us", "ve", "za"}
	allowedCategories = []string{"business", "entertainment", "general", "health", "science", "sports", "technology"}
	allowedSearchIn   = []string{"title", "description", "content"}
)

func (rep *newsAPI) GetTopHeadlines(queryParams map[string]interface{}) (NewsResp, error) {
	var (
		apiKey   = rep.apikey
		apiURL   = everyThingUrl
		newsResp = NewsResp{}
	)
	apiURL, err := constructURL(queryParams, apiURL)
	if err != nil {
		return newsResp, err
	}
	resp, err := getRequest(apiURL, apiKey)
	if err != nil {
		return newsResp, err
	}
	err = json.Unmarshal(resp, &newsResp)
	if err != nil {
		return newsResp, err
	}
	if newsResp.Status == "error" {
		return NewsResp{}, errors.New(newsResp.Message)
	}
	return newsResp, nil
}

func (rep *newsAPI) GetEveryThing(queryParams map[string]interface{}) (NewsResp, error) {
	var (
		apiKey   = rep.apikey
		apiURL   = topHeadlinesUrl
		newsResp = NewsResp{}
	)
	apiURL, err := constructURL(queryParams, apiURL)
	if err != nil {
		return newsResp, err
	}
	resp, err := getRequest(apiURL, apiKey)
	if err != nil {
		return newsResp, err
	}
	err = json.Unmarshal(resp, &newsResp)
	if err != nil {
		return newsResp, err
	}
	if newsResp.Status == "error" {
		return NewsResp{}, errors.New(newsResp.Message)
	}
	return newsResp, nil
}

func (rep *newsAPI) GetSources(queryParams map[string]interface{}) (SourcesResp, error) {
	var (
		apiKey     = rep.apikey
		apiURL     = sourceUrl
		sourceResp = SourcesResp{}
	)
	apiURL, err := constructURL(queryParams, apiURL)
	if err != nil {
		return sourceResp, err
	}
	resp, err := getRequest(apiURL, apiKey)
	if err != nil {
		return sourceResp, err
	}
	err = json.Unmarshal(resp, &sourceResp)
	if err != nil {
		return sourceResp, err
	}
	if sourceResp.Status == "error" {
		return SourcesResp{}, errors.New(sourceResp.Message)
	}
	return sourceResp, nil
}

func urlEncodeString(str string) string {
	encodedString := url.QueryEscape(str)
	return encodedString
}

func constructURL(queryParams map[string]interface{}, apiURL string) (string, error) {
	if q, ok := queryParams["q"].(string); ok {
		if len(q) > 500 {
			return "", errors.New("query string length should be lessthan equalto 500")
		} else if len(q) < 1 {
			return "", errors.New("query string length should be greaterthan equalto 1")
		}
		apiURL = fmt.Sprintf("%s%s%s%s", apiURL, "&", "q=", urlEncodeString(q))
	} else {
		return "", errors.New("query string is required")
	}
	if _, ok := queryParams["searchIn"].([]string); ok {
		searchIn := checkIfValueAllowedInStringArray(queryParams["searchIn"].([]string), allowedSearchIn)
		if searchIn != "" {
			apiURL = fmt.Sprintf("%s%s%s%s", apiURL, "&", "searchIn=", urlEncodeString(searchIn))
		}
	}
	if _, ok := queryParams["sources"].([]string); ok {
		sources := checkIfValueAllowedInStringArray(queryParams["sources"].([]string), []string{})
		if sources != "" {
			apiURL = fmt.Sprintf("%s%s%s%s", apiURL, "&", "sources=", urlEncodeString(sources))
		}
	} else {
		if _, ok := queryParams["country"].([]string); ok {
			sources := checkIfValueAllowedInStringArray(queryParams["country"].([]string), allowedCountries)
			if sources != "" {
				apiURL = fmt.Sprintf("%s%s%s%s", apiURL, "&", "country=", urlEncodeString(sources))
			}
		}
		if _, ok := queryParams["category"].([]string); ok {
			sources := checkIfValueAllowedInStringArray(queryParams["category"].([]string), allowedCategories)
			if sources != "" {
				apiURL = fmt.Sprintf("%s%s%s%s", apiURL, "&", "category=", urlEncodeString(sources))
			}
		}
	}
	if _, ok := queryParams["domains"].([]string); ok {
		domains := checkIfValueAllowedInStringArray(queryParams["domains"].([]string), []string{})
		if domains != "" {
			apiURL = fmt.Sprintf("%s%s%s%s", apiURL, "&", "domains=", urlEncodeString(domains))
		}
	}
	if _, ok := queryParams["excludeDomains"].([]string); ok {
		excludeDomains := checkIfValueAllowedInStringArray(queryParams["excludeDomains"].([]string), []string{})
		if excludeDomains != "" {
			apiURL = fmt.Sprintf("%s%s%s%s", apiURL, "&", "excludeDomains=", urlEncodeString(excludeDomains))
		}
	}
	if _, ok := queryParams["from"].(string); ok {
		from, err := parseDTString(queryParams["from"].(string))
		if err == nil {
			apiURL = fmt.Sprintf("%s%s%s%s", apiURL, "&", "from=", urlEncodeString(from))
		} else {
			log.Printf("unable to parse from date retrieving to defaults")
		}
	}
	if _, ok := queryParams["to"].(string); ok {
		to, err := parseDTString(queryParams["to"].(string))
		if err == nil {
			apiURL = fmt.Sprintf("%s%s%s%s", apiURL, "&", "to=", urlEncodeString(to))
		} else {
			log.Printf("unable to parse to date retrieving to defaults")
		}
	}
	if _, ok := queryParams["language"].([]string); ok {
		languages := checkIfValueAllowedInStringArray(queryParams["language"].([]string), allowedLanguage)
		if languages != "" {
			apiURL = fmt.Sprintf("%s%s%s%s", apiURL, "&", "language=", urlEncodeString(languages))
		}
	}
	if _, ok := queryParams["sortBy"].(string); ok {
		sortBy := defaultSortBy
		for _, allowedValue := range allowedSortBys {
			if strings.EqualFold(allowedValue, queryParams["sortBy"].(string)) {
				sortBy = allowedValue
				break
			}
		}
		apiURL = fmt.Sprintf("%s%s%s%s", apiURL, "&", "sortBy=", urlEncodeString(sortBy))
	}
	if pageSize, ok := queryParams["pageSize"].(int64); ok {
		apiURL = fmt.Sprintf("%s%s%s%d", apiURL, "&", "pageSize=", pageSize)
	}
	if page, ok := queryParams["page"].(int64); ok {
		if page < 1 {
			log.Printf("page number is less than 1 using default value")
			page = defaultPage
		}
		if page > 100 {
			log.Printf("page number is greater than maxPage size allowed using maxAllowed Value of 100")
			page = maxpageSize
		}
		apiURL = fmt.Sprintf("%s%s%s%d", apiURL, "&", "page=", page)
	}
	return apiURL, nil
}

// ISO 8601
func parseDTString(dateTimeString string) (string, error) {
	parsedDT, err := time.Parse(time.RFC3339, dateTimeString)
	if err != nil {
		return "", err
	}
	iso8601String := parsedDT.Format(time.RFC3339)
	return iso8601String, nil
}

func checkIfValueAllowedInStringArray(strArr []string, allowedArray []string) string {
	queryValArr := []string{}
	queryString := ""
	for _, value := range strArr {
		if len(allowedArray) > 0 {
			for _, allowedValue := range allowedArray {
				if allowedValue == strings.ToLower(value) {
					queryValArr = append(queryValArr, allowedValue)
					break
				}
			}
		} else {
			queryValArr = append(queryValArr, value)
		}

	}
	if len(queryValArr) > 0 {
		queryString = strings.Join(queryValArr, ",")
	}
	return queryString
}

func getRequest(url, apiKey string) ([]byte, error) {
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Add("X-Api-Key", apiKey)
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}
