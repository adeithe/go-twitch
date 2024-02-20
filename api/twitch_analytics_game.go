package api

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type AnalyticsGameResource struct {
	client *Client
}

func NewAnalyticsGameResource(client *Client) *AnalyticsGameResource {
	return &AnalyticsGameResource{client}
}

type AnalyticsGameListRequest struct {
	client *Client
	opts   []RequestOption
}

type AnalyticsGameListResponse struct {
	Header     http.Header
	Data       []GameAnalytics
	Pagination Pagination
}

type GameAnalytics struct {
	GameID   string             `json:"game_id"`
	URL      string             `json:"URL"`
	DateRate AnalyticsDateRange `json:"date_range"`
}

func (r *AnalyticsGameResource) List() *AnalyticsGameListRequest {
	return &AnalyticsGameListRequest{client: r.client}
}

// GameID If specified, the response contains a report for the specified game.
// If not specified, the response includes a report for each game that the authenticated user has played.
func (r *AnalyticsGameListRequest) GameID(id string) *AnalyticsGameListRequest {
	r.opts = append(r.opts, SetQueryParameter("game_id", id))
	return r
}

// StartedAt The start of the date range for the report.
func (r *AnalyticsGameListRequest) StartedAt(t time.Time) *AnalyticsGameListRequest {
	r.opts = append(r.opts, SetQueryParameter("started_at", t.Format(time.RFC3339)))
	return r
}

// EndedAt The end of the date range for the report.
func (r *AnalyticsGameListRequest) EndedAt(t time.Time) *AnalyticsGameListRequest {
	r.opts = append(r.opts, SetQueryParameter("ended_at", t.Format(time.RFC3339)))
	return r
}

// First The number of records to return. Maximum: 100. Default: 20.
func (r *AnalyticsGameListRequest) First(f int) *AnalyticsGameListRequest {
	r.opts = append(r.opts, SetQueryParameter("first", fmt.Sprint(f)))
	return r
}

// After A cursor for forward pagination: the first set of results to return. Provide this value in the after query parameter.
func (r *AnalyticsGameListRequest) After(a string) *AnalyticsGameListRequest {
	r.opts = append(r.opts, SetQueryParameter("after", a))
	return r
}

// Do executes the request.
func (r *AnalyticsGameListRequest) Do(ctx context.Context, opts ...RequestOption) (*AnalyticsGameListResponse, error) {
	res, err := r.client.doRequest(ctx, http.MethodGet, "/analytics/games", nil, append(r.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[GameAnalytics](res)
	if err != nil {
		return nil, err
	}

	return &AnalyticsGameListResponse{
		Header:     res.Header,
		Data:       data.Data,
		Pagination: data.Pagination,
	}, nil
}
