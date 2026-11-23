package sign

import (
	"fmt"
	"strconv"
)

// SignPlayerGET binds the current player's query to its full Argus/Ladon headers.
// Wire fields: hongguoTV 77a99f5, ContentRepository.playbackRequest / Signer.signPlayback.
// Account signing stays separate; do not mix its legacy headers into playback.
func SignPlayerGET(pairs [][2]string, nowMillis int64) (string, map[string]string, error) {
	body := [16]byte{}
	for offset := int64(0); offset < 128; offset++ {
		ticket := nowMillis + offset
		ts := uint32(ticket / 1000)
		values := append(append([][2]string{}, pairs...), [2]string{"_rticket", strconv.FormatInt(ticket, 10)})
		query := encodeQuery(values)
		hash := sm3Sum([]byte(query))
		if branchOf(hash[:], body[:], le32(ts)) != 0 {
			continue
		}
		aid := playbackAID(query)
		return query, map[string]string{
			"Accept":      "application/json; charset=utf-8,application/x-protobuf",
			"sdk-version": "2", "x-xs-from-web": "0",
			"x-ss-req-ticket": strconv.FormatInt(ticket, 10), "x-khronos": strconv.FormatUint(uint64(ts), 10),
			"x-argus": argusHeader(query, body, ts, randUint32(), aid, queryValue(query, "device_id"), ""),
			"x-ladon": ladonHeader(ts, randUint32(), aid),
		}, nil
	}
	return "", nil, fmt.Errorf("playback signing unavailable")
}
