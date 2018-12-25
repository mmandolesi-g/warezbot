package warez

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"warezbot/emby"
	"warezbot/radarr"
	"warezbot/slack"
)

type SlackEvent struct {
	Token    string `json:"token"`
	TeamID   string `json:"team_id"`
	APIAppID string `json:"api_app_id"`
	Event    struct {
		Type        string `json:"type"`
		Subtype     string `json:"subtype"`
		Text        string `json:"text"`
		Ts          string `json:"ts"`
		Username    string `json:"username"`
		BotID       string `json:"bot_id"`
		Attachments []struct {
			ImageURL    string `json:"image_url"`
			ImageWidth  int    `json:"image_width"`
			ImageHeight int    `json:"image_height"`
			ImageBytes  int    `json:"image_bytes"`
			Fields      []struct {
				Title string `json:"title"`
				Value string `json:"value"`
				Short bool   `json:"short"`
			} `json:"fields"`
			CallbackID string `json:"callback_id"`
			Text       string `json:"text"`
			ID         int    `json:"id"`
			Color      string `json:"color"`
			Actions    []struct {
				ID      string `json:"id"`
				Name    string `json:"name"`
				Text    string `json:"text"`
				Type    string `json:"type"`
				Value   string `json:"value"`
				Style   string `json:"style"`
				Confirm struct {
					Text        string `json:"text"`
					Title       string `json:"title"`
					OkText      string `json:"ok_text"`
					DismissText string `json:"dismiss_text"`
				} `json:"confirm"`
			} `json:"actions"`
		} `json:"attachments"`
		Channel     string `json:"channel"`
		EventTs     string `json:"event_ts"`
		ChannelType string `json:"channel_type"`
	} `json:"event"`
	Type        string   `json:"type"`
	EventID     string   `json:"event_id"`
	EventTime   int      `json:"event_time"`
	AuthedUsers []string `json:"authed_users"`
}

type SlackAction struct {
	Type    string `json:"type"`
	Actions []struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"actions"`
	CallbackID string `json:"callback_id"`
	Team       struct {
		ID     string `json:"id"`
		Domain string `json:"domain"`
	} `json:"team"`
	Channel struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"channel"`
	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
	ActionTs        string `json:"action_ts"`
	MessageTs       string `json:"message_ts"`
	AttachmentID    string `json:"attachment_id"`
	Token           string `json:"token"`
	IsAppUnfurl     bool   `json:"is_app_unfurl"`
	OriginalMessage struct {
		Type        string `json:"type"`
		Subtype     string `json:"subtype"`
		Text        string `json:"text"`
		Ts          string `json:"ts"`
		Username    string `json:"username"`
		BotID       string `json:"bot_id"`
		Attachments []struct {
			CallbackID string `json:"callback_id"`
			Text       string `json:"text"`
			ID         int    `json:"id"`
			Color      string `json:"color"`
			Actions    []struct {
				ID      string `json:"id"`
				Name    string `json:"name"`
				Text    string `json:"text"`
				Type    string `json:"type"`
				Value   string `json:"value"`
				Style   string `json:"style"`
				Confirm struct {
					Text        string `json:"text"`
					Title       string `json:"title"`
					OkText      string `json:"ok_text"`
					DismissText string `json:"dismiss_text"`
				} `json:"confirm"`
			} `json:"actions"`
		} `json:"attachments"`
	} `json:"original_message"`
	ResponseURL string `json:"response_url"`
	TriggerID   string `json:"trigger_id"`
}

type SlackEventFunc func(context.Context, SlackEvent) (Response, error)

type SlackActionFunc func(context.Context, SlackAction) (Response, error)

type Response struct {
	EventType  string
	StatusCode int
}

type Service interface {
	ProcessSlackEvents(context.Context, SlackEvent) (Response, error)
	ProcessSlackActions(context.Context, SlackAction) (Response, error)
}

type service struct {
	emby   *emby.Client
	radarr *radarr.Client
	slack  *slack.Client
}

func NewService(embyClient *emby.Client, radarrClient *radarr.Client, slackClient *slack.Client) (Service, error) {
	return &service{
		emby:   embyClient,
		radarr: radarrClient,
		slack:  slackClient,
	}, nil
}

func (s *service) ProcessSlackEvents(ctx context.Context, request SlackEvent) (Response, error) {
	fmt.Printf("%+v\n", request)
	if request.Event.Type == "app_mention" {
		if strings.Contains(request.Event.Text, "ping") {
			s.slack.Ping()
		}
		if strings.Contains(request.Event.Text, "now playing") {
			sessions, err := s.emby.Sessions(ctx)
			if err != nil {
				fmt.Print(err)
			}
			s.slack.NowPlaying(sessions)
		}
		if strings.Contains(request.Event.Text, "add movie") {
			text := strings.Split(request.Event.Text, " ")
			if len(text) >= 4 {
				movies, err := s.radarr.Search(ctx, text[3:])
				if err != nil {
					fmt.Print(err)
				}

				s.slack.PostSearch(ctx, movies)
			}
		}
	}

	resp := Response{
		EventType:  request.Event.Type,
		StatusCode: http.StatusOK,
	}
	return resp, nil
}

func (s *service) ProcessSlackActions(ctx context.Context, request SlackAction) (Response, error) {
	fmt.Printf("%+v\n", request)

	if request.Type == "interactive_message" {
		if request.CallbackID == "movieDownloadPrompt" {
			s.slack.MsgUpdate(ctx, request.OriginalMessage.Ts, request.User.Name, request.Actions[0].Name)
		}
	}
	// resp := Response{
	// 	StatusCode: http.StatusOK,
	// }
	return Response{}, nil
}
