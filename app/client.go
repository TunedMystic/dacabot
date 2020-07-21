package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// NewNewsAPIClient creates a new api client for News API.
func NewNewsAPIClient() *NewsAPIClient {
	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		log.Fatal("NEWS_API_KEY is not set")
	}
	return &NewsAPIClient{
		APIKey:  apiKey,
		BaseURL: "https://newsapi.org/v2",
		Sources: []string{
			"abc-news", "bloomberg", "cbs-news",
			"cnn", "fox-news", "google-news",
			"msnbc", "nbc-news", "newsweek",
			"the-hill", "the-huffington-post",
			"the-next-web", "the-wall-street-journal",
			"the-washington-post", "the-washington-times",
			"usa-today",
		},
	}
}

// NewsAPIClient is an api client for NewsAPI.
type NewsAPIClient struct {
	APIKey  string
	BaseURL string
	Sources []string
}

// GetArticles fetches new articles with the given params.
func (n *NewsAPIClient) GetArticles(q string, from, to time.Time) ([]*Article, error) {
	u, err := url.Parse(n.BaseURL + "/everything")
	if err != nil {
		log.Fatalf("Error parsing url: %v\n", err)
	}

	// Prepare url params.
	params := url.Values{}
	params.Add("qInTitle", q)
	params.Add("from", from.Format("2006-01-02"))
	params.Add("to", to.Format("2006-01-02"))
	params.Add("language", "en")
	params.Add("sortBy", "relevancy")
	params.Add("pageSize", "100")
	params.Add("sources", strings.Join(n.Sources, ","))
	params.Add("apiKey", n.APIKey)
	u.RawQuery = params.Encode()

	// Make HTTP request.
	res, err := http.Get(u.String())
	if err != nil {
		fmt.Println("Error making request")
		return nil, err
	}

	// 426 Error when date range is too far back.
	fmt.Printf("Got status code: %v\n", res.StatusCode)

	// Read the response body.
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading the response body")
		return nil, err
	}
	res.Body.Close()

	// Unmarshal the JSON data.
	apiResponse := NewsAPIResponse{}
	err = json.Unmarshal(data, &apiResponse)
	if err != nil {
		fmt.Println("Error unmarshalling the json data")
		return nil, err
	}

	// Transform api response into a slice of *Article objects.
	articles := apiResponse.transform()

	return articles, nil
}

type sourceJSON struct {
	Name string `json:"id"`
}

type articleJSON struct {
	Article
	Source sourceJSON `json:"source"`
}

// NewsAPIResponse is the response format for News API.
/* {
    "status": "ok",
    "totalResults": 22,
    "articles": [{
        "source": {
            "id": "source-id",
            "name": "Source Name"
        },
        "author": "Some author",
        "title": "Some title",
        "description": "Some description",
        "url": "https://some-url.com/article",
        "urlToImage": "https://some-url-image.png",
        "publishedAt": "2020-06-25T23:44:16Z",
        "content": "Truncated content text [+2023 chars]"
	}]
}
*/
type NewsAPIResponse struct {
	Status       string        `json:"status"`
	TotalResults int           `json:"totalResults"`
	Articles     []articleJSON `json:"articles"`
}

// Transform converts article data to objects of type *Article.
func (n NewsAPIResponse) transform() []*Article {
	articles := []*Article{}

	for _, articleJSON := range n.Articles {
		article := &Article{
			URL:         articleJSON.URL,
			Title:       articleJSON.Title,
			Description: articleJSON.Description,
			Source:      articleJSON.Source.Name,
			Author:      articleJSON.Author,
			LedeImg:     articleJSON.LedeImg,
			PublishedAt: articleJSON.PublishedAt,
			CreatedAt:   time.Now().UTC(),
		}
		articles = append(articles, article)
	}

	return articles
}
