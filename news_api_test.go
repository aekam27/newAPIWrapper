package news_api_test

import (
	"math/rand"
	news_api "newAPIWrapper"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type NewsAPIDAO struct {
	mock.Mock
}

func TestNewsAPI(t *testing.T) {
	t.Run("Initialize New API Fail", func(t *testing.T) {
		_, err := news_api.InitializeNewsAPI("    ")
		assert.NotNil(t, err)
	})
	t.Run("Initialize with invalid key", func(t *testing.T) {
		invalidKey, err := news_api.InitializeNewsAPI("xxxxxxxxxxxxxxxxxxxxxxxx")
		assert.Nil(t, err)
		_, err = invalidKey.GetEveryThing(map[string]interface{}{"q": "apple"})
		assert.NotNil(t, err)
		_, err = invalidKey.GetTopHeadlines(map[string]interface{}{"q": "apple"})
		assert.NotNil(t, err)
		_, err = invalidKey.GetSources(map[string]interface{}{"q": "apple"})
		assert.NotNil(t, err)
	})

	t.Run("Initialize New API Success With no search string", func(t *testing.T) {
		validKey, err := news_api.InitializeNewsAPI("xxxxxxxxxxxxxxxxxxxxxxxx")
		assert.Nil(t, err)
		queryParams := map[string]interface{}{}
		_, err = validKey.GetEveryThing(queryParams)
		assert.EqualError(t, err, "query string is required")
		_, err = validKey.GetTopHeadlines(queryParams)
		assert.EqualError(t, err, "query string is required")
		_, err = validKey.GetSources(queryParams)
		assert.EqualError(t, err, "query string is required")
	})
	t.Run("Initialize New API Success With empty search string", func(t *testing.T) {
		validKey, err := news_api.InitializeNewsAPI("xxxxxxxxxxxxxxxxxxxxxxxx")
		assert.Nil(t, err)
		queryParams := map[string]interface{}{
			"q": "",
		}
		_, err = validKey.GetEveryThing(queryParams)
		assert.EqualError(t, err, "query string length should be greaterthan equalto 1")
		_, err = validKey.GetTopHeadlines(queryParams)
		assert.EqualError(t, err, "query string length should be greaterthan equalto 1")
		_, err = validKey.GetSources(queryParams)
		assert.EqualError(t, err, "query string length should be greaterthan equalto 1")
	})

	t.Run("Initialize New API Success With search string length greater than 500 characters", func(t *testing.T) {
		validKey, err := news_api.InitializeNewsAPI("xxxxxxxxxxxxxxxxxxxxxxxx")
		assert.Nil(t, err)
		charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		randStr := make([]byte, 501)
		for i := 0; i < 501; i++ {
			randStr[i] = charset[rand.Intn(len(charset))]
		}
		queryParams := map[string]interface{}{
			"q": string(randStr),
		}
		_, err = validKey.GetEveryThing(queryParams)
		assert.EqualError(t, err, "query string length should be lessthan equalto 500")
		_, err = validKey.GetTopHeadlines(queryParams)
		assert.EqualError(t, err, "query string length should be lessthan equalto 500")
		_, err = validKey.GetSources(queryParams)
		assert.EqualError(t, err, "query string length should be lessthan equalto 500")
	})

	t.Run("Initialize New API Success With All Params", func(t *testing.T) {
		validKey, err := news_api.InitializeNewsAPI("xxxxxxxxxxxxxxxxxxxxxxxx")
		assert.Nil(t, err)
		queryParams := map[string]interface{}{
			"q":        "apple",
			"page":     int64(-1),
			"pageSize": int64(120),
		}
		_, err = validKey.GetEveryThing(queryParams)
		assert.Nil(t, err)
		_, err = validKey.GetTopHeadlines(queryParams)
		assert.Nil(t, err)
		_, err = validKey.GetSources(queryParams)
		assert.Nil(t, err)
	})
}
