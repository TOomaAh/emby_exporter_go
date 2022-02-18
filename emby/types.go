package emby

import "time"

//SystemInfo return all server info
type SystemInfo struct {
	SystemUpdateLevel                    string        `json:"SystemUpdateLevel"`
	OperatingSystemDisplayName           string        `json:"OperatingSystemDisplayName"`
	HasPendingRestart                    bool          `json:"HasPendingRestart"`
	IsShuttingDown                       bool          `json:"IsShuttingDown"`
	SupportsLibraryMonitor               bool          `json:"SupportsLibraryMonitor"`
	WebSocketPortNumber                  int           `json:"WebSocketPortNumber"`
	CompletedInstallations               []interface{} `json:"CompletedInstallations"`
	CanSelfRestart                       bool          `json:"CanSelfRestart"`
	CanSelfUpdate                        bool          `json:"CanSelfUpdate"`
	CanLaunchWebBrowser                  bool          `json:"CanLaunchWebBrowser"`
	ProgramDataPath                      string        `json:"ProgramDataPath"`
	ItemsByNamePath                      string        `json:"ItemsByNamePath"`
	CachePath                            string        `json:"CachePath"`
	LogPath                              string        `json:"LogPath"`
	InternalMetadataPath                 string        `json:"InternalMetadataPath"`
	TranscodingTempPath                  string        `json:"TranscodingTempPath"`
	HTTPServerPortNumber                 int           `json:"HttpServerPortNumber"`
	SupportsHTTPS                        bool          `json:"SupportsHttps"`
	HTTPSPortNumber                      int           `json:"HttpsPortNumber"`
	HasUpdateAvailable                   bool          `json:"HasUpdateAvailable"`
	SupportsAutoRunAtStartup             bool          `json:"SupportsAutoRunAtStartup"`
	HardwareAccelerationRequiresPremiere bool          `json:"HardwareAccelerationRequiresPremiere"`
	LocalAddress                         string        `json:"LocalAddress"`
	WanAddress                           string        `json:"WanAddress"`
	ServerName                           string        `json:"ServerName"`
	Version                              string        `json:"Version"`
	OperatingSystem                      string        `json:"OperatingSystem"`
	ID                                   string        `json:"Id"`
}

// UserView represents a single user-visible view
type UserView struct {
	Name                    string        `json:"Name"`
	ServerID                string        `json:"ServerId"`
	ID                      string        `json:"Id"`
	Etag                    string        `json:"Etag"`
	DateCreated             string        `json:"DateCreated"`
	CanDelete               bool          `json:"CanDelete"`
	CanDownload             bool          `json:"CanDownload"`
	SortName                string        `json:"SortName"`
	ExternalUrls            []string      `json:"ExternalUrls"`
	Path                    string        `json:"Path"`
	Taglines                []string      `json:"Taglines"`
	Genres                  []string      `json:"Genres"`
	PlayAccess              string        `json:"PlayAccess"`
	RemoteTrailers          []interface{} `json:"RemoteTrailers"`
	ProviderIds             interface{}   `json:"ProviderIds"`
	IsFolder                bool          `json:"IsFolder"`
	ParentID                string        `json:"ParentId"`
	Type                    string        `json:"Type"`
	Studios                 []interface{} `json:"Studios"`
	GenreItems              []interface{} `json:"GenreItems"`
	UserData                UserData      `json:"UserData"`
	ChildCount              int           `json:"ChildCount"`
	DisplayPreferencesID    string        `json:"DisplayPreferencesId"`
	Tags                    []string      `json:"Tags"`
	PrimaryImageAspectRatio float32       `json:"PrimaryImageAspectRatio,omitempty"`
	CollectionType          string        `json:"CollectionType"`
	ImageTags               ImageTags     `json:"ImageTags"`
	BackdropImageTags       []interface{} `json:"BackdropImageTags"`
	LockedFields            []interface{} `json:"LockedFields"`
	LockData                bool          `json:"LockData"`
}

type LibraryInfo struct {
	Name      string `json:"Name"`
	ServerID  string `json:"ServerId"`
	ID        string `json:"Id"`
	IsFolder  bool   `json:"IsFolder"`
	Type      string `json:"Type"`
	ImageTags struct {
	} `json:"ImageTags"`
	BackdropImageTags []interface{} `json:"BackdropImageTags"`
}

type Library struct {
	Items            []*Item `json:"Items"`
	TotalRecordCount int     `json:"TotalRecordCount"`
}

type Item struct {
	Name         string `json:"Name"`
	ServerID     string `json:"ServerId"`
	ID           string `json:"Id"`
	SupportsSync bool   `json:"SupportsSync"`
	RunTimeTicks int64  `json:"RunTimeTicks"`
	IsFolder     bool   `json:"IsFolder"`
	Type         string `json:"Type"`
	UserData     struct {
		PlaybackPositionTicks int       `json:"PlaybackPositionTicks"`
		PlayCount             int       `json:"PlayCount"`
		IsFavorite            bool      `json:"IsFavorite"`
		LastPlayedDate        time.Time `json:"LastPlayedDate"`
		Played                bool      `json:"Played"`
	} `json:"UserData"`
	ImageTags struct {
		Primary string `json:"Primary"`
	} `json:"ImageTags"`
	MediaType string `json:"MediaType"`
}

type Sessions struct {
	PlayState struct {
		CanSeek             bool   `json:"CanSeek"`
		IsPaused            bool   `json:"IsPaused"`
		IsMuted             bool   `json:"IsMuted"`
		PositionTicks       int64  `json:"PositionTicks"`
		VolumeLevel         int    `json:"VolumeLevel"`
		AudioStreamIndex    int    `json:"AudioStreamIndex"`
		SubtitleStreamIndex int    `json:"SubtitleStreamIndex"`
		MediaSourceID       string `json:"MediaSourceId"`
		PlayMethod          string `json:"PlayMethod"`
		RepeatMode          string `json:"RepeatMode"`
		SubtitleOffset      int    `json:"SubtitleOffset"`
		PlaybackRate        int    `json:"PlaybackRate"`
	} `json:"PlayState,omitempty"`
	AdditionalUsers       []interface{} `json:"AdditionalUsers"`
	RemoteEndPoint        string        `json:"RemoteEndPoint"`
	PlayableMediaTypes    []string      `json:"PlayableMediaTypes"`
	PlaylistIndex         int           `json:"PlaylistIndex"`
	PlaylistLength        int           `json:"PlaylistLength"`
	ID                    string        `json:"Id"`
	ServerID              string        `json:"ServerId"`
	UserID                string        `json:"UserId,omitempty"`
	UserName              string        `json:"UserName,omitempty"`
	Client                string        `json:"Client"`
	LastActivityDate      time.Time     `json:"LastActivityDate"`
	DeviceName            string        `json:"DeviceName"`
	DeviceID              string        `json:"DeviceId"`
	ApplicationVersion    string        `json:"ApplicationVersion"`
	AppIconURL            string        `json:"AppIconUrl,omitempty"`
	SupportedCommands     []string      `json:"SupportedCommands"`
	SupportsRemoteControl bool          `json:"SupportsRemoteControl"`
	NowPlayingItem        struct {
		Name                  string    `json:"Name"`
		OriginalTitle         string    `json:"OriginalTitle"`
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
		Path            string        `json:"Path"`
		OfficialRating  string        `json:"OfficialRating"`
		Overview        string        `json:"Overview"`
		Taglines        []interface{} `json:"Taglines"`
		Genres          []string      `json:"Genres"`
		CommunityRating float64       `json:"CommunityRating"`
		RunTimeTicks    int64         `json:"RunTimeTicks"`
		ProductionYear  int           `json:"ProductionYear"`
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
		LocalTrailerCount       int     `json:"LocalTrailerCount"`
		PrimaryImageAspectRatio float64 `json:"PrimaryImageAspectRatio"`
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
			Profile                string  `json:"Profile,omitempty"`
			Type                   string  `json:"Type"`
			AspectRatio            string  `json:"AspectRatio,omitempty"`
			Index                  int     `json:"Index"`
			IsExternal             bool    `json:"IsExternal"`
			IsTextSubtitleStream   bool    `json:"IsTextSubtitleStream"`
			SupportsExternalStream bool    `json:"SupportsExternalStream"`
			Protocol               string  `json:"Protocol"`
			PixelFormat            string  `json:"PixelFormat,omitempty"`
			Level                  int     `json:"Level,omitempty"`
			IsAnamorphic           bool    `json:"IsAnamorphic,omitempty"`
			Language               string  `json:"Language,omitempty"`
			Title                  string  `json:"Title,omitempty"`
			DisplayLanguage        string  `json:"DisplayLanguage,omitempty"`
			ChannelLayout          string  `json:"ChannelLayout,omitempty"`
			Channels               int     `json:"Channels,omitempty"`
			SampleRate             int     `json:"SampleRate,omitempty"`
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
	TranscodingInfo struct {
		AudioCodec                    string   `json:"AudioCodec"`
		VideoCodec                    string   `json:"VideoCodec"`
		SubProtocol                   string   `json:"SubProtocol"`
		Container                     string   `json:"Container"`
		IsVideoDirect                 bool     `json:"IsVideoDirect"`
		IsAudioDirect                 bool     `json:"IsAudioDirect"`
		Bitrate                       int      `json:"Bitrate"`
		Framerate                     int      `json:"Framerate"`
		CompletionPercentage          float64  `json:"CompletionPercentage"`
		TranscodingPositionTicks      int64    `json:"TranscodingPositionTicks"`
		TranscodingStartPositionTicks int64    `json:"TranscodingStartPositionTicks"`
		Width                         int      `json:"Width"`
		Height                        int      `json:"Height"`
		AudioChannels                 int      `json:"AudioChannels"`
		TranscodeReasons              []string `json:"TranscodeReasons"`
		CurrentCPUUsage               float64  `json:"CurrentCpuUsage"`
		AverageCPUUsage               float64  `json:"AverageCpuUsage"`
		CPUHistory                    []struct {
			Item1 float64 `json:"Item1"`
			Item2 float64 `json:"Item2"`
		} `json:"CpuHistory"`
		CurrentThrottle        int  `json:"CurrentThrottle"`
		VideoDecoderIsHardware bool `json:"VideoDecoderIsHardware"`
		VideoEncoderIsHardware bool `json:"VideoEncoderIsHardware"`
		VideoPipelineInfo      []struct {
			HardwareContextName string `json:"HardwareContextName"`
			IsHardwareContext   bool   `json:"IsHardwareContext"`
			Name                string `json:"Name"`
			Short               string `json:"Short"`
			StepType            string `json:"StepType"`
			StepTypeName        string `json:"StepTypeName"`
			FfmpegName          string `json:"FfmpegName,omitempty"`
			FfmpegDescription   string `json:"FfmpegDescription,omitempty"`
			FfmpegOptions       string `json:"FfmpegOptions,omitempty"`
			Param               string `json:"Param"`
			ParamShort          string `json:"ParamShort"`
		} `json:"VideoPipelineInfo"`
	} `json:"TranscodingInfo,omitempty"`
}

// UserData is user-specific data for that media item
type UserData struct {
	PlaybackPositionTicks int    `json:"PlaybackPositionTicks"`
	PlayCount             int    `json:"PlayCount"`
	IsFavorite            bool   `json:"IsFavorite"`
	Played                bool   `json:"Played"`
	Key                   string `json:"Key"`
}

// ImageTags are image tagging details for a media item
type ImageTags struct {
	Primary string `json:"Primary"`
	Logo    string `json:"Logo"`
	Thumb   string `json:"Thumb"`
}

type MediaItemList struct {
	Items []UserView `json:"Items"`
}
