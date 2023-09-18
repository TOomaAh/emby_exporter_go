package entity

import "strings"

type TranscodeReasons string
type PlayMethod string

// TranscodeReasons ContainerNotSupported = "Container Not Supported"
const (
	ContainerNotSupported        TranscodeReasons = "Container Not Supported"
	VideoCodecNotSupported       TranscodeReasons = "Video Codec Not Supported"
	AudioCodecNotSupported       TranscodeReasons = "Audio Codec Not Supported"
	ContainerBitrateExceedsLimit TranscodeReasons = "Container Bitrate Exceeds Limit"
	AudioBitrateNotSupported     TranscodeReasons = "Audio Bitrate Not Supported"
	AudioChannelsNotSupported    TranscodeReasons = "Audio Channels Not Supported"
	VideoResolutionNotSupported  TranscodeReasons = "Video Resolution Not Supported"
	UnknownVideoStreamInfo       TranscodeReasons = "Unknown Video Stream Info"
	UnknownAudioStreamInfo       TranscodeReasons = "Unknown Audio Stream Info"
	AudioProfileNotSupported     TranscodeReasons = "Audio Profile Not Supported"
	AudioSampleRateNotSupported  TranscodeReasons = "Audio Sample Rate Not Supported"
	AnamorphicVideoNotSupported  TranscodeReasons = "Anamorphic Video Not Supported"
	InterlacedVideoNotSupported  TranscodeReasons = "Interlaced Video Not Supported"
	SecondaryAudioNotSupported   TranscodeReasons = "Secondary Audio Not Supported"
	RefFramesNotSupported        TranscodeReasons = "Ref Frames Not Supported"
	VideoBitDepthNotSupported    TranscodeReasons = "Video Bit Depth Not Supported"
	VideoBitrateNotSupported     TranscodeReasons = "Video Bitrate Not Supported"
	VideoFramerateNotSupported   TranscodeReasons = "Video Framerate Not Supported"
	VideoLevelNotSupported       TranscodeReasons = "Video Level Not Supported"
	VideoProfileNotSupported     TranscodeReasons = "Video Profile Not Supported"
	AudioBitDepthNotSupported    TranscodeReasons = "Audio Bit Depth Not Supported"
	SubtitleCodecNotSupported    TranscodeReasons = "Subtitle Codec Not Supported"
	DirectPlayError              TranscodeReasons = "Direct Play Error"
	VideoRangeNotSupported       TranscodeReasons = "Video Range Not Supported"
)

var transcodeReasons = map[string]TranscodeReasons{
	"ContainerNotSupported":        ContainerNotSupported,
	"VideoCodecNotSupported":       VideoCodecNotSupported,
	"AudioCodecNotSupported":       AudioCodecNotSupported,
	"ContainerBitrateExceedsLimit": ContainerBitrateExceedsLimit,
	"AudioBitrateNotSupported":     AudioBitrateNotSupported,
	"AudioChannelsNotSupported":    AudioChannelsNotSupported,
	"VideoResolutionNotSupported":  VideoResolutionNotSupported,
	"UnknownVideoStreamInfo":       UnknownVideoStreamInfo,
	"UnknownAudioStreamInfo":       UnknownAudioStreamInfo,
	"AudioProfileNotSupported":     AudioProfileNotSupported,
	"AudioSampleRateNotSupported":  AudioSampleRateNotSupported,
	"AnamorphicVideoNotSupported":  AnamorphicVideoNotSupported,
	"InterlacedVideoNotSupported":  InterlacedVideoNotSupported,
	"SecondaryAudioNotSupported":   SecondaryAudioNotSupported,
	"RefFramesNotSupported":        RefFramesNotSupported,
	"VideoBitDepthNotSupported":    VideoBitDepthNotSupported,
	"VideoBitrateNotSupported":     VideoBitrateNotSupported,
	"VideoFramerateNotSupported":   VideoFramerateNotSupported,
	"VideoLevelNotSupported":       VideoLevelNotSupported,
	"VideoProfileNotSupported":     VideoProfileNotSupported,
	"AudioBitDepthNotSupported":    AudioBitDepthNotSupported,
	"SubtitleCodecNotSupported":    SubtitleCodecNotSupported,
	"DirectPlayError":              DirectPlayError,
	"VideoRangeNotSupported":       VideoRangeNotSupported,
}

const (
	Transcoding PlayMethod = "Transcoding"
	DirectPlay  PlayMethod = "Direct Play"
)

type PlayState struct {
	IsPaused      bool   `json:"IsPaused"`
	PositionTicks int64  `json:"PositionTicks"`
	PlayMethod    string `json:"PlayMethod"`
}

type NowPlayingItem struct {
	Name         string `json:"Name"`
	RunTimeTicks int64  `json:"RunTimeTicks"`
	SeriesName   string `json:"SeriesName"`
	SeasonName   string `json:"SeasonName"`
	MediaType    string `json:"MediaType"`
	Type         string `json:"Type"`
}

type TranscodingInfo struct {
	AudioCodec                    string             `json:"AudioCodec"`
	VideoCodec                    string             `json:"VideoCodec"`
	IsVideoDirect                 bool               `json:"IsVideoDirect"`
	IsAudioDirect                 bool               `json:"IsAudioDirect"`
	Bitrate                       int                `json:"Bitrate"`
	TranscodingPositionTicks      int64              `json:"TranscodingPositionTicks"`
	TranscodingStartPositionTicks int64              `json:"TranscodingStartPositionTicks"`
	TranscodeReasons              []TranscodeReasons `json:"TranscodeReasons"`
	CurrentCPUUsage               float64            `json:"CurrentCpuUsage"`
	CurrentThrottle               int                `json:"CurrentThrottle"`
	VideoDecoderIsHardware        bool               `json:"VideoDecoderIsHardware"`
	VideoEncoderIsHardware        bool               `json:"VideoEncoderIsHardware"`
}

type Sessions struct {
	NowPlayingItem  *NowPlayingItem  `json:"NowPlayingItem,omitempty"`
	TranscodingInfo *TranscodingInfo `json:"TranscodingInfo,omitempty"`
	PlayState       PlayState        `json:"PlayState,omitempty"`
	RemoteEndPoint  string           `json:"RemoteEndPoint"`
	UserName        string           `json:"UserName,omitempty"`
	Client          string           `json:"Client"`
}

type SessionsMetrics struct {
	Username           string
	Client             string
	RemoteEndPoint     string
	Region             string
	City               string
	CountryCode        string
	NowPlayingItemName string
	NowPlayingItemType string
	TVShow             string
	Season             string
	PlayMethod         string
	TranscodeReasons   string
	Latitude           float64
	Longitude          float64
	PlaybackPosition   int64
	MediaDuration      int64
	PlaybackPercent    int64
	IsPaused           bool
}

func JoinTranscodeReasons(transcodeReasons []TranscodeReasons) string {
	var reasons []string = make([]string, len(transcodeReasons))
	for i, reason := range transcodeReasons {
		reasons[i] = reason.String()
	}
	return strings.Join(reasons, ", ")
}

func (t TranscodeReasons) String() string {
	return string(transcodeReasons[string(t)])
}

func (pm PlayMethod) equal(s string) bool {
	return string(pm) == s
}

func (pm PlayMethod) String() string {
	return string(pm)
}

func (s *Sessions) isEpisode() bool {
	return s.NowPlayingItem.Type == "Episode"
}

func (s *Sessions) GetTranscodeReason() string {
	// Join the slice of strings into a single string.
	return JoinTranscodeReasons(s.TranscodingInfo.TranscodeReasons)
}

func (s *Sessions) To() *SessionsMetrics {
	sessionsMetrics := &SessionsMetrics{
		Username:           s.UserName,
		Client:             s.Client,
		IsPaused:           s.PlayState.IsPaused,
		RemoteEndPoint:     s.RemoteEndPoint,
		NowPlayingItemName: s.NowPlayingItem.Name,
		NowPlayingItemType: s.NowPlayingItem.Type,
		MediaDuration:      s.NowPlayingItem.RunTimeTicks,
		PlaybackPosition:   s.PlayState.PositionTicks,
		PlaybackPercent:    s.getPercentPlayed(),
		PlayMethod:         s.getPlayMethod().String(),
		TranscodeReasons:   s.GetTranscodeReason(),
	}

	if s.isEpisode() {
		sessionsMetrics.TVShow = s.NowPlayingItem.SeriesName
		sessionsMetrics.Season = s.NowPlayingItem.SeasonName
	}

	return sessionsMetrics

}

func (s *Sessions) getPlayMethod() PlayMethod {
	if s.TranscodingInfo == nil {
		return DirectPlay
	} else {
		return Transcoding
	}
}

func (s *Sessions) getPercentPlayed() int64 {
	if s.NowPlayingItem.RunTimeTicks > 0 {
		return s.PlayState.PositionTicks * 100 / s.NowPlayingItem.RunTimeTicks
	}
	return 0
}

func (s *Sessions) HasPlayMethod() bool {
	return s.PlayState.PlayMethod != ""
}
