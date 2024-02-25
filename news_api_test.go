package news_api_test

import (
	"fmt"
	"math/rand"
	"net/url"
	news_api "newAPIWrapper"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type NewsAPIDAO struct {
	mock.Mock
}

func TestNewsAPI(t *testing.T) {

	t.Run("Construct url with invalid query type", func(t *testing.T) {
		_, err := news_api.ConstructQueryURL("all-news", map[string]interface{}{})
		assert.EqualError(t, err, "invalid query type")
	})
	t.Run("Construct url with valid query type", func(t *testing.T) {
		queryTypes := []string{"everything", "top-headlines", "sources"}
		for _, i := range queryTypes {
			qurl, err := news_api.ConstructQueryURL(i, map[string]interface{}{"q": "apple"})
			assert.Nil(t, err)
			assert.Equal(t, true, strings.Contains(qurl, i))
		}
	})
	t.Run("Construct url with no search string", func(t *testing.T) {
		queryTypes := []string{"everything", "top-headlines", "sources"}
		for _, i := range queryTypes {
			_, err := news_api.ConstructQueryURL(i, map[string]interface{}{})
			assert.EqualError(t, err, "query string is required")
		}
	})

	t.Run("Construct url with no search string", func(t *testing.T) {
		queryTypes := []string{"everything", "top-headlines", "sources"}
		for _, i := range queryTypes {
			_, err := news_api.ConstructQueryURL(i, map[string]interface{}{"q": "    "})
			assert.EqualError(t, err, "query string length should be greaterthan equalto 1")
		}
	})

	t.Run("Construct url with search string length greater than 500 characters", func(t *testing.T) {
		charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		randStr := make([]byte, 501)
		for i := 0; i < 501; i++ {
			randStr[i] = charset[rand.Intn(len(charset))]
		}
		queryParams := map[string]interface{}{
			"q": string(randStr),
		}
		queryTypes := []string{"everything", "top-headlines", "sources"}
		for _, i := range queryTypes {
			_, err := news_api.ConstructQueryURL(i, queryParams)
			assert.EqualError(t, err, "query string length should be lessthan equalto 500")
		}
	})

	t.Run("Construct url with valid search string and searchIn value", func(t *testing.T) {
		queryTypes := []string{"everything", "top-headlines", "sources"}
		for _, i := range queryTypes {
			qurl, err := news_api.ConstructQueryURL(i, map[string]interface{}{"q": "apple", "searchIn": []string{"title", "content"}})
			assert.Nil(t, err)
			assert.Equal(t, true, strings.Contains(qurl, i))
			assert.Equal(t, true, strings.Contains(qurl, "searchIn=title,content"))
		}
	})

	t.Run("Construct url with valid search string and invalid searchIn value", func(t *testing.T) {
		queryTypes := []string{"everything", "top-headlines", "sources"}
		for _, i := range queryTypes {
			qurl, err := news_api.ConstructQueryURL(i, map[string]interface{}{"q": "apple", "searchIn": []string{"title", "contenting"}})
			assert.Nil(t, err)
			assert.Equal(t, true, strings.Contains(qurl, i))
			assert.Equal(t, true, strings.Contains(qurl, "searchIn=title"))
			assert.Equal(t, false, strings.Contains(qurl, "contenting"))
		}
	})

	t.Run("Construct url with valid search string, sources, country, category value", func(t *testing.T) {
		queryTypes := []string{"everything", "top-headlines", "sources"}
		for _, i := range queryTypes {
			qurl, err := news_api.ConstructQueryURL(i, map[string]interface{}{"q": "apple", "sources": []string{"a", "b"},
				"country": []string{"us", "uk"}, "category": []string{"cat1", "cat2"}})
			assert.Nil(t, err)
			assert.Equal(t, true, strings.Contains(qurl, i))
			assert.Equal(t, true, strings.Contains(qurl, fmt.Sprintf("%s%s", "sources=", url.QueryEscape("a,b"))))
			assert.Equal(t, false, strings.Contains(qurl, "country="))
			assert.Equal(t, false, strings.Contains(qurl, "category="))
		}
	})

	t.Run("Construct url with valid search string, empty sources, country, category value", func(t *testing.T) {
		queryTypes := []string{"everything", "top-headlines", "sources"}
		for _, i := range queryTypes {
			qurl, err := news_api.ConstructQueryURL(i, map[string]interface{}{"q": "apple", "sources": []string{},
				"country": []string{"us", "uk"}, "category": []string{"business"}})
			assert.Nil(t, err)
			assert.Equal(t, true, strings.Contains(qurl, i))
			assert.Equal(t, false, strings.Contains(qurl, "sources="))
			assert.Equal(t, true, strings.Contains(qurl, "country=us"))
			assert.Equal(t, true, strings.Contains(qurl, "category=business"))
		}
	})

	t.Run("Construct url with valid search string, domains, excludedDomain and langaugae value", func(t *testing.T) {
		queryTypes := []string{"everything", "top-headlines", "sources"}
		for _, i := range queryTypes {
			qurl, err := news_api.ConstructQueryURL(i, map[string]interface{}{"q": "apple", "domains": []string{"abc.com", "xyz.com"},
				"excludeDomains": []string{"twitter.com", "x.com"}, "language": []string{"en", "ep"}})
			assert.Nil(t, err)
			assert.Equal(t, true, strings.Contains(qurl, i))
			assert.Equal(t, true, strings.Contains(qurl, fmt.Sprintf("%s%s", "domains=", url.QueryEscape("abc.com,xyz.com"))))
			assert.Equal(t, true, strings.Contains(qurl, fmt.Sprintf("%s%s", "excludeDomains=", url.QueryEscape("twitter.com,x.com"))))
			assert.Equal(t, true, strings.Contains(qurl, "language=en"))
		}
	})

	t.Run("Construct url with valid search string, sortBy, pageSize and page value", func(t *testing.T) {
		queryTypes := []string{"everything", "top-headlines", "sources"}
		for _, i := range queryTypes {
			qurl, err := news_api.ConstructQueryURL(i, map[string]interface{}{"q": "apple", "sortBy": []string{"date"},
				"pageSize": int64(-100), "page": int64(-1)})
			assert.Nil(t, err)
			assert.Equal(t, true, strings.Contains(qurl, i))
			assert.Equal(t, false, strings.Contains(qurl, "sortBy="))
			assert.Equal(t, true, strings.Contains(qurl, "pageSize=100"))
			assert.Equal(t, true, strings.Contains(qurl, "page=1"))

			qurl, err = news_api.ConstructQueryURL(i, map[string]interface{}{"q": "apple", "sortBy": "publishedAt",
				"pageSize": int64(101)})
			assert.Nil(t, err)
			assert.Equal(t, true, strings.Contains(qurl, i))
			assert.Equal(t, true, strings.Contains(qurl, "sortBy=publishedAt"))
			assert.Equal(t, true, strings.Contains(qurl, "pageSize=100"))
		}
	})

	t.Run("Construct url with valid search string and dates", func(t *testing.T) {
		queryTypes := []string{"everything", "top-headlines", "sources"}
		for _, i := range queryTypes {
			qurl, err := news_api.ConstructQueryURL(i, map[string]interface{}{"q": "apple", "from": "2024-01-02T00:00:00Z",
				"to": "2024-01-05T15:04:05Z"})
			assert.Nil(t, err)
			assert.Equal(t, true, strings.Contains(qurl, i))
			assert.Equal(t, true, strings.Contains(qurl, "from=2024-01-02T00:00:00Z"))
			assert.Equal(t, true, strings.Contains(qurl, "to=2024-01-05T15:04:05Z"))
		}
	})

	t.Run("Construct url with valid search string and invalid date stamps", func(t *testing.T) {
		queryTypes := []string{"everything", "top-headlines", "sources"}
		for _, i := range queryTypes {
			qurl, err := news_api.ConstructQueryURL(i, map[string]interface{}{"q": "apple", "from": "2025-01-02T00:00:00Z",
				"to": "2024-01-05T15:04:05Z"})
			assert.Nil(t, err)
			assert.Equal(t, true, strings.Contains(qurl, i))
			assert.Equal(t, false, strings.Contains(qurl, "from=2025-01-02T00:00:00Z"))
			assert.Equal(t, false, strings.Contains(qurl, "to=2024-01-05T15:04:05Z"))
		}
	})

	t.Run("Construct url with valid search string and invalid dates", func(t *testing.T) {
		queryTypes := []string{"everything", "top-headlines", "sources"}
		for _, i := range queryTypes {
			qurl, err := news_api.ConstructQueryURL(i, map[string]interface{}{"q": "apple",
				"to": "2024-01-05T15:04:05Z"})
			assert.Nil(t, err)
			assert.Equal(t, true, strings.Contains(qurl, i))
			assert.Equal(t, true, strings.Contains(qurl, "to=2024-01-05T15:04:05Z"))

			qurl, err = news_api.ConstructQueryURL(i, map[string]interface{}{"q": "apple", "from": "2024-01-05T15:04:05Z"})
			assert.Nil(t, err)
			assert.Equal(t, true, strings.Contains(qurl, "from=2024-01-05T15:04:05Z"))
		}
	})

	t.Run("Construct url with valid search string and invalid dates format", func(t *testing.T) {
		queryTypes := []string{"everything", "top-headlines", "sources"}
		for _, i := range queryTypes {
			qurl, err := news_api.ConstructQueryURL(i, map[string]interface{}{"q": "apple",
				"to": "2024-01-05T15:04:05"})
			assert.Nil(t, err)
			assert.Equal(t, true, strings.Contains(qurl, i))
			assert.Equal(t, false, strings.Contains(qurl, "to=2024-01-05T15:04:05"))

			qurl, err = news_api.ConstructQueryURL(i, map[string]interface{}{"q": "apple", "from": "2024-01-05"})
			assert.Nil(t, err)
			assert.Equal(t, false, strings.Contains(qurl, "from=2024-01-05"))
		}
	})

	t.Run("Initialize News API With Invalid Key", func(t *testing.T) {
		_, err := news_api.InitializeNewsAPI("    ")
		assert.NotNil(t, err)
		assert.EqualError(t, err, "invalid api key")
	})

	t.Run("Initialize with invalid key and Valid API URL", func(t *testing.T) {
		invalidKey, err := news_api.InitializeNewsAPI("xxxxxxxxxxxxxxxxxxxxxxxx")
		assert.Nil(t, err)
		queryTypes := []string{"everything", "top-headlines"}
		for _, i := range queryTypes {
			qurl, err := news_api.ConstructQueryURL(i, map[string]interface{}{"q": "apple", "from": "2023-01-02T00:00:00Z",
				"to": "2023-01-15T15:04:05Z", "sortBy": "publishedAt",
				"pageSize": int64(10), "page": 1, "searchIn": []string{"title", "contenting"}})
			assert.Nil(t, err)
			assert.Equal(t, true, strings.Contains(qurl, i))
			_, err = invalidKey.GetNews(qurl)
			assert.EqualError(t, err, "Your API key is invalid or incorrect. Check your key, or go to https://newsapi.org to create a free API key.")
			assert.NotNil(t, err)
		}

		queryTypes = []string{"sources"}
		for _, i := range queryTypes {
			qurl, err := news_api.ConstructQueryURL(i, map[string]interface{}{"q": "apple", "from": "2023-01-02T00:00:00Z",
				"to": "2023-01-15T15:04:05Z", "sortBy": "publishedAt",
				"pageSize": int64(10), "page": 1, "searchIn": []string{"title", "contenting"}})
			assert.Nil(t, err)
			assert.Equal(t, true, strings.Contains(qurl, i))
			_, err = invalidKey.GetSources(qurl)
			assert.EqualError(t, err, "Your API key is invalid or incorrect. Check your key, or go to https://newsapi.org to create a free API key.")
			assert.NotNil(t, err)
		}
	})
}
