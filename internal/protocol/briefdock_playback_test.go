package protocol

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestQualityKeepsSupportedCodecAndItsOwnMaterial(t *testing.T) {
	items := []any{
		map[string]any{"main_url": "https://fixture.qznovelvod.com/unsupported", "codec_type": "bytevc2", "definition": "720p", "bitrate": 1, "spade_a": "unsupported-fixture"},
		map[string]any{"main_url": "https://fixture.qznovelvod.com/low", "codec_type": "bytevc1", "definition": "360p", "bitrate": 2, "spade_a": "low-fixture"},
		map[string]any{"main_url": "https://fixture.qznovelvod.com/hd", "video_meta": map[string]any{"codec_type": "bytevc1", "definition": "normal", "vwidth": 1280, "vheight": 720}, "bitrate": 3, "spade_a": "hd-fixture"},
	}
	raw, _ := json.Marshal(map[string]any{"code": 0, "data": map[string]any{"video_model": map[string]any{"video_list": items}}})
	value, err := ParsePlaybackForQuality(raw, "12345", "720p")
	if err != nil || len(value.Streams) != 1 || value.Streams[0].Codec != "bytevc1" || value.KeyMaterial != "hd-fixture" {
		t.Fatal("selected stream and material diverged")
	}
	public, _ := json.Marshal(value)
	if strings.Contains(string(public), "fixture\"") || strings.Contains(string(public), "spade_a") {
		t.Fatal("material leaked into public JSON")
	}
	value, err = ParsePlaybackForQuality(raw, "12345", "1080p")
	if err != nil || value.Streams[0].Height != 720 || value.Streams[0].Definition != "720p" || value.KeyMaterial != "hd-fixture" {
		t.Fatal("missing preferred quality did not select the highest compatible lower stream")
	}
}

func TestHigherCompatibleRenditionKeepsActualResolution(t *testing.T) {
	items := []any{
		map[string]any{"main_url": "https://fixture.qznovelvod.com/new", "codec_type": "bytevc2", "vwidth": 1280, "vheight": 720, "spade_a": "unsupported-fixture"},
		map[string]any{"main_url": "https://fixture.qznovelvod.com/hd", "codec_type": "bytevc1", "vwidth": 1920, "vheight": 1080, "spade_a": "hd-fixture"},
	}
	raw, _ := json.Marshal(map[string]any{"data": map[string]any{"video_model": map[string]any{"video_list": items}}})
	for _, wanted := range []string{"360p", "480p", "720p"} {
		value, err := ParsePlaybackForQuality(raw, "12345", wanted)
		if err != nil || len(value.Streams) != 1 || value.Streams[0].Height != 1080 || value.Streams[0].Definition != "1080p" || value.KeyMaterial != "hd-fixture" {
			t.Fatal("compatible resolution or material was misreported")
		}
	}
}

func TestQualityDoesNotFallBackToUnsupportedStream(t *testing.T) {
	items := []any{
		map[string]any{"main_url": "https://fixture.qznovelvod.com/hd", "codec_type": "bytevc2", "vwidth": 1920, "vheight": 1080},
		map[string]any{"main_url": "https://fixture.qznovelvod.com/mid", "codec_type": "h264", "vwidth": 1280, "vheight": 720},
		map[string]any{"main_url": "https://fixture.qznovelvod.com/low", "codec_type": "h264", "vwidth": 640, "vheight": 360},
	}
	raw, _ := json.Marshal(map[string]any{"data": map[string]any{"video_model": map[string]any{"video_list": items}}})
	value, err := ParsePlaybackForQuality(raw, "12345", "1080p")
	if err != nil || value.Streams[0].Height != 720 || len(value.AvailableQualities) != 2 {
		t.Fatal("fallback used an unsupported codec or wrong quality")
	}
	raw, _ = json.Marshal(map[string]any{"data": map[string]any{"video_model": map[string]any{"video_list": items[:1]}}})
	if _, err = ParsePlaybackForQuality(raw, "12345", "1080p"); err == nil {
		t.Fatal("unsupported stream accepted")
	}
}

func TestQualityPreferenceSurvivesSeriesSwitch(t *testing.T) {
	for _, height := range []int{1080, 720, 1080} {
		items := []any{map[string]any{"main_url": "https://fixture.qznovelvod.com/video", "codec_type": "h264", "vwidth": height, "vheight": height * 16 / 9}}
		raw, _ := json.Marshal(map[string]any{"data": map[string]any{"video_model": map[string]any{"video_list": items}}})
		value, err := ParsePlaybackForQuality(raw, "12345", "1080p")
		if err != nil || value.Streams[0].Width != height || value.Streams[0].Definition != map[int]string{720: "720p", 1080: "1080p"}[height] {
			t.Fatal("series switch changed the saved preference or actual resolution")
		}
	}
}
