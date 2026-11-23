package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/beilunyang/openhongguo/internal/protocol"
	"github.com/beilunyang/openhongguo/internal/protocol/sign"
	"github.com/beilunyang/openhongguo/sdk/model"
)

// BriefDockPlayback uses one current player RPC with the existing device only.
func (c *Client) BriefDockPlayback(ctx context.Context, episodeID, quality string) (*model.Playback, error) {
	if !validID(episodeID) {
		return nil, fmt.Errorf("invalid episode ID")
	}
	p, err := c.ensureSearchGuest(ctx)
	if err != nil {
		return nil, err
	}
	pairs := append(protocol.DeviceQueryPairs(c.identity(p)),
		[2]string{"video_id", episodeID}, [2]string{"video_id_type", "0"}, [2]string{"device_level", "1"},
		[2]string{"need_all_video_definition", "1"}, [2]string{"need_mp4_align", "0"},
		[2]string{"use_os_player", "0"}, [2]string{"video_platform", "1024"}, [2]string{"source", "0"}, [2]string{"from_video_id", ""})
	query, headers, err := sign.SignPlayerGET(pairs, c.clock().UnixMilli())
	if err != nil {
		return nil, err
	}
	d := c.deviceOf(p)
	headers["User-Agent"] = fmt.Sprintf("%s/%s (Linux; U; Android %s; zh_CN; %s; Build/%s;tt-ok/3.12.13.20)", protocol.Package, protocol.VersionCode, d.OSVersion, d.Model, d.Build)
	if p.TTReq != "" {
		headers["Cookie"] = guestCookie(p)
	}
	raw, _, err := c.doRaw(ctx, "GET", "https://api5-normal-sinfonlineb.fqnovel.com/reading/bookapi/video_tab/multi_video_model/v1/?"+query, headers, nil)
	if err != nil {
		return nil, err
	}
	video, err := playerResponseModel(raw, episodeID)
	if err != nil {
		return nil, err
	}
	return parsePlayerModel(video, episodeID, quality)
}

func playerResponseModel(raw []byte, episodeID string) (json.RawMessage, error) {
	var env struct {
		Data map[string]json.RawMessage `json:"data"`
	}
	if json.Unmarshal(raw, &env) != nil {
		return nil, fmt.Errorf("invalid playback response")
	}
	rows := env.Data
	if nested, ok := rows["video_model_data"]; ok {
		if json.Unmarshal(nested, &rows) != nil {
			return nil, fmt.Errorf("invalid playback response")
		}
	}
	var row struct {
		VideoModel json.RawMessage `json:"video_model"`
	}
	if json.Unmarshal(rows[episodeID], &row) != nil || len(row.VideoModel) == 0 {
		return nil, fmt.Errorf("playback_episode_mismatch")
	}
	if len(row.VideoModel) > 0 && row.VideoModel[0] == '"' {
		var text string
		if json.Unmarshal(row.VideoModel, &text) != nil {
			return nil, fmt.Errorf("invalid playback model")
		}
		row.VideoModel = json.RawMessage(text)
	}
	return row.VideoModel, nil
}

func parsePlayerModel(video json.RawMessage, episodeID, quality string) (*model.Playback, error) {
	wrapped, _ := json.Marshal(map[string]any{"code": 0, "data": map[string]any{"video_model": video}})
	playback, err := protocol.ParsePlaybackForQuality(wrapped, episodeID, quality)
	if err != nil {
		return nil, err
	}
	if playback.Protected {
		return playback, ErrDRMProtected
	}
	return playback, nil
}
