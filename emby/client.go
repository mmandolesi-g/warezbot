package emby

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	httpTimeout = 10 * time.Second
)

type Sessions []struct {
	ItemDetail ItemDetail
	ItemImages ItemImages
	PlayState  struct {
		PositionTicks int64  `json:"PositionTicks"`
		CanSeek       bool   `json:"CanSeek"`
		IsPaused      bool   `json:"IsPaused"`
		IsMuted       bool   `json:"IsMuted"`
		RepeatMode    string `json:"RepeatMode"`
	} `json:"PlayState"`
	AdditionalUsers []interface{} `json:"AdditionalUsers"`
	Capabilities    struct {
		PlayableMediaTypes           []interface{} `json:"PlayableMediaTypes"`
		SupportedCommands            []interface{} `json:"SupportedCommands"`
		SupportsMediaControl         bool          `json:"SupportsMediaControl"`
		SupportsPersistentIdentifier bool          `json:"SupportsPersistentIdentifier"`
		SupportsSync                 bool          `json:"SupportsSync"`
	} `json:"Capabilities"`
	RemoteEndPoint        string        `json:"RemoteEndPoint"`
	PlayableMediaTypes    []interface{} `json:"PlayableMediaTypes"`
	ID                    string        `json:"Id"`
	ServerID              string        `json:"ServerId"`
	Client                string        `json:"Client"`
	LastActivityDate      time.Time     `json:"LastActivityDate"`
	DeviceName            string        `json:"DeviceName"`
	DeviceID              string        `json:"DeviceId"`
	ApplicationVersion    string        `json:"ApplicationVersion"`
	SupportedCommands     []interface{} `json:"SupportedCommands"`
	SupportsRemoteControl bool          `json:"SupportsRemoteControl"`
	PlaylistItemID        string        `json:"PlaylistItemId,omitempty"`
	UserID                string        `json:"UserId,omitempty"`
	UserName              string        `json:"UserName,omitempty"`
	UserPrimaryImageTag   string        `json:"UserPrimaryImageTag,omitempty"`
	NowPlayingItem        struct {
		Name         string    `json:"Name"`
		ServerID     string    `json:"ServerId"`
		ID           string    `json:"Id"`
		DateCreated  time.Time `json:"DateCreated"`
		HasSubtitles bool      `json:"HasSubtitles"`
		Container    string    `json:"Container"`
		PremiereDate time.Time `json:"PremiereDate"`
		ExternalUrls []struct {
			Name string `json:"Name"`
			URL  string `json:"Url"`
		} `json:"ExternalUrls"`
		Path                     string        `json:"Path"`
		EnableMediaSourceDisplay bool          `json:"EnableMediaSourceDisplay"`
		Overview                 string        `json:"Overview"`
		Taglines                 []interface{} `json:"Taglines"`
		Genres                   []interface{} `json:"Genres"`
		RunTimeTicks             int64         `json:"RunTimeTicks"`
		ProductionYear           int           `json:"ProductionYear"`
		IndexNumber              int           `json:"IndexNumber"`
		ParentIndexNumber        int           `json:"ParentIndexNumber"`
		ProviderIds              struct {
			Tvdb string `json:"Tvdb"`
		} `json:"ProviderIds"`
		IsFolder                bool          `json:"IsFolder"`
		ParentID                string        `json:"ParentId"`
		Type                    string        `json:"Type"`
		Studios                 []interface{} `json:"Studios"`
		GenreItems              []interface{} `json:"GenreItems"`
		ParentLogoItemID        string        `json:"ParentLogoItemId"`
		ParentBackdropItemID    string        `json:"ParentBackdropItemId"`
		ParentBackdropImageTags []string      `json:"ParentBackdropImageTags"`
		LocalTrailerCount       int           `json:"LocalTrailerCount"`
		SeriesName              string        `json:"SeriesName"`
		SeriesID                string        `json:"SeriesId"`
		SeasonID                string        `json:"SeasonId"`
		SpecialFeatureCount     int           `json:"SpecialFeatureCount"`
		PrimaryImageAspectRatio float64       `json:"PrimaryImageAspectRatio"`
		SeriesPrimaryImageTag   string        `json:"SeriesPrimaryImageTag"`
		SeasonName              string        `json:"SeasonName"`
		MediaStreams            []struct {
			Codec                  string  `json:"Codec"`
			TimeBase               string  `json:"TimeBase"`
			CodecTimeBase          string  `json:"CodecTimeBase"`
			VideoRange             string  `json:"VideoRange,omitempty"`
			DisplayTitle           string  `json:"DisplayTitle"`
			NalLengthSize          string  `json:"NalLengthSize,omitempty"`
			IsInterlaced           bool    `json:"IsInterlaced"`
			IsAVC                  bool    `json:"IsAVC,omitempty"`
			BitRate                int     `json:"BitRate,omitempty"`
			BitDepth               int     `json:"BitDepth,omitempty"`
			RefFrames              int     `json:"RefFrames,omitempty"`
			IsDefault              bool    `json:"IsDefault"`
			IsForced               bool    `json:"IsForced"`
			Height                 int     `json:"Height,omitempty"`
			Width                  int     `json:"Width,omitempty"`
			AverageFrameRate       float64 `json:"AverageFrameRate,omitempty"`
			RealFrameRate          float64 `json:"RealFrameRate,omitempty"`
			Profile                string  `json:"Profile,omitempty"`
			Type                   string  `json:"Type"`
			AspectRatio            string  `json:"AspectRatio,omitempty"`
			Index                  int     `json:"Index"`
			IsExternal             bool    `json:"IsExternal"`
			IsTextSubtitleStream   bool    `json:"IsTextSubtitleStream"`
			SupportsExternalStream bool    `json:"SupportsExternalStream"`
			PixelFormat            string  `json:"PixelFormat,omitempty"`
			Level                  int     `json:"Level"`
			IsAnamorphic           bool    `json:"IsAnamorphic"`
			Language               string  `json:"Language,omitempty"`
			ChannelLayout          string  `json:"ChannelLayout,omitempty"`
			Channels               int     `json:"Channels,omitempty"`
			SampleRate             int     `json:"SampleRate,omitempty"`
		} `json:"MediaStreams"`
		VideoType string `json:"VideoType"`
		ImageTags struct {
			Primary string `json:"Primary"`
		} `json:"ImageTags"`
		BackdropImageTags   []interface{} `json:"BackdropImageTags"`
		ScreenshotImageTags []interface{} `json:"ScreenshotImageTags"`
		ParentLogoImageTag  string        `json:"ParentLogoImageTag"`
		SeriesStudio        string        `json:"SeriesStudio"`
		ParentThumbItemID   string        `json:"ParentThumbItemId"`
		ParentThumbImageTag string        `json:"ParentThumbImageTag"`
		Chapters            []struct {
			StartPositionTicks int    `json:"StartPositionTicks"`
			Name               string `json:"Name"`
		} `json:"Chapters"`
		LocationType string `json:"LocationType"`
		MediaType    string `json:"MediaType"`
	} `json:"NowPlayingItem,omitempty"`
	AppIconURL      string `json:"AppIconUrl,omitempty"`
	TranscodingInfo struct {
		AudioCodec           string   `json:"AudioCodec"`
		VideoCodec           string   `json:"VideoCodec"`
		Container            string   `json:"Container"`
		IsVideoDirect        bool     `json:"IsVideoDirect"`
		IsAudioDirect        bool     `json:"IsAudioDirect"`
		Bitrate              int      `json:"Bitrate"`
		CompletionPercentage float64  `json:"CompletionPercentage"`
		Width                int      `json:"Width"`
		Height               int      `json:"Height"`
		AudioChannels        int      `json:"AudioChannels"`
		TranscodeReasons     []string `json:"TranscodeReasons"`
	} `json:"TranscodingInfo,omitempty"`
}

type ItemDetail struct {
	Name          string    `json:"Name"`
	OriginalTitle string    `json:"OriginalTitle"`
	ServerID      string    `json:"ServerId"`
	ID            string    `json:"Id"`
	Etag          string    `json:"Etag"`
	DateCreated   time.Time `json:"DateCreated"`
	CanDelete     bool      `json:"CanDelete"`
	CanDownload   bool      `json:"CanDownload"`
	HasSubtitles  bool      `json:"HasSubtitles"`
	SupportsSync  bool      `json:"SupportsSync"`
	Container     string    `json:"Container"`
	SortName      string    `json:"SortName"`
	PremiereDate  time.Time `json:"PremiereDate"`
	ExternalUrls  []struct {
		Name string `json:"Name"`
		URL  string `json:"Url"`
	} `json:"ExternalUrls"`
	MediaSources []struct {
		Protocol              string `json:"Protocol"`
		ID                    string `json:"Id"`
		Path                  string `json:"Path"`
		Type                  string `json:"Type"`
		Container             string `json:"Container"`
		Name                  string `json:"Name"`
		IsRemote              bool   `json:"IsRemote"`
		ETag                  string `json:"ETag"`
		RunTimeTicks          int64  `json:"RunTimeTicks"`
		ReadAtNativeFramerate bool   `json:"ReadAtNativeFramerate"`
		IgnoreDts             bool   `json:"IgnoreDts"`
		IgnoreIndex           bool   `json:"IgnoreIndex"`
		GenPtsInput           bool   `json:"GenPtsInput"`
		SupportsTranscoding   bool   `json:"SupportsTranscoding"`
		SupportsDirectStream  bool   `json:"SupportsDirectStream"`
		SupportsDirectPlay    bool   `json:"SupportsDirectPlay"`
		IsInfiniteStream      bool   `json:"IsInfiniteStream"`
		RequiresOpening       bool   `json:"RequiresOpening"`
		RequiresClosing       bool   `json:"RequiresClosing"`
		RequiresLooping       bool   `json:"RequiresLooping"`
		SupportsProbing       bool   `json:"SupportsProbing"`
		VideoType             string `json:"VideoType"`
		MediaStreams          []struct {
			Codec                  string  `json:"Codec"`
			Language               string  `json:"Language"`
			TimeBase               string  `json:"TimeBase"`
			CodecTimeBase          string  `json:"CodecTimeBase"`
			VideoRange             string  `json:"VideoRange,omitempty"`
			DisplayTitle           string  `json:"DisplayTitle"`
			NalLengthSize          string  `json:"NalLengthSize,omitempty"`
			IsInterlaced           bool    `json:"IsInterlaced"`
			IsAVC                  bool    `json:"IsAVC,omitempty"`
			BitRate                int     `json:"BitRate,omitempty"`
			BitDepth               int     `json:"BitDepth,omitempty"`
			RefFrames              int     `json:"RefFrames,omitempty"`
			IsDefault              bool    `json:"IsDefault"`
			IsForced               bool    `json:"IsForced"`
			Height                 int     `json:"Height,omitempty"`
			Width                  int     `json:"Width,omitempty"`
			AverageFrameRate       float64 `json:"AverageFrameRate,omitempty"`
			RealFrameRate          float64 `json:"RealFrameRate,omitempty"`
			Profile                string  `json:"Profile,omitempty"`
			Type                   string  `json:"Type"`
			AspectRatio            string  `json:"AspectRatio,omitempty"`
			Index                  int     `json:"Index"`
			IsExternal             bool    `json:"IsExternal"`
			IsTextSubtitleStream   bool    `json:"IsTextSubtitleStream"`
			SupportsExternalStream bool    `json:"SupportsExternalStream"`
			PixelFormat            string  `json:"PixelFormat,omitempty"`
			Level                  int     `json:"Level"`
			IsAnamorphic           bool    `json:"IsAnamorphic,omitempty"`
			ChannelLayout          string  `json:"ChannelLayout,omitempty"`
			Channels               int     `json:"Channels,omitempty"`
			SampleRate             int     `json:"SampleRate,omitempty"`
			Title                  string  `json:"Title,omitempty"`
		} `json:"MediaStreams"`
		Formats             []interface{} `json:"Formats"`
		Bitrate             int           `json:"Bitrate"`
		RequiredHTTPHeaders struct {
		} `json:"RequiredHttpHeaders"`
		DefaultAudioStreamIndex int `json:"DefaultAudioStreamIndex"`
	} `json:"MediaSources"`
	CriticRating             int      `json:"CriticRating"`
	ProductionLocations      []string `json:"ProductionLocations"`
	Path                     string   `json:"Path"`
	EnableMediaSourceDisplay bool     `json:"EnableMediaSourceDisplay"`
	OfficialRating           string   `json:"OfficialRating"`
	Overview                 string   `json:"Overview"`
	Taglines                 []string `json:"Taglines"`
	Genres                   []string `json:"Genres"`
	CommunityRating          float64  `json:"CommunityRating"`
	RunTimeTicks             int64    `json:"RunTimeTicks"`
	PlayAccess               string   `json:"PlayAccess"`
	ProductionYear           int      `json:"ProductionYear"`
	RemoteTrailers           []struct {
		URL  string `json:"Url"`
		Name string `json:"Name"`
	} `json:"RemoteTrailers"`
	ProviderIds struct {
		Imdb           string `json:"Imdb"`
		Tmdb           string `json:"Tmdb"`
		TmdbCollection string `json:"TmdbCollection"`
	} `json:"ProviderIds"`
	IsFolder bool   `json:"IsFolder"`
	ParentID string `json:"ParentId"`
	Type     string `json:"Type"`
	People   []struct {
		Name            string `json:"Name"`
		ID              string `json:"Id"`
		Role            string `json:"Role"`
		Type            string `json:"Type"`
		PrimaryImageTag string `json:"PrimaryImageTag,omitempty"`
	} `json:"People"`
	Studios []struct {
		Name string `json:"Name"`
		ID   string `json:"Id"`
	} `json:"Studios"`
	GenreItems []struct {
		Name string `json:"Name"`
		ID   string `json:"Id"`
	} `json:"GenreItems"`
	LocalTrailerCount int `json:"LocalTrailerCount"`
	UserData          struct {
		PlaybackPositionTicks int       `json:"PlaybackPositionTicks"`
		PlayCount             int       `json:"PlayCount"`
		IsFavorite            bool      `json:"IsFavorite"`
		LastPlayedDate        time.Time `json:"LastPlayedDate"`
		Played                bool      `json:"Played"`
		Key                   string    `json:"Key"`
	} `json:"UserData"`
	SpecialFeatureCount     int           `json:"SpecialFeatureCount"`
	DisplayPreferencesID    string        `json:"DisplayPreferencesId"`
	Tags                    []interface{} `json:"Tags"`
	PrimaryImageAspectRatio float64       `json:"PrimaryImageAspectRatio"`
	MediaStreams            []struct {
		Codec                  string  `json:"Codec"`
		Language               string  `json:"Language"`
		TimeBase               string  `json:"TimeBase"`
		CodecTimeBase          string  `json:"CodecTimeBase"`
		VideoRange             string  `json:"VideoRange,omitempty"`
		DisplayTitle           string  `json:"DisplayTitle"`
		NalLengthSize          string  `json:"NalLengthSize,omitempty"`
		IsInterlaced           bool    `json:"IsInterlaced"`
		IsAVC                  bool    `json:"IsAVC,omitempty"`
		BitRate                int     `json:"BitRate,omitempty"`
		BitDepth               int     `json:"BitDepth,omitempty"`
		RefFrames              int     `json:"RefFrames,omitempty"`
		IsDefault              bool    `json:"IsDefault"`
		IsForced               bool    `json:"IsForced"`
		Height                 int     `json:"Height,omitempty"`
		Width                  int     `json:"Width,omitempty"`
		AverageFrameRate       float64 `json:"AverageFrameRate,omitempty"`
		RealFrameRate          float64 `json:"RealFrameRate,omitempty"`
		Profile                string  `json:"Profile,omitempty"`
		Type                   string  `json:"Type"`
		AspectRatio            string  `json:"AspectRatio,omitempty"`
		Index                  int     `json:"Index"`
		IsExternal             bool    `json:"IsExternal"`
		IsTextSubtitleStream   bool    `json:"IsTextSubtitleStream"`
		SupportsExternalStream bool    `json:"SupportsExternalStream"`
		PixelFormat            string  `json:"PixelFormat,omitempty"`
		Level                  int     `json:"Level"`
		IsAnamorphic           bool    `json:"IsAnamorphic,omitempty"`
		ChannelLayout          string  `json:"ChannelLayout,omitempty"`
		Channels               int     `json:"Channels,omitempty"`
		SampleRate             int     `json:"SampleRate,omitempty"`
		Title                  string  `json:"Title,omitempty"`
	} `json:"MediaStreams"`
	VideoType string `json:"VideoType"`
	ImageTags struct {
		Primary string `json:"Primary"`
		Logo    string `json:"Logo"`
		Thumb   string `json:"Thumb"`
	} `json:"ImageTags"`
	BackdropImageTags   []string      `json:"BackdropImageTags"`
	ScreenshotImageTags []interface{} `json:"ScreenshotImageTags"`
	Chapters            []struct {
		StartPositionTicks int    `json:"StartPositionTicks"`
		Name               string `json:"Name"`
	} `json:"Chapters"`
	LocationType string        `json:"LocationType"`
	MediaType    string        `json:"MediaType"`
	LockedFields []interface{} `json:"LockedFields"`
	LockData     bool          `json:"LockData"`
}

type ItemImages struct {
	Images []struct {
		ProviderName    string  `json:"ProviderName"`
		URL             string  `json:"Url"`
		Height          int     `json:"Height"`
		Width           int     `json:"Width"`
		Type            string  `json:"Type"`
		RatingType      string  `json:"RatingType"`
		CommunityRating float64 `json:"CommunityRating,omitempty"`
		VoteCount       int     `json:"VoteCount,omitempty"`
	} `json:"Images"`
	TotalRecordCount int      `json:"TotalRecordCount"`
	Providers        []string `json:"Providers"`
}

type Client struct {
	adminID string
	token   string
	baseURL *url.URL
	http    http.Client
}

func NewClient(host, token string, adminID string) (*Client, error) {
	base, err := url.Parse(host)
	if err != nil {
		return nil, fmt.Errorf("failed to parse host %q: %v", host, err)
	}
	base.Scheme = "https"
	httpClient := http.Client{
		Timeout: httpTimeout,
	}
	return &Client{
		adminID: adminID,
		baseURL: base,
		token:   token,
		http:    httpClient,
	}, nil
}

func (c *Client) Sessions(ctx context.Context) (Sessions, error) {
	body, err := c.do(ctx, "GET", fmt.Sprintf("Sessions"))
	if err != nil {
		return nil, err
	}

	var sessions Sessions
	if err := json.Unmarshal(body, &sessions); err != nil {
		return nil, err
	}

	for i, session := range sessions {
		if session.NowPlayingItem.ID != "" {
			details, err := c.itemDetails(ctx, session.NowPlayingItem.ID)
			if err != nil {
				return nil, err
			}
			sessions[i].ItemDetail = details

			images, err := c.itemImages(ctx, session.NowPlayingItem.ID)
			if err != nil {
				return nil, err
			}
			sessions[i].ItemImages = images
		}
	}

	return sessions, nil
}

func (c *Client) itemDetails(ctx context.Context, id string) (ItemDetail, error) {
	body, err := c.do(ctx, "GET", fmt.Sprintf("Users/%s/Items/%s", c.adminID, id))
	if err != nil {
		return ItemDetail{}, err
	}

	var itemDetail ItemDetail
	if err := json.Unmarshal(body, &itemDetail); err != nil {
		return ItemDetail{}, err
	}

	return itemDetail, nil
}

func (c *Client) itemImages(ctx context.Context, id string) (ItemImages, error) {
	body, err := c.do(ctx, "GET", fmt.Sprintf("Items/%s/RemoteImages/", id))
	if err != nil {
		return ItemImages{}, err
	}

	var images ItemImages
	if err := json.Unmarshal(body, &images); err != nil {
		return ItemImages{}, err
	}

	return images, nil
}

func (c *Client) do(ctx context.Context, method string, path string) ([]byte, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.baseURL, path), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-MediaBrowser-Token", c.token)
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

func (c *Client) Search(ctx context.Context, searchTerm string) ([]byte, error) {
	return nil, nil
}
