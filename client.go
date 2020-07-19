package main

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

func newNewsAPIClient() *NewsAPIClient {
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

type NewsAPIClient struct {
	APIKey  string
	BaseURL string
	Sources []string
}

func (n *NewsAPIClient) fetchArticles(q string, from, to time.Time) ([]*Article, error) {
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
	params.Add("sources", strings.Join(n.Sources, ","))
	params.Add("apiKey", n.APIKey)
	u.RawQuery = params.Encode()

	// Make HTTP request.
	res, err := http.Get(u.String())
	if err != nil {
		fmt.Println("Error making request")
		return nil, err
	}

	fmt.Printf("Got status code: %v\n", res.StatusCode) // 426 Error when date range is too far back.

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

	// Transform api response into a slice of Article objects.
	articles := apiResponse.Transform()

	return articles, nil
}
