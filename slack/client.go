package slack

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/nlopes/slack"

	"github.com/mmandolesi-g/warezbot/emby"
	"github.com/mmandolesi-g/warezbot/radarr"
)

var (
	title      string
	titleValue string
	playStatus string
	footer     string
	imageURL   string
	image404   = "https://www.howtogeek.com/wp-content/uploads/2018/05/2018-06-03-2.png"
)

type Client struct {
	channel string
	botID   string
	client  *slack.Client
}

func NewClient(token string, channel string, botID string) (*Client, error) {
	return &Client{
		channel: channel,
		botID:   botID,
		client:  slack.New(token),
	}, nil
}

func (s *Client) MsgUpdate(ctx context.Context, ts string, name string, title string) error {
	attachment := slack.Attachment{
		Text:  fmt.Sprintf("Download process started by %s for %s", name, title),
		Color: makeHexColor(),
	}
	s.UpdateMessage(ts, slack.MsgOptionAttachments(attachment))

	return nil
}

func (s *Client) PostEmbySearch(ctx context.Context, results emby.SearchResults) {
	var count int
	for _, result := range results.SearchHints {
		if result.Type == "Movie" || result.Type == "Episode" || result.Type == "Series" {
			count++

			var image string
			if result.ItemImages.TotalRecordCount > 0 {
				image = result.ItemImages.Images[0].URL
			} else {
				image = image404
			}

			attachment := slack.Attachment{
				Color:      makeHexColor(),
				CallbackID: "embySearchResult",
				Footer:     result.ItemDetail.Overview,
				ImageURL:   image,
				Fields: []slack.AttachmentField{
					{
						Title: fmt.Sprintf("%s - %s", result.Name, result.Type),
						Value: strconv.Itoa(result.ProductionYear),
					},
				},
			}
			s.PostMessage(slack.MsgOptionText("", false), slack.MsgOptionAttachments(attachment))
		}
	}

	attachment := slack.Attachment{
		Color:    makeHexColor(),
		ImageURL: "",
		Fields: []slack.AttachmentField{
			{
				Title: "Total Results found:",
				Value: strconv.Itoa(count),
			},
		},
	}
	s.PostMessage(slack.MsgOptionText("", false), slack.MsgOptionAttachments(attachment))
}

func (s *Client) PostSearch(ctx context.Context, movies radarr.Movies) {
	var attachmentActions []slack.AttachmentAction
	// Only post the first 5 movies in the search
	if len(movies) > 5 {
		movies = movies[:5]
	}
	for i, movie := range movies {
		i++
		attachmentActions = append(attachmentActions, slack.AttachmentAction{
			Name:  strconv.Itoa(movie.TmdbID),
			Type:  "button",
			Text:  fmt.Sprintf("%d.) %s - %d", i, movie.Title, movie.Year),
			Value: fmt.Sprintf("%d.) %s - %d", i, movie.Title, movie.Year),
			Confirm: &slack.ConfirmationField{
				Text: fmt.Sprintf("Are you sure you want to download %s (%s)?", movie.Title, strconv.Itoa(movie.Year)),
			},
		})
		attachment := slack.Attachment{
			Color:      makeHexColor(),
			CallbackID: "movieSearchResult",
			Text:       fmt.Sprintf("TMDB ID: %s", strconv.Itoa(movie.TmdbID)),
			ImageURL:   movie.Images[0].URL,
			Footer:     movie.Overview,
			Fields: []slack.AttachmentField{
				{
					Title: fmt.Sprintf("%d.) %s", i, movie.Title),
					Value: strconv.Itoa(movie.Year),
				},
			},
		}
		s.PostMessage(slack.MsgOptionText("", false), slack.MsgOptionAttachments(attachment))
	}

	attachment := slack.Attachment{
		Color:      makeHexColor(),
		Text:       "Select movie to download",
		CallbackID: "movieDownloadPrompt",
		Actions:    attachmentActions,
	}
	s.PostMessage(slack.MsgOptionText("", false), slack.MsgOptionAttachments(attachment))
}

func (s *Client) NowPlaying(sessions emby.Sessions) {
	for _, ses := range sessions {
		if ses.NowPlayingItem.Name != "" {
			if ses.NowPlayingItem.Type == "Episode" {
				title = fmt.Sprintf("%s is playing the TV show:", ses.UserName)
				titleValue = fmt.Sprintf("%s - %s (Season %d - %d)", ses.NowPlayingItem.SeriesName,
					ses.NowPlayingItem.Name,
					ses.NowPlayingItem.ParentIndexNumber,
					ses.NowPlayingItem.IndexNumber)
			} else {
				title = fmt.Sprintf("%s is playing the film:", ses.UserName)
				titleValue = ses.NowPlayingItem.Name
			}

			if ses.PlayState.IsPaused {
				playStatus = "Paused"
			} else {
				playStatus = "Playing"
			}

			if ses.ItemDetail.Overview != "" {
				footer = ses.ItemDetail.Overview
			} else {
				footer = "No overview found..."
			}

			if ses.ItemImages.TotalRecordCount > 0 {
				imageURL = ses.ItemImages.Images[0].URL
			} else {
				imageURL = image404
			}

			attachment := slack.Attachment{
				Text:       fmt.Sprintf("%s - %g%%", playStatus, percentComplete(ses.PlayState.PositionTicks, ses.NowPlayingItem.RunTimeTicks)),
				ImageURL:   imageURL,
				Color:      makeHexColor(),
				AuthorName: fmt.Sprintf("%s - %s", ses.DeviceName, ses.Client),
				AuthorIcon: ses.AppIconURL,
				Footer:     footer,
				Fields: []slack.AttachmentField{
					{
						Title: title,
						Value: titleValue,
					},
				},
			}
			s.PostMessage(slack.MsgOptionText("", false), slack.MsgOptionAttachments(attachment))
		}
	}
}

func (s *Client) Ping() {
	attachment := slack.Attachment{
		ImageURL:   "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQ3cSyJ752nVIUGaR4QRzf1qNfKA2XsFAtqWZ78c3xIJlFOvpqh",
		Color:      makeHexColor(),
		AuthorName: "warezbot",
		AuthorIcon: "https://i.imgur.com/s0F5TJA.jpg",
	}
	s.PostMessage(slack.MsgOptionText("pong", false), slack.MsgOptionAttachments(attachment))
}

func (s *Client) PostMessage(options ...slack.MsgOption) (string, string, error) {
	return s.client.PostMessage(s.channel, options...)
}

func (s *Client) UpdateMessage(timestamp string, options ...slack.MsgOption) (string, string, string, error) {
	return s.client.UpdateMessage(s.channel, timestamp, options...)
}

func makeHexColor() string {
	src := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 6/2)
	src.Read(b)
	return hex.EncodeToString(b)[:6]
}

func percentComplete(positionTicks int64, runTimeTicks int64) float64 {
	return math.Round((float64(positionTicks) * 100) / float64(runTimeTicks))
}
