package warez

import (
	"context"
	"net/http"
	"strings"
	"time"

	"warezbot/emby"
	"warezbot/radarr"
	"warezbot/slack"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	slackClient "github.com/nlopes/slack"
)

const (
	addMovie   = "add movie"
	nowPlaying = "now playing"
	ping       = "ping"
	search     = "search"
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

type EmbyEvent struct {
	Event string `json:"Event"`
	User  struct {
		Name                      string    `json:"Name"`
		ServerID                  string    `json:"ServerId"`
		ConnectUserName           string    `json:"ConnectUserName"`
		ConnectLinkType           string    `json:"ConnectLinkType"`
		ID                        string    `json:"Id"`
		PrimaryImageTag           string    `json:"PrimaryImageTag"`
		HasPassword               bool      `json:"HasPassword"`
		HasConfiguredPassword     bool      `json:"HasConfiguredPassword"`
		HasConfiguredEasyPassword bool      `json:"HasConfiguredEasyPassword"`
		LastLoginDate             time.Time `json:"LastLoginDate"`
		LastActivityDate          time.Time `json:"LastActivityDate"`
		Configuration             struct {
			AudioLanguagePreference    string        `json:"AudioLanguagePreference"`
			PlayDefaultAudioTrack      bool          `json:"PlayDefaultAudioTrack"`
			SubtitleLanguagePreference string        `json:"SubtitleLanguagePreference"`
			DisplayMissingEpisodes     bool          `json:"DisplayMissingEpisodes"`
			GroupedFolders             []interface{} `json:"GroupedFolders"`
			SubtitleMode               string        `json:"SubtitleMode"`
			DisplayCollectionsView     bool          `json:"DisplayCollectionsView"`
			EnableLocalPassword        bool          `json:"EnableLocalPassword"`
			OrderedViews               []string      `json:"OrderedViews"`
			LatestItemsExcludes        []interface{} `json:"LatestItemsExcludes"`
			MyMediaExcludes            []interface{} `json:"MyMediaExcludes"`
			HidePlayedInLatest         bool          `json:"HidePlayedInLatest"`
			RememberAudioSelections    bool          `json:"RememberAudioSelections"`
			RememberSubtitleSelections bool          `json:"RememberSubtitleSelections"`
			EnableNextEpisodeAutoPlay  bool          `json:"EnableNextEpisodeAutoPlay"`
		} `json:"Configuration"`
		Policy struct {
			IsAdministrator                  bool          `json:"IsAdministrator"`
			IsHidden                         bool          `json:"IsHidden"`
			IsHiddenRemotely                 bool          `json:"IsHiddenRemotely"`
			IsDisabled                       bool          `json:"IsDisabled"`
			BlockedTags                      []interface{} `json:"BlockedTags"`
			IsTagBlockingModeInclusive       bool          `json:"IsTagBlockingModeInclusive"`
			EnableUserPreferenceAccess       bool          `json:"EnableUserPreferenceAccess"`
			AccessSchedules                  []interface{} `json:"AccessSchedules"`
			BlockUnratedItems                []interface{} `json:"BlockUnratedItems"`
			EnableRemoteControlOfOtherUsers  bool          `json:"EnableRemoteControlOfOtherUsers"`
			EnableSharedDeviceControl        bool          `json:"EnableSharedDeviceControl"`
			EnableRemoteAccess               bool          `json:"EnableRemoteAccess"`
			EnableLiveTvManagement           bool          `json:"EnableLiveTvManagement"`
			EnableLiveTvAccess               bool          `json:"EnableLiveTvAccess"`
			EnableMediaPlayback              bool          `json:"EnableMediaPlayback"`
			EnableAudioPlaybackTranscoding   bool          `json:"EnableAudioPlaybackTranscoding"`
			EnableVideoPlaybackTranscoding   bool          `json:"EnableVideoPlaybackTranscoding"`
			EnablePlaybackRemuxing           bool          `json:"EnablePlaybackRemuxing"`
			EnableContentDeletion            bool          `json:"EnableContentDeletion"`
			EnableContentDeletionFromFolders []interface{} `json:"EnableContentDeletionFromFolders"`
			EnableContentDownloading         bool          `json:"EnableContentDownloading"`
			EnableSubtitleDownloading        bool          `json:"EnableSubtitleDownloading"`
			EnableSubtitleManagement         bool          `json:"EnableSubtitleManagement"`
			EnableSyncTranscoding            bool          `json:"EnableSyncTranscoding"`
			EnableMediaConversion            bool          `json:"EnableMediaConversion"`
			EnabledDevices                   []interface{} `json:"EnabledDevices"`
			EnableAllDevices                 bool          `json:"EnableAllDevices"`
			EnabledChannels                  []interface{} `json:"EnabledChannels"`
			EnableAllChannels                bool          `json:"EnableAllChannels"`
			EnabledFolders                   []interface{} `json:"EnabledFolders"`
			EnableAllFolders                 bool          `json:"EnableAllFolders"`
			InvalidLoginAttemptCount         int           `json:"InvalidLoginAttemptCount"`
			EnablePublicSharing              bool          `json:"EnablePublicSharing"`
			RemoteClientBitrateLimit         int           `json:"RemoteClientBitrateLimit"`
			AuthenticationProviderID         string        `json:"AuthenticationProviderId"`
			ExcludedSubFolders               []interface{} `json:"ExcludedSubFolders"`
			DisablePremiumFeatures           bool          `json:"DisablePremiumFeatures"`
			SimultaneousStreamLimit          int           `json:"SimultaneousStreamLimit"`
		} `json:"Policy"`
		PrimaryImageAspectRatio float64 `json:"PrimaryImageAspectRatio"`
	} `json:"User"`
	Item struct {
		Name                  string    `json:"Name"`
		ServerID              string    `json:"ServerId"`
		ID                    string    `json:"Id"`
		DateCreated           time.Time `json:"DateCreated"`
		PresentationUniqueKey string    `json:"PresentationUniqueKey"`
		Container             string    `json:"Container"`
		PremiereDate          time.Time `json:"PremiereDate"`
		ExternalUrls          []struct {
			Name string `json:"Name"`
			URL  string `json:"Url"`
		} `json:"ExternalUrls"`
		Path              string        `json:"Path"`
		Overview          string        `json:"Overview"`
		Taglines          []interface{} `json:"Taglines"`
		Genres            []interface{} `json:"Genres"`
		RunTimeTicks      int64         `json:"RunTimeTicks"`
		ProductionYear    int           `json:"ProductionYear"`
		IndexNumber       int           `json:"IndexNumber"`
		ParentIndexNumber int           `json:"ParentIndexNumber"`
		ProviderIds       struct {
			Tvdb string `json:"Tvdb"`
			Imdb string `json:"Imdb"`
		} `json:"ProviderIds"`
		IsFolder                bool          `json:"IsFolder"`
		ParentID                string        `json:"ParentId"`
		Type                    string        `json:"Type"`
		Studios                 []interface{} `json:"Studios"`
		GenreItems              []interface{} `json:"GenreItems"`
		ParentBackdropItemID    string        `json:"ParentBackdropItemId"`
		ParentBackdropImageTags []string      `json:"ParentBackdropImageTags"`
		SeriesName              string        `json:"SeriesName"`
		SeriesID                string        `json:"SeriesId"`
		SeasonID                string        `json:"SeasonId"`
		PrimaryImageAspectRatio float64       `json:"PrimaryImageAspectRatio"`
		SeriesPrimaryImageTag   string        `json:"SeriesPrimaryImageTag"`
		SeasonName              string        `json:"SeasonName"`
		MediaStreams            []struct {
			Codec                  string  `json:"Codec"`
			ColorTransfer          string  `json:"ColorTransfer,omitempty"`
			ColorPrimaries         string  `json:"ColorPrimaries,omitempty"`
			ColorSpace             string  `json:"ColorSpace,omitempty"`
			TimeBase               string  `json:"TimeBase"`
			CodecTimeBase          string  `json:"CodecTimeBase"`
			VideoRange             string  `json:"VideoRange,omitempty"`
			DisplayTitle           string  `json:"DisplayTitle"`
			NalLengthSize          string  `json:"NalLengthSize,omitempty"`
			IsInterlaced           bool    `json:"IsInterlaced"`
			IsAVC                  bool    `json:"IsAVC,omitempty"`
			BitRate                int     `json:"BitRate"`
			BitDepth               int     `json:"BitDepth,omitempty"`
			RefFrames              int     `json:"RefFrames,omitempty"`
			IsDefault              bool    `json:"IsDefault"`
			IsForced               bool    `json:"IsForced"`
			Height                 int     `json:"Height,omitempty"`
			Width                  int     `json:"Width,omitempty"`
			AverageFrameRate       float64 `json:"AverageFrameRate,omitempty"`
			RealFrameRate          float64 `json:"RealFrameRate,omitempty"`
			Profile                string  `json:"Profile"`
			Type                   string  `json:"Type"`
			AspectRatio            string  `json:"AspectRatio,omitempty"`
			Index                  int     `json:"Index"`
			IsExternal             bool    `json:"IsExternal"`
			IsTextSubtitleStream   bool    `json:"IsTextSubtitleStream"`
			SupportsExternalStream bool    `json:"SupportsExternalStream"`
			Protocol               string  `json:"Protocol"`
			PixelFormat            string  `json:"PixelFormat,omitempty"`
			Level                  int     `json:"Level"`
			IsAnamorphic           bool    `json:"IsAnamorphic,omitempty"`
			ChannelLayout          string  `json:"ChannelLayout,omitempty"`
			Channels               int     `json:"Channels,omitempty"`
			SampleRate             int     `json:"SampleRate,omitempty"`
		} `json:"MediaStreams"`
		ImageTags struct {
			Primary string `json:"Primary"`
		} `json:"ImageTags"`
		BackdropImageTags []interface{} `json:"BackdropImageTags"`
		Chapters          []struct {
			StartPositionTicks int    `json:"StartPositionTicks"`
			Name               string `json:"Name"`
		} `json:"Chapters"`
		MediaType string `json:"MediaType"`
		Width     int    `json:"Width"`
		Height    int    `json:"Height"`
	} `json:"Item"`
	Server struct {
		Name string `json:"Name"`
		ID   string `json:"Id"`
	} `json:"Server"`
	Session struct {
		RemoteEndPoint     string `json:"RemoteEndPoint"`
		Client             string `json:"Client"`
		DeviceName         string `json:"DeviceName"`
		DeviceID           string `json:"DeviceId"`
		ApplicationVersion string `json:"ApplicationVersion"`
		ID                 string `json:"Id"`
	} `json:"Session"`
}

type SlackEventFunc func(context.Context, SlackEvent) (Response, error)

type EmbyEventFunc func(context.Context, EmbyEvent) (Response, error)

type SlackActionFunc func(context.Context, SlackAction) (Response, error)

type Response struct {
	EventType  string
	StatusCode int
}

type Service interface {
	ProcessSlackEvents(context.Context, SlackEvent) (Response, error)
	ProcessSlackActions(context.Context, SlackAction) (Response, error)
	ProcessEmbyEvents(context.Context, EmbyEvent) (Response, error)
}

type service struct {
	emby   *emby.Client
	radarr *radarr.Client
	slack  *slack.Client
	logger log.Logger
}

func NewService(embyClient *emby.Client, radarrClient *radarr.Client, slackClient *slack.Client, log log.Logger) (Service, error) {
	return &service{
		emby:   embyClient,
		radarr: radarrClient,
		slack:  slackClient,
		logger: log,
	}, nil
}

func (s *service) ProcessSlackEvents(ctx context.Context, request SlackEvent) (Response, error) {
	if request.Event.Type == "message" {
		if strings.Contains(request.Event.Text, ping) {
			s.slack.Ping()
		}
		if strings.Contains(request.Event.Text, nowPlaying) {
			sessions, err := s.emby.Sessions(ctx)
			if err != nil {
				level.Error(s.logger).Log("error", err)
				return Response{}, err
			}
			s.slack.NowPlaying(sessions)
		}
		if strings.Contains(request.Event.Text, addMovie) {
			go func() {
				text := strings.Split(request.Event.Text, " ")
				if len(text) >= 4 {
					movies, err := s.radarr.Search(ctx, text[3:])
					if err != nil {
						level.Error(s.logger).Log("error", err)
					}
					s.slack.PostSearch(ctx, movies)
				}
			}()
		}
		if strings.Contains(request.Event.Text, search) {
			go func() {
				text := strings.Split(request.Event.Text, " ")
				if len(text) >= 3 {
					sResults, err := s.emby.Search(ctx, text[2:])
					if err != nil {
						level.Error(s.logger).Log("error", err)
					}
					s.slack.PostEmbySearch(ctx, sResults)
				}
			}()
		}
	}

	return Response{
		EventType:  request.Event.Type,
		StatusCode: http.StatusOK,
	}, nil
}

func (s *service) ProcessSlackActions(ctx context.Context, request SlackAction) (Response, error) {
	if request.Type == "interactive_message" {
		if request.CallbackID == "movieDownloadPrompt" {
			go func() {
				s.slack.MsgUpdate(ctx, request.OriginalMessage.Ts, request.User.Name, request.Actions[0].Value)
				s.radarr.Download(ctx, request.Actions[0].Name)
			}()
		}
	}

	return Response{
		StatusCode: http.StatusAccepted,
	}, nil
}

func (s *service) ProcessEmbyEvents(ctx context.Context, request EmbyEvent) (Response, error) {
	var url string
	if len(request.Item.ExternalUrls) > 0 {
		url = request.Item.ExternalUrls[0].URL
	}

	attachment := slackClient.Attachment{
		CallbackID: "embyEvent",
		ImageURL:   url,
		Fields: []slackClient.AttachmentField{
			{
				Title: request.Event,
				Value: request.User.Name,
			},
		},
	}
	_, _, err := s.slack.PostMessage(slackClient.MsgOptionText("", false), slackClient.MsgOptionAttachments(attachment))
	if err != nil {
		level.Error(s.logger).Log("error", err)
		return Response{}, err
	}

	return Response{
		EventType:  request.Event,
		StatusCode: http.StatusOK,
	}, nil
}
