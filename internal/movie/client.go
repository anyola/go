package movie

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var ErrNotFound = errors.New("movie not found")

type Client struct {
	baseURL string
	http    *http.Client
	timeout time.Duration
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		http:    &http.Client{},
		timeout: 10 * time.Second,
	}
}

func (c *Client) Get(ctx context.Context, id int) (*Movie, error) {
	reqCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet,
		fmt.Sprintf("%s/%d/info.0.json", c.baseURL, id), nil)
	if err != nil {
		return nil, fmt.Errorf("create movie request: %w", err)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("movie stand request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("movie stand returned HTTP %d", resp.StatusCode)
	}
	var result Movie
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode movie %d: %w", id, err)
	}
	return &result, nil
}
