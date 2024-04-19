package faketime

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ClientInterface 時刻偽装モード API のクライアントインターフェース。
type ClientInterface interface {
	// Get 固定されたサーバー時刻を取得する。
	Get(ctx context.Context) (time.Time, error)
	// Set サーバー時刻を t に変更する。
	Set(ctx context.Context, t time.Time) error
}

var _ ClientInterface = (*Client)(nil)

// Client 時刻偽装モード API のクライアントを表す。
type Client struct {
	client   *http.Client
	endpoint string
}

// NewClient Client を生成して返す。
func NewClient(client *http.Client, endpoint string) *Client {
	return &Client{
		client:   client,
		endpoint: endpoint,
	}
}

// Get 固定されたサーバー時刻を取得する。
func (s *Client) Get(ctx context.Context) (time.Time, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.endpoint+getPath, nil)
	if err != nil {
		return time.Time{}, fmt.Errorf("new request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return time.Time{}, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return time.Time{}, fmt.Errorf("http error, %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return time.Time{}, fmt.Errorf("read body: %w", err)
	}

	t, err := time.Parse(time.RFC3339, string(body))
	if err != nil {
		return time.Time{}, fmt.Errorf("parse response: %w", err)
	}

	return t, nil
}

// Set サーバー時刻を t に変更する。
func (s *Client) Set(ctx context.Context, t time.Time) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.endpoint+setPath, nil)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	query := req.URL.Query()
	query.Add("t", t.Format(time.RFC3339))
	req.URL.RawQuery = query.Encode()

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http error, %s", resp.Status)
	}

	return nil
}
