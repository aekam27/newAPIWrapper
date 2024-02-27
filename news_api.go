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
	GetNews(apiURL string) (NewsResp, error)
	GetSources(apiURL string) (SourcesResp, error)
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
	allowedQueryTypes = []string{"everything", "top-headlines", "sources"}
	queryTypeUrlMap   = map[string]string{
		"everything":    everyThingUrl,
		"top-headlines": topHeadlinesUrl,
		"sources":       sourceUrl,
	}
)

func ConstructQueryURL(queryType string, queryParams map[string]interface{}) (string, error) {
	apiURL := ""
	for _, allowedValue := range allowedQueryTypes {
		if strings.EqualFold(allowedValue, queryType) {
			apiURL = queryTypeUrlMap[allowedValue]
			break
		}
	}
	if apiURL == "" {
		return "", errors.New("invalid query type")
	}

	apiURL, err := constructURL(queryParams, apiURL)
	if err != nil {
		return "", err
	}
	return apiURL, nil
}

func (rep *newsAPI) GetNews(apiURL string) (NewsResp, error) {
	var (
		apiKey   = rep.apikey
		newsResp = NewsResp{}
	)
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

func (rep *newsAPI) GetSources(apiURL string) (SourcesResp, error) {
	var (
		apiKey     = rep.apikey
		sourceResp = SourcesResp{}
	)
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
		q := strings.TrimSpace(q)
		if len(q) > 500 {
			return "", errors.New("query string length should be lessthan equalto 500")
		} else if len(q) < 1 {
			return "", errors.New("query string length should be greaterthan equalto 1")
		}
		apiURL = fmt.Sprintf("%s%s%s", apiURL, "q=", urlEncodeString(q))
	} else {
		return "", errors.New("query string is required")
	}
	if _, ok := queryParams["searchIn"].([]string); ok {
		searchIn := checkIfValueAllowedInStringArray(queryParams["searchIn"].([]string), allowedSearchIn)
		if searchIn != "" {
			apiURL = fmt.Sprintf("%s%s%s%s", apiURL, "&", "searchIn=", searchIn)
		}
	}
	if sourcesArr, ok := queryParams["sources"].([]string); ok {
		if len(sourcesArr) > 0 {
			sources := checkIfValueAllowedInStringArray(sourcesArr, []string{})
			if sources != "" {
				apiURL = fmt.Sprintf("%s%s%s%s", apiURL, "&", "sources=", urlEncodeString(sources))
			}
		} else {
			apiURL = checkForCountryAndCategory(queryParams, apiURL)
		}

	} else {
		apiURL = checkForCountryAndCategory(queryParams, apiURL)
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
		if _, ok := queryParams["to"].(string); ok {
			datesQuery, err := compareForValidtoAndFromDate(queryParams["from"].(string), queryParams["to"].(string))
			if err == nil {
				apiURL = fmt.Sprintf("%s%s%s", apiURL, "&", datesQuery)
			} else {
				log.Print(err.Error())
				log.Printf("unable to parse dates rolling back to defaults")
			}
		} else {
			to, err := parseDTString(queryParams["from"].(string))
			if err == nil {
				apiURL = fmt.Sprintf("%s%s%s%s", apiURL, "&", "from=", to)
			} else {
				log.Print(err.Error())
				log.Printf("unable to parse to date rolling back to defaults")
			}
		}
	} else if _, ok := queryParams["to"].(string); ok {
		to, err := parseDTString(queryParams["to"].(string))
		if err == nil {
			apiURL = fmt.Sprintf("%s%s%s%s", apiURL, "&", "to=", to)
		} else {
			log.Print(err.Error())
			log.Printf("unable to parse to date rolling back to defaults")
		}
	}
	if _, ok := queryParams["language"].([]string); ok {
		languages := checkIfValueAllowedInStringArray(queryParams["language"].([]string), allowedLanguage)
		if languages != "" {
			apiURL = fmt.Sprintf("%s%s%s%s", apiURL, "&", "language=", languages)
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
		apiURL = fmt.Sprintf("%s%s%s%s", apiURL, "&", "sortBy=", sortBy)
	}
	if pageSize, ok := queryParams["pageSize"].(int64); ok {
		if pageSize < 1 {
			pageSize = defaultPageSize
		}
		if pageSize > 100 {
			log.Printf("page number is greater than maxPage size allowed using maxAllowed Value of 100")
			pageSize = maxpageSize
		}
		apiURL = fmt.Sprintf("%s%s%s%d", apiURL, "&", "pageSize=", pageSize)
	}
	if page, ok := queryParams["page"].(int64); ok {
		if page < 1 {
			log.Printf("page number is less than 1 using default value")
			page = defaultPage
		}
		apiURL = fmt.Sprintf("%s%s%s%d", apiURL, "&", "page=", page)
	}
	return apiURL, nil
}

func checkForCountryAndCategory(queryParams map[string]interface{}, apiURL string) string {
	if _, ok := queryParams["country"].([]string); ok {
		country := checkIfValueAllowedInStringArray(queryParams["country"].([]string), allowedCountries)
		if country != "" {
			apiURL = fmt.Sprintf("%s%s%s%s", apiURL, "&", "country=", country)
		}
	}
	if _, ok := queryParams["category"].([]string); ok {
		category := checkIfValueAllowedInStringArray(queryParams["category"].([]string), allowedCategories)
		if category != "" {
			apiURL = fmt.Sprintf("%s%s%s%s", apiURL, "&", "category=", category)
		}
	}
	return apiURL
}

// ISO 8601
func parseDTString(dateTimeString string) (string, error) {
	_, err := time.Parse(time.RFC3339, dateTimeString)
	if err != nil {
		return "", err
	}
	return dateTimeString, nil
}

func compareForValidtoAndFromDate(fromDate, toDate string) (string, error) {
	parsedToDate, err := time.Parse(time.RFC3339, toDate)
	if err != nil {
		return "", err
	}
	parsedFromDate, err := time.Parse(time.RFC3339, fromDate)
	if err != nil {
		return "", err
	}
	if parsedToDate.Before(parsedFromDate) {
		return "", errors.New("invalid dated: to date timestamp is before from date timestamp")
	}
	return fmt.Sprintf("%s%s%s%s", "from=", fromDate, "&to=", toDate), nil
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
