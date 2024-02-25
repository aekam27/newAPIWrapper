# news_api_wrapper

 Import:

    import "github.com/aekam27/news_api"

 Types:

 The package defines the following types:

        Articles represents a news article.
        type Articles struct {
            Source      interface{} `json:"source,omitempty"`    The identifier id and a display name name for the source this article came from..
            Author      string      `json:"author,omitempty"`    The author of the article.
            Title       string      `json:"title,omitempty"`     The title of the article.
            Description string      `json:"description,omitempty"`  The description or summary of the article.
            Url         string      `json:"url,omitempty"`       The URL to the full article.
            UrlToImage  string      `json:"urlToImage,omitempty"`  The URL to the main image associated with the article.
            PublishedAt string      `json:"publishedAt,omitempty"`  The publication date and time of the article.
            Content     string      `json:"content,omitempty"`    The unformatted content of the article, where available. This is truncated to 200 chars.
        }

        Sources represents a news source.
        type Sources struct {
            Id          string `json:"id,omitempty"`           The identifier of the news source. You can use this with our other endpoints.
            Name        string `json:"name,omitempty"`         The name of the news source.
            Description string `json:"description,omitempty"`  A brief description of the news source.
            Url         string `json:"url,omitempty"`          The URL to the news source's website.
            Category    string `json:"category,omitempty"`     The category of the news source.
            Language    string `json:"language,omitempty"`     The language used by the news source.
            Country     string `json:"country,omitempty"`      The country in which the news source is based in (and primarily writes about).
        }

        NewsResp represents a response containing news articles.
        type NewsResp struct {
            Status       string     `json:"status,omitempty"`        The status of the response. Options: ok, error. In the case of error a code and message property will be populated
            TotalResults int        `json:"totalResults,omitempty"`  The total number of results available for your request. Only a limited number are shown at a time though, so use the page parameter in your requests to page through them.
            Articles     []Articles `json:"articles,omitempty"`      A slice of news articles.
            Code         string     `json:"code,omitempty"`          Error code.
            Message      string     `json:"message,omitempty"`       Error message.
        }

        SourcesResp represents a response containing news sources.
        type SourcesResp struct {
            Status  string    `json:"status,omitempty"`    The status of the response. Options: ok, error. In the case of error a code and message property will be populated
            Sources []Sources `json:"sources,omitempty"`   A slice of news sources.
            Code    string    `json:"code,omitempty"`      Error code.
            Message string    `json:"message,omitempty"`   Error message.
        }


 Functions:

 The package provides the following functions:

    InitializeNewsAPI initializes a new instance of the News API with the provided API key.

    It returns a NewsAPIDAO interface and an error. The NewsAPIDAO
    interface can be used to interact with the News API.

    Parameters:
    - apikey: A string representing the API key to authenticate with the News API.

    Returns:
    - NewsAPIDAO: An interface providing access to News API functionality.
    - error: An error, if any, encountered during initialization. Returns an error
                if the provided API key is empty or contains only whitespace.


    ConstructQueryURL constructs a query URL for a specified query type using the provided parameters.

    It takes a queryType string and a map of queryParams, and returns a constructed URL string
    based on the specified query type and parameters. The constructed URL is specific to the
    query type and may include additional parameters as needed.

    Parameters:
    - queryType: A string representing the type of the query (options: "everything", "top-headlines", "sources").
    - queryParams: A map[string]interface{} containing additional parameters for the query.

    Returns:
    - string: The constructed URL for the specified query type and parameters.
    - error: An error, if any, encountered during the construction process. Returns an error
                if the provided query type is not allowed or if there is an issue constructing the URL.

    Example:
    url, err := ConstructQueryURL("everything", map[string]interface{}{"q": "apple", "category": "business"})
    if err != nil {
        fmt.Println("Error constructing URL:", err)
        return
    }
    fmt.Println("Constructed URL:", url)


Methods:

The package provides the following methods:

    GetNews fetches news articles from the News API based on the provided API URL.

    It takes an API URL as a parameter and returns a NewsResp (News Response) and an error.
    The NewsResp contains information about the fetched news articles, and the error
    indicates if there was any issue during the retrieval or processing of the data.

    Parameters:
    - apiURL: A string representing the API URL for fetching news articles.

    Returns:
    - NewsResp: A struct representing the response containing news articles.
    - error: An error, if any, encountered during the API request or response handling.

    GetSources fetches news sources from the News API based on the provided API URL.

    It takes an API URL as a parameter and returns a SourcesResp (Sources Response) and an error.
    The SourcesResp contains information about the fetched news sources, and the error indicates
    if there was any issue during the retrieval or processing of the data.

    Parameters:
    - apiURL: A string representing the API URL for fetching news sources.

    Returns:
    - SourcesResp: A struct representing the response containing news sources.
    - error: An error, if any, encountered during the API request or response handling.

Constants:

The package defines the following constants:

    - SortBy: "publishedAt", "popularity", "relevancy".
    - Languages: "ar", "de", "en", "es", "fr", "he", "it", "nl", "no", "pt", "ru", "sv", "ud", "zh".
    - Countries: "ae", "ar", "at", "au", "be", "bg", "br", "ca", "ch", "cn", "co", "cu", "cz", "de", "eg", "fr", "gb", "gr", "hk", "hu", "id", "ie", "il", "in", "it", "jp", "kr", "lt", "lv", "ma", "mx", "my", "ng", "nl", "no", "nz", "ph", "pl", "pt", "ro", "rs", "ru", "sa", "se", "sg", "si", "sk", "th", "tr", "tw", "ua", "us", "ve", "za"
    - Categories: "business", "entertainment", "general", "health", "science", "sports", "technology"
    - SearchIn: "title", "description", "content"
    - QueryTypes: "everything", "top-headlines", "sources"


Examples:
            newsAPI, err := news_api.InitializeNewsAPI("xxxxxxxxxxxxxxxxxxxxxxxx")
            qurl, err := news_api.ConstructQueryURL(i, map[string]interface{}{"q": "apple", "from": "2023-01-02T00:00:00Z",
				"to": "2023-01-15T15:04:05Z", "sortBy": "publishedAt",
				"pageSize": int64(10), "page": 1, "searchIn": []string{"title", "contenting"}})
			resp, err = newsAPI.GetSources(qurl)
