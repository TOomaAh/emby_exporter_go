package entity

import (
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

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
	SeriesName   string `json:"SeriesName"`
	SeasonName   string `json:"SeasonName"`
	MediaType    string `json:"MediaType"`
	Type         string `json:"Type"`
	RunTimeTicks int64  `json:"RunTimeTicks"`
	Bitrate      uint64 `json:"Bitrate"`
	IndexNumber  int    `json:"IndexNumber"`
}

type TranscodingInfo struct {
	TranscodeReasons              []TranscodeReasons `json:"TranscodeReasons"`
	CurrentCPUUsage               float64            `json:"CurrentCpuUsage"`
	TranscodingPositionTicks      int64              `json:"TranscodingPositionTicks"`
	TranscodingStartPositionTicks int64              `json:"TranscodingStartPositionTicks"`
	Bitrate                       uint64             `json:"Bitrate"`
	CurrentThrottle               int64              `json:"CurrentThrottle"`
	AudioCodec                    string             `json:"AudioCodec"`
	VideoCodec                    string             `json:"VideoCodec"`
	VideoDecoderIsHardware        bool               `json:"VideoDecoderIsHardware"`
	VideoEncoderIsHardware        bool               `json:"VideoEncoderIsHardware"`
	IsVideoDirect                 bool               `json:"IsVideoDirect"`
	IsAudioDirect                 bool               `json:"IsAudioDirect"`
}

type Sessions struct {
	NowPlayingItem  *NowPlayingItem  `json:"NowPlayingItem,omitempty"`
	TranscodingInfo *TranscodingInfo `json:"TranscodingInfo,omitempty"`
	PlayState       PlayState        `json:"PlayState,omitempty"`
	RemoteEndPoint  string           `json:"RemoteEndPoint"`
	UserName        string           `json:"UserName,omitempty"`
	Client          string           `json:"Client"`
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

func (pm PlayMethod) String() string {
	return string(pm)
}

func (s *Sessions) IsEpisode() bool {
	return s.NowPlayingItem.Type == "Episode"
}

func formatDuration(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	return formatTimeComponent(hours) + "h" + formatTimeComponent(minutes) + "m" + formatTimeComponent(seconds) + "s"
}

func formatTimeComponent(component int) string {
	if component < 10 {
		return "0" + strconv.Itoa(component)
	}
	return strconv.Itoa(component)
}

func (s *Sessions) GetDuration(tick int64) string {
	if tick == 0 {
		return ""
	}
	duration := time.Duration(tick * 100)
	// return HH:MM:SS
	return formatDuration(duration)
}

func (s *Sessions) GetTranscodeReason() string {
	// Join the slice of strings into a single string.
	if s.TranscodingInfo == nil {
		return ""
	}
	return JoinTranscodeReasons(s.TranscodingInfo.TranscodeReasons)
}

func (s *Sessions) GetEpisodeNumber() string {
	if s.IsEpisode() {
		if s.NowPlayingItem.IndexNumber > 0 {
			return "Ep. " + strconv.Itoa(s.NowPlayingItem.IndexNumber) + " - "
		}
	}
	return ""
}

func (s *Sessions) GetRuntimeTick() int64 {
	if s.NowPlayingItem == nil || s.NowPlayingItem.RunTimeTicks == 0 {
		return 0
	}
	return s.NowPlayingItem.RunTimeTicks
}

func byteToMb(b uint64) string {
	return humanize.Bytes(b)
}

func (s *Sessions) GetBitrateValue() uint64 {
	if s.TranscodingInfo == nil {
		return s.NowPlayingItem.Bitrate
	}
	return s.TranscodingInfo.Bitrate
}

func (s *Sessions) GetBitrateFormat() string {
	return byteToMb(s.GetBitrateValue())
}

func (s *Sessions) GetPlayMethod() string {
	if s.TranscodingInfo == nil {
		return DirectPlay.String()
	} else {
		return Transcoding.String()
	}
}

func (s *Sessions) GetPercentPlayed() int64 {
	if s.NowPlayingItem.RunTimeTicks > 0 {
		return s.PlayState.PositionTicks * 100 / s.NowPlayingItem.RunTimeTicks
	}
	return 0
}

func (s *Sessions) HasPlayMethod() bool {
	return s.PlayState.PlayMethod != ""
}
