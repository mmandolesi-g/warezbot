package radarr

import (
	"bytes"
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
	httpTimeout         = 10 * time.Second
	qualityProfileID    = 3
	rootFolderPath      = "/movies/"
	minimumAvailability = "announced"
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

type AddMovieRequest struct {
	Title               string `json:"title"`
	MinimumAvailability string `json:"minimumAvailability"`
	Images              []struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
	} `json:"images"`
	TMDBID           int    `json:"tmdbId"`
	Year             int    `json:"year"`
	QualityProfileID int    `json:"qualityProfileID"`
	TitleSlug        string `json:"titleSlug"`
	RootFolderPath   string `json:"rootFolderPath"`
	Monitored        bool   `json:"monitored"`
	AddOptions       struct {
		SearchForMovie             bool `json:"searchForMovie"`
		IgnoreEpisodesWithFiles    bool `json:"ignoreEpisodesWithFiles"`
		IgnoreEpisodesWithoutFiles bool `json:"ignoreEpisodesWithoutFiles"`
	} `json:"addOptions"`
}

type AddMovieResponse struct {
	Title      string `json:"title"`
	SortTitle  string `json:"sortTitle"`
	SizeOnDisk int    `json:"sizeOnDisk"`
	Status     string `json:"status"`
	Images     []struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
	} `json:"images"`
	Downloaded          bool          `json:"downloaded"`
	Year                int           `json:"year"`
	HasFile             bool          `json:"hasFile"`
	Path                string        `json:"path"`
	ProfileID           int           `json:"profileId"`
	Monitored           bool          `json:"monitored"`
	MinimumAvailability string        `json:"minimumAvailability"`
	Runtime             int           `json:"runtime"`
	CleanTitle          string        `json:"cleanTitle"`
	ImdbID              string        `json:"imdbId"`
	TmdbID              int           `json:"tmdbId"`
	TitleSlug           string        `json:"titleSlug"`
	Genres              []interface{} `json:"genres"`
	Tags                []interface{} `json:"tags"`
	Added               time.Time     `json:"added"`
	AlternativeTitles   []interface{} `json:"alternativeTitles"`
	QualityProfileID    int           `json:"qualityProfileId"`
	ID                  int           `json:"id"`
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
	body, err := c.do(ctx, "GET", fmt.Sprintf("movie/lookup?term=%s&apikey=%s", s, c.token), nil)
	if err != nil {
		return nil, err
	}

	var movies Movies
	if err := json.Unmarshal(body, &movies); err != nil {
		return nil, err
	}

	return movies, nil
}

func (c *Client) Download(ctx context.Context, id string) (AddMovieResponse, error) {
	x, err := c.do(ctx, "GET", fmt.Sprintf("movie/lookup?term=tmdb:%s&apikey=%s", id, c.token), nil)
	var r []AddMovieRequest
	if err := json.Unmarshal(x, &r); err != nil {
		return AddMovieResponse{}, err
	}

	r[0].QualityProfileID = qualityProfileID
	r[0].Monitored = true
	r[0].RootFolderPath = rootFolderPath
	r[0].AddOptions.SearchForMovie = true
	r[0].MinimumAvailability = minimumAvailability

	input, err := json.Marshal(r[0])
	if err != nil {
		return AddMovieResponse{}, err
	}

	resp, err := c.do(ctx, "POST", fmt.Sprintf("movie?apikey=%s", c.token), input)
	if err != nil {
		return AddMovieResponse{}, err
	}

	var b AddMovieResponse
	if err := json.Unmarshal(resp, &b); err != nil {
		return AddMovieResponse{}, err
	}

	if err := c.runCommand(ctx, "MoviesSearch", b.ID); err != nil {
		return AddMovieResponse{}, err
	}

	return b, nil
}

func (c *Client) runCommand(ctx context.Context, name string, id int) error {

	i := struct {
		Name     string
		MovieIds []int
	}{
		Name:     name,
		MovieIds: []int{id},
	}
	im, err := json.Marshal(i)
	time.Sleep(5 * time.Second)
	resp, err := c.do(ctx, "POST", fmt.Sprintf("command?apikey=%s", c.token), im)
	if err != nil {
		return err
	}
	fmt.Print(string(resp))

	return nil
}

func (c *Client) do(ctx context.Context, method string, path string, input []byte) ([]byte, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.baseURL, path), bytes.NewBuffer(input))
	if err != nil {
		return nil, err
	}

	// req = req.WithContext(ctx)
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
