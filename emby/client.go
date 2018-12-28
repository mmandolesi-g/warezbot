package emby

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	httpTimeout = 10 * time.Second
)

type Sessions []struct {
	ItemDetail ItemDetail
	ItemImages ItemImages
	PlayState  struct {
		PositionTicks    int64  `json:"PositionTicks"`
		CanSeek          bool   `json:"CanSeek"`
		IsPaused         bool   `json:"IsPaused"`
		IsMuted          bool   `json:"IsMuted"`
		VolumeLevel      int    `json:"VolumeLevel"`
		AudioStreamIndex int    `json:"AudioStreamIndex"`
		MediaSourceID    string `json:"MediaSourceId"`
		PlayMethod       string `json:"PlayMethod"`
		RepeatMode       string `json:"RepeatMode"`
	} `json:"PlayState"`
	AdditionalUsers []interface{} `json:"AdditionalUsers"`
	Capabilities    struct {
		ID                           string   `json:"Id"`
		PlayableMediaTypes           []string `json:"PlayableMediaTypes"`
		SupportedCommands            []string `json:"SupportedCommands"`
		SupportsMediaControl         bool     `json:"SupportsMediaControl"`
		SupportsPersistentIdentifier bool     `json:"SupportsPersistentIdentifier"`
		SupportsSync                 bool     `json:"SupportsSync"`
		DeviceProfile                struct {
			EnableAlbumArtInDidl             bool          `json:"EnableAlbumArtInDidl"`
			EnableSingleAlbumArtLimit        bool          `json:"EnableSingleAlbumArtLimit"`
			EnableSingleSubtitleLimit        bool          `json:"EnableSingleSubtitleLimit"`
			SupportedMediaTypes              string        `json:"SupportedMediaTypes"`
			MaxAlbumArtWidth                 int           `json:"MaxAlbumArtWidth"`
			MaxAlbumArtHeight                int           `json:"MaxAlbumArtHeight"`
			MaxStreamingBitrate              int           `json:"MaxStreamingBitrate"`
			MaxStaticBitrate                 int           `json:"MaxStaticBitrate"`
			MusicStreamingTranscodingBitrate int           `json:"MusicStreamingTranscodingBitrate"`
			MaxStaticMusicBitrate            int           `json:"MaxStaticMusicBitrate"`
			TimelineOffsetSeconds            int           `json:"TimelineOffsetSeconds"`
			RequiresPlainVideoItems          bool          `json:"RequiresPlainVideoItems"`
			RequiresPlainFolders             bool          `json:"RequiresPlainFolders"`
			EnableMSMediaReceiverRegistrar   bool          `json:"EnableMSMediaReceiverRegistrar"`
			IgnoreTranscodeByteRangeRequests bool          `json:"IgnoreTranscodeByteRangeRequests"`
			XMLRootAttributes                []interface{} `json:"XmlRootAttributes"`
			DirectPlayProfiles               []struct {
				Container  string `json:"Container"`
				AudioCodec string `json:"AudioCodec,omitempty"`
				VideoCodec string `json:"VideoCodec,omitempty"`
				Type       string `json:"Type"`
			} `json:"DirectPlayProfiles"`
			TranscodingProfiles []struct {
				Container             string `json:"Container"`
				Type                  string `json:"Type"`
				AudioCodec            string `json:"AudioCodec"`
				Protocol              string `json:"Protocol,omitempty"`
				EstimateContentLength bool   `json:"EstimateContentLength"`
				EnableMpegtsM2TsMode  bool   `json:"EnableMpegtsM2TsMode"`
				TranscodeSeekInfo     string `json:"TranscodeSeekInfo"`
				CopyTimestamps        bool   `json:"CopyTimestamps"`
				Context               string `json:"Context"`
				MaxAudioChannels      string `json:"MaxAudioChannels,omitempty"`
				MinSegments           int    `json:"MinSegments"`
				SegmentLength         int    `json:"SegmentLength"`
				BreakOnNonKeyFrames   bool   `json:"BreakOnNonKeyFrames"`
				VideoCodec            string `json:"VideoCodec,omitempty"`
			} `json:"TranscodingProfiles"`
			ContainerProfiles []interface{} `json:"ContainerProfiles"`
			CodecProfiles     []struct {
				Type       string `json:"Type"`
				Conditions []struct {
					Condition  string `json:"Condition"`
					Property   string `json:"Property"`
					Value      string `json:"Value"`
					IsRequired bool   `json:"IsRequired"`
				} `json:"Conditions"`
				ApplyConditions []interface{} `json:"ApplyConditions"`
				Codec           string        `json:"Codec,omitempty"`
			} `json:"CodecProfiles"`
			ResponseProfiles []struct {
				Container  string        `json:"Container"`
				Type       string        `json:"Type"`
				MimeType   string        `json:"MimeType"`
				Conditions []interface{} `json:"Conditions"`
			} `json:"ResponseProfiles"`
			SubtitleProfiles []struct {
				Format string `json:"Format"`
				Method string `json:"Method"`
			} `json:"SubtitleProfiles"`
		} `json:"DeviceProfile"`
		IconURL string `json:"IconUrl"`
	} `json:"Capabilities"`
	RemoteEndPoint      string    `json:"RemoteEndPoint"`
	PlayableMediaTypes  []string  `json:"PlayableMediaTypes"`
	PlaylistItemID      string    `json:"PlaylistItemId,omitempty"`
	ID                  string    `json:"Id"`
	ServerID            string    `json:"ServerId"`
	UserID              string    `json:"UserId"`
	UserName            string    `json:"UserName"`
	UserPrimaryImageTag string    `json:"UserPrimaryImageTag,omitempty"`
	Client              string    `json:"Client"`
	LastActivityDate    time.Time `json:"LastActivityDate"`
	DeviceName          string    `json:"DeviceName"`
	NowPlayingItem      struct {
		Name          string    `json:"Name"`
		OriginalTitle string    `json:"OriginalTitle"`
		ServerID      string    `json:"ServerId"`
		ID            string    `json:"Id"`
		DateCreated   time.Time `json:"DateCreated"`
		Container     string    `json:"Container"`
		PremiereDate  time.Time `json:"PremiereDate"`
		ExternalUrls  []struct {
			Name string `json:"Name"`
			URL  string `json:"Url"`
		} `json:"ExternalUrls"`
		CriticRating    int      `json:"CriticRating"`
		Path            string   `json:"Path"`
		OfficialRating  string   `json:"OfficialRating"`
		Overview        string   `json:"Overview"`
		Taglines        []string `json:"Taglines"`
		Genres          []string `json:"Genres"`
		CommunityRating float64  `json:"CommunityRating"`
		RunTimeTicks    int64    `json:"RunTimeTicks"`
		ProductionYear  int      `json:"ProductionYear"`
		ProviderIds     struct {
			Tmdb string `json:"Tmdb"`
			Imdb string `json:"Imdb"`
		} `json:"ProviderIds"`
		IsFolder bool   `json:"IsFolder"`
		ParentID string `json:"ParentId"`
		Type     string `json:"Type"`
		Studios  []struct {
			Name string `json:"Name"`
			ID   int    `json:"Id"`
		} `json:"Studios"`
		GenreItems []struct {
			Name string `json:"Name"`
			ID   int    `json:"Id"`
		} `json:"GenreItems"`
		SeriesName              string  `json:"SeriesName"`
		SeriesId                string  `json:"SeriesId"`
		SeasonId                string  `json:"SeasonId"`
		IndexNumber             int     `json:"IndexNumber"`
		ParentIndexNumber       int     `json:"ParentIndexNumber"`
		LocalTrailerCount       int     `json:"LocalTrailerCount"`
		PrimaryImageAspectRatio float64 `json:"PrimaryImageAspectRatio"`
		MediaStreams            []struct {
			Codec                  string  `json:"Codec"`
			Language               string  `json:"Language"`
			TimeBase               string  `json:"TimeBase"`
			CodecTimeBase          string  `json:"CodecTimeBase"`
			DisplayTitle           string  `json:"DisplayTitle"`
			IsInterlaced           bool    `json:"IsInterlaced"`
			ChannelLayout          string  `json:"ChannelLayout,omitempty"`
			BitRate                int     `json:"BitRate"`
			Channels               int     `json:"Channels,omitempty"`
			SampleRate             int     `json:"SampleRate,omitempty"`
			IsDefault              bool    `json:"IsDefault"`
			IsForced               bool    `json:"IsForced"`
			Type                   string  `json:"Type"`
			Index                  int     `json:"Index"`
			IsExternal             bool    `json:"IsExternal"`
			IsTextSubtitleStream   bool    `json:"IsTextSubtitleStream"`
			SupportsExternalStream bool    `json:"SupportsExternalStream"`
			Level                  int     `json:"Level"`
			ColorTransfer          string  `json:"ColorTransfer,omitempty"`
			ColorPrimaries         string  `json:"ColorPrimaries,omitempty"`
			ColorSpace             string  `json:"ColorSpace,omitempty"`
			VideoRange             string  `json:"VideoRange,omitempty"`
			NalLengthSize          string  `json:"NalLengthSize,omitempty"`
			IsAVC                  bool    `json:"IsAVC,omitempty"`
			BitDepth               int     `json:"BitDepth,omitempty"`
			RefFrames              int     `json:"RefFrames,omitempty"`
			Height                 int     `json:"Height,omitempty"`
			Width                  int     `json:"Width,omitempty"`
			AverageFrameRate       float64 `json:"AverageFrameRate,omitempty"`
			RealFrameRate          float64 `json:"RealFrameRate,omitempty"`
			Profile                string  `json:"Profile,omitempty"`
			AspectRatio            string  `json:"AspectRatio,omitempty"`
			PixelFormat            string  `json:"PixelFormat,omitempty"`
			IsAnamorphic           bool    `json:"IsAnamorphic,omitempty"`
		} `json:"MediaStreams"`
		ImageTags struct {
			Primary string `json:"Primary"`
			Logo    string `json:"Logo"`
		} `json:"ImageTags"`
		BackdropImageTags []string `json:"BackdropImageTags"`
		Chapters          []struct {
			StartPositionTicks int    `json:"StartPositionTicks"`
			Name               string `json:"Name"`
		} `json:"Chapters"`
		MediaType string `json:"MediaType"`
		Width     int    `json:"Width"`
		Height    int    `json:"Height"`
	} `json:"NowPlayingItem,omitempty"`
	DeviceID           string   `json:"DeviceId"`
	ApplicationVersion string   `json:"ApplicationVersion"`
	AppIconURL         string   `json:"AppIconUrl,omitempty"`
	SupportedCommands  []string `json:"SupportedCommands"`
	TranscodingInfo    struct {
		AudioCodec                    string   `json:"AudioCodec"`
		VideoCodec                    string   `json:"VideoCodec"`
		Container                     string   `json:"Container"`
		IsVideoDirect                 bool     `json:"IsVideoDirect"`
		IsAudioDirect                 bool     `json:"IsAudioDirect"`
		Bitrate                       int      `json:"Bitrate"`
		Framerate                     int      `json:"Framerate"`
		CompletionPercentage          float64  `json:"CompletionPercentage"`
		TranscodingPositionTicks      int64    `json:"TranscodingPositionTicks"`
		TranscodingStartPositionTicks int      `json:"TranscodingStartPositionTicks"`
		Width                         int      `json:"Width"`
		Height                        int      `json:"Height"`
		AudioChannels                 int      `json:"AudioChannels"`
		TranscodeReasons              []string `json:"TranscodeReasons"`
		CurrentThrottle               int      `json:"CurrentThrottle"`
		VideoDecoderIsHardware        bool     `json:"VideoDecoderIsHardware"`
		VideoEncoderIsHardware        bool     `json:"VideoEncoderIsHardware"`
	} `json:"TranscodingInfo,omitempty"`
	SupportsRemoteControl bool `json:"SupportsRemoteControl"`
}

type ItemDetail struct {
	Name         string    `json:"Name"`
	ServerID     string    `json:"ServerId"`
	ID           string    `json:"Id"`
	Etag         string    `json:"Etag"`
	DateCreated  time.Time `json:"DateCreated"`
	CanDelete    bool      `json:"CanDelete"`
	CanDownload  bool      `json:"CanDownload"`
	HasSubtitles bool      `json:"HasSubtitles"`
	SupportsSync bool      `json:"SupportsSync"`
	Container    string    `json:"Container"`
	SortName     string    `json:"SortName"`
	PremiereDate time.Time `json:"PremiereDate"`
	ExternalUrls []struct {
		Name string `json:"Name"`
		URL  string `json:"Url"`
	} `json:"ExternalUrls"`
	MediaSources []struct {
		Protocol              string `json:"Protocol"`
		ID                    string `json:"Id"`
		Path                  string `json:"Path"`
		Type                  string `json:"Type"`
		Container             string `json:"Container"`
		Size                  int    `json:"Size"`
		Name                  string `json:"Name"`
		IsRemote              bool   `json:"IsRemote"`
		RunTimeTicks          int64  `json:"RunTimeTicks"`
		ReadAtNativeFramerate bool   `json:"ReadAtNativeFramerate"`
		DiscardCorruptPts     bool   `json:"DiscardCorruptPts"`
		FillWallClockDts      bool   `json:"FillWallClockDts"`
		IgnoreDts             bool   `json:"IgnoreDts"`
		IgnoreIndex           bool   `json:"IgnoreIndex"`
		SupportsTranscoding   bool   `json:"SupportsTranscoding"`
		SupportsDirectStream  bool   `json:"SupportsDirectStream"`
		SupportsDirectPlay    bool   `json:"SupportsDirectPlay"`
		IsInfiniteStream      bool   `json:"IsInfiniteStream"`
		RequiresOpening       bool   `json:"RequiresOpening"`
		RequiresClosing       bool   `json:"RequiresClosing"`
		RequiresLooping       bool   `json:"RequiresLooping"`
		SupportsProbing       bool   `json:"SupportsProbing"`
		MediaStreams          []struct {
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
			Language               string  `json:"Language,omitempty"`
			ChannelLayout          string  `json:"ChannelLayout,omitempty"`
			Channels               int     `json:"Channels,omitempty"`
			SampleRate             int     `json:"SampleRate,omitempty"`
		} `json:"MediaStreams"`
		Formats             []interface{} `json:"Formats"`
		Bitrate             int           `json:"Bitrate"`
		RequiredHTTPHeaders struct {
		} `json:"RequiredHttpHeaders"`
		DefaultAudioStreamIndex int `json:"DefaultAudioStreamIndex"`
	} `json:"MediaSources"`
	Path              string        `json:"Path"`
	Overview          string        `json:"Overview"`
	Taglines          []interface{} `json:"Taglines"`
	Genres            []interface{} `json:"Genres"`
	RunTimeTicks      int64         `json:"RunTimeTicks"`
	PlayAccess        string        `json:"PlayAccess"`
	ProductionYear    int           `json:"ProductionYear"`
	IndexNumber       int           `json:"IndexNumber"`
	ParentIndexNumber int           `json:"ParentIndexNumber"`
	RemoteTrailers    []interface{} `json:"RemoteTrailers"`
	ProviderIds       struct {
		Tvdb string `json:"Tvdb"`
		Imdb string `json:"Imdb"`
	} `json:"ProviderIds"`
	IsFolder bool   `json:"IsFolder"`
	ParentID string `json:"ParentId"`
	Type     string `json:"Type"`
	People   []struct {
		Name string `json:"Name"`
		ID   string `json:"Id"`
		Type string `json:"Type"`
	} `json:"People"`
	Studios                 []interface{} `json:"Studios"`
	GenreItems              []interface{} `json:"GenreItems"`
	ParentLogoItemID        string        `json:"ParentLogoItemId"`
	ParentBackdropItemID    string        `json:"ParentBackdropItemId"`
	ParentBackdropImageTags []string      `json:"ParentBackdropImageTags"`
	UserData                struct {
		PlaybackPositionTicks int    `json:"PlaybackPositionTicks"`
		PlayCount             int    `json:"PlayCount"`
		IsFavorite            bool   `json:"IsFavorite"`
		Played                bool   `json:"Played"`
		Key                   string `json:"Key"`
	} `json:"UserData"`
	SeriesName              string        `json:"SeriesName"`
	SeriesID                string        `json:"SeriesId"`
	SeasonID                string        `json:"SeasonId"`
	DisplayPreferencesID    string        `json:"DisplayPreferencesId"`
	Tags                    []interface{} `json:"Tags"`
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
		Language               string  `json:"Language,omitempty"`
		ChannelLayout          string  `json:"ChannelLayout,omitempty"`
		Channels               int     `json:"Channels,omitempty"`
		SampleRate             int     `json:"SampleRate,omitempty"`
	} `json:"MediaStreams"`
	ImageTags struct {
		Primary string `json:"Primary"`
	} `json:"ImageTags"`
	BackdropImageTags   []interface{} `json:"BackdropImageTags"`
	ParentLogoImageTag  string        `json:"ParentLogoImageTag"`
	ParentThumbItemID   string        `json:"ParentThumbItemId"`
	ParentThumbImageTag string        `json:"ParentThumbImageTag"`
	Chapters            []struct {
		StartPositionTicks int    `json:"StartPositionTicks"`
		Name               string `json:"Name"`
	} `json:"Chapters"`
	MediaType    string        `json:"MediaType"`
	LockedFields []interface{} `json:"LockedFields"`
	LockData     bool          `json:"LockData"`
	Width        int           `json:"Width"`
	Height       int           `json:"Height"`
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

type SearchResults struct {
	SearchHints []struct {
		ItemImages              ItemImages
		ItemDetail              ItemDetail
		ItemID                  int      `json:"ItemId"`
		ID                      int      `json:"Id"`
		Name                    string   `json:"Name"`
		IndexNumber             int      `json:"IndexNumber,omitempty"`
		ProductionYear          int      `json:"ProductionYear,omitempty"`
		PrimaryImageTag         string   `json:"PrimaryImageTag,omitempty"`
		Type                    string   `json:"Type"`
		RunTimeTicks            int64    `json:"RunTimeTicks,omitempty"`
		MediaType               string   `json:"MediaType,omitempty"`
		Album                   string   `json:"Album,omitempty"`
		AlbumID                 int      `json:"AlbumId"`
		AlbumArtist             string   `json:"AlbumArtist,omitempty"`
		Artists                 []string `json:"Artists,omitempty"`
		PrimaryImageAspectRatio float64  `json:"PrimaryImageAspectRatio,omitempty"`
		ParentIndexNumber       int      `json:"ParentIndexNumber,omitempty"`
		ThumbImageTag           string   `json:"ThumbImageTag,omitempty"`
		ThumbImageItemID        string   `json:"ThumbImageItemId,omitempty"`
		BackdropImageTag        string   `json:"BackdropImageTag,omitempty"`
		BackdropImageItemID     string   `json:"BackdropImageItemId,omitempty"`
		Series                  string   `json:"Series,omitempty"`
	} `json:"SearchHints"`
	TotalRecordCount int `json:"TotalRecordCount"`
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

func (c *Client) Search(ctx context.Context, searchTerm []string) (SearchResults, error) {
	s := strings.Join(searchTerm, " ")
	body, err := c.do(ctx, "GET", fmt.Sprintf("Search/Hints?searchTerm=%s", s))
	if err != nil {
		fmt.Println(err)
	}

	var sr SearchResults
	if err := json.Unmarshal(body, &sr); err != nil {
		fmt.Println(err)
	}

	for i, hint := range sr.SearchHints {
		if hint.ID != 0 {
			images, err := c.itemImages(ctx, strconv.Itoa(hint.ID))
			if err != nil {
				fmt.Println(err)
			}
			sr.SearchHints[i].ItemImages = images

			details, err := c.itemDetails(ctx, strconv.Itoa(hint.ID))
			if err != nil {
				fmt.Println(err)
			}
			sr.SearchHints[i].ItemDetail = details
		}
	}
	return sr, nil
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
