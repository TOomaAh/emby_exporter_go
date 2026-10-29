package entity

import "testing"

const ticksPerSecond = 10_000_000

func TestGetDuration(t *testing.T) {
	s := &Sessions{}
	tests := []struct {
		ticks int64
		want  string
	}{
		{0, ""},
		{5 * ticksPerSecond, "00h00m05s"},
		{(3600 + 2*60 + 3) * ticksPerSecond, "01h02m03s"},
		{(12*3600 + 34*60 + 56) * ticksPerSecond, "12h34m56s"},
	}
	for _, tt := range tests {
		if got := s.GetDuration(tt.ticks); got != tt.want {
			t.Errorf("GetDuration(%d) = %q, want %q", tt.ticks, got, tt.want)
		}
	}
}

func TestGetPercentPlayed(t *testing.T) {
	tests := []struct {
		name    string
		session Sessions
		want    int64
	}{
		{"no item", Sessions{PlayState: PlayState{PositionTicks: 100}}, 0},
		{"unknown runtime", Sessions{NowPlayingItem: &NowPlayingItem{}, PlayState: PlayState{PositionTicks: 100}}, 0},
		{"start", Sessions{NowPlayingItem: &NowPlayingItem{RunTimeTicks: 1000}}, 0},
		{"half", Sessions{NowPlayingItem: &NowPlayingItem{RunTimeTicks: 1000}, PlayState: PlayState{PositionTicks: 500}}, 50},
		{"end", Sessions{NowPlayingItem: &NowPlayingItem{RunTimeTicks: 1000}, PlayState: PlayState{PositionTicks: 1000}}, 100},
	}
	for _, tt := range tests {
		if got := tt.session.GetPercentPlayed(); got != tt.want {
			t.Errorf("%s: GetPercentPlayed() = %d, want %d", tt.name, got, tt.want)
		}
	}
}

func TestBitsToHumanReadable(t *testing.T) {
	tests := []struct {
		bits uint64
		want string
	}{
		{0, "0 b"},
		{512, "512 b"},
		{1024, "1 Kb"},
		{1536, "1.50 Kb"},
		{10 * 1024 * 1024, "10 Mb"},
		{3 * 1024 * 1024 * 1024, "3 Gb"},
	}
	for _, tt := range tests {
		if got := bitsToHumanReadable(tt.bits); got != tt.want {
			t.Errorf("bitsToHumanReadable(%d) = %q, want %q", tt.bits, got, tt.want)
		}
	}
}

func TestGetBitrateValue(t *testing.T) {
	tests := []struct {
		name    string
		session Sessions
		want    uint64
	}{
		{"no item", Sessions{}, 0},
		{"direct play", Sessions{NowPlayingItem: &NowPlayingItem{Bitrate: 8000}}, 8000},
		{"transcoding", Sessions{NowPlayingItem: &NowPlayingItem{Bitrate: 8000}, TranscodingInfo: &TranscodingInfo{Bitrate: 3000}}, 3000},
	}
	for _, tt := range tests {
		if got := tt.session.GetBitrateValue(); got != tt.want {
			t.Errorf("%s: GetBitrateValue() = %d, want %d", tt.name, got, tt.want)
		}
	}
}

func TestGetPlayMethod(t *testing.T) {
	if got := (&Sessions{}).GetPlayMethod(); got != "Direct Play" {
		t.Errorf("without transcoding: %q, want Direct Play", got)
	}
	if got := (&Sessions{TranscodingInfo: &TranscodingInfo{}}).GetPlayMethod(); got != "Transcoding" {
		t.Errorf("with transcoding: %q, want Transcoding", got)
	}
}

func TestGetTranscodeReason(t *testing.T) {
	if got := (&Sessions{}).GetTranscodeReason(); got != "" {
		t.Errorf("without transcoding: %q, want empty", got)
	}

	s := &Sessions{TranscodingInfo: &TranscodingInfo{
		TranscodeReasons: []TranscodeReasons{"VideoCodecNotSupported", "AudioCodecNotSupported"},
	}}
	want := "Video Codec Not Supported, Audio Codec Not Supported"
	if got := s.GetTranscodeReason(); got != want {
		t.Errorf("GetTranscodeReason() = %q, want %q", got, want)
	}
}

func TestIsEpisode(t *testing.T) {
	tests := []struct {
		name string
		item *NowPlayingItem
		want bool
	}{
		{"no item", nil, false},
		{"movie", &NowPlayingItem{Type: "Movie"}, false},
		{"episode", &NowPlayingItem{Type: "Episode"}, true},
	}
	for _, tt := range tests {
		s := &Sessions{NowPlayingItem: tt.item}
		if got := s.IsEpisode(); got != tt.want {
			t.Errorf("%s: IsEpisode() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
