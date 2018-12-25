package radarr

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	httpTimeout = 10 * time.Second
)

type Movies []struct {
	Title                 string        `json:"title"`
	AlternativeTitles     []interface{} `json:"alternativeTitles"`
	SecondaryYearSourceID int           `json:"secondaryYearSourceId"`
	SortTitle             string        `json:"sortTitle"`
	SizeOnDisk            int           `json:"sizeOnDisk"`
	Status                string        `json:"status"`
	Overview              string        `json:"overview"`
	InCinemas             time.Time     `json:"inCinemas,omitempty"`
	Images                []struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
	} `json:"images"`
	Downloaded          bool          `json:"downloaded"`
	RemotePoster        string        `json:"remotePoster"`
	Year                int           `json:"year"`
	HasFile             bool          `json:"hasFile"`
	ProfileID           int           `json:"profileId"`
	PathState           string        `json:"pathState"`
	Monitored           bool          `json:"monitored"`
	MinimumAvailability string        `json:"minimumAvailability"`
	IsAvailable         bool          `json:"isAvailable"`
	FolderName          string        `json:"folderName"`
	Runtime             int           `json:"runtime"`
	TmdbID              int           `json:"tmdbId"`
	TitleSlug           string        `json:"titleSlug"`
	Genres              []interface{} `json:"genres"`
	Tags                []interface{} `json:"tags"`
	Added               time.Time     `json:"added"`
	Ratings             struct {
		Votes int     `json:"votes"`
		Value float64 `json:"value"`
	} `json:"ratings"`
	QualityProfileID int `json:"qualityProfileId"`
}

type Client struct {
	token   string
	baseURL *url.URL
	http    http.Client
}

func NewClient(host, token string) (*Client, error) {
	base, err := url.Parse(host)
	if err != nil {
		return nil, fmt.Errorf("failed to parse host %q: %v", host, err)
	}
	base.Scheme = "https"
	httpClient := http.Client{
		Timeout: httpTimeout,
	}
	return &Client{
		baseURL: base,
		token:   token,
		http:    httpClient,
	}, nil
}

func (c *Client) Search(ctx context.Context, searchTerm []string) (Movies, error) {
	s := strings.Join(searchTerm, "%20")
	body, err := c.do(ctx, "GET", fmt.Sprintf("movie/lookup?term=%s&apikey=%s", s, c.token))
	if err != nil {
		return nil, err
	}

	var movies Movies
	if err := json.Unmarshal(body, &movies); err != nil {
		return nil, err
	}

	return movies, nil
}

func (c *Client) do(ctx context.Context, method string, path string) ([]byte, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.baseURL, path), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	response, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
