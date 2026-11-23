// BriefDock local additions: fixed authenticated video-shelf and watch-progress endpoints.
package sdk

import (
	"context"
	"github.com/beilunyang/openhongguo/internal/protocol"
	"github.com/beilunyang/openhongguo/internal/protocol/sign"
)

func (c *Client) BriefDockVideoShelf(ctx context.Context) ([]byte, error) {
	p, err := c.ensureSearchGuest(ctx)
	if err != nil {
		return nil, err
	}
	const path = "/reading/bookapi/bookshelf/video/info/v1/"
	pairs := protocol.DeviceQueryPairs(c.identity(p))
	ticket := c.clock().UnixMilli()
	query, headers := sign.SignPOST(path, pairs, ticket, p.DeviceID, protocol.VersionName, nil, nil)
	raw, _, err := c.doRaw(ctx, "GET", c.apiOrigin+path+"?"+query, anonymousHeaders(c, p, ticket, headers), nil)
	return raw, err
}

// BriefDockReadProgress reads official video watch records without changing them.
func (c *Client) BriefDockReadProgress(ctx context.Context) ([]byte, error) {
	p, err := c.ensureSearchGuest(ctx)
	if err != nil {
		return nil, err
	}
	const path = "/reading/bookapi/read_progress/list/v1/"
	pairs := append(protocol.DeviceQueryPairs(c.identity(p)), [2]string{"book_types", "2"}, [2]string{"limit", "500"})
	ticket := c.clock().UnixMilli()
	query, headers := sign.SignPOST(path, pairs, ticket, p.DeviceID, protocol.VersionName, nil, nil)
	raw, _, err := c.doRaw(ctx, "GET", c.apiOrigin+path+"?"+query, anonymousHeaders(c, p, ticket, headers), nil)
	return raw, err
}

// BriefDockUploadProgress changes only an explicitly supplied watch checkpoint.
func (c *Client) BriefDockUploadProgress(ctx context.Context, body []byte) ([]byte, error) {
	return c.briefDockAccountPOST(ctx, "/reading/bookapi/read_progress/upload/v1/", body)
}

func (c *Client) BriefDockUpdateVideoShelf(ctx context.Context, body []byte) ([]byte, error) {
	return c.briefDockAccountPOST(ctx, "/reading/bookapi/bookshelf/video/update/v1/", body)
}

func (c *Client) briefDockAccountPOST(ctx context.Context, path string, body []byte) ([]byte, error) {
	p, err := c.ensureSearchGuest(ctx)
	if err != nil {
		return nil, err
	}
	ticket := c.clock().UnixMilli()
	query, headers := sign.SignPOST(path, protocol.DeviceQueryPairs(c.identity(p)), ticket, p.DeviceID, protocol.VersionName, body, nil)
	headers = anonymousHeaders(c, p, ticket, headers)
	headers["Content-Type"] = "application/json; charset=utf-8"
	raw, _, err := c.doRaw(ctx, "POST", c.apiOrigin+path+"?"+query, headers, body)
	return raw, err
}
