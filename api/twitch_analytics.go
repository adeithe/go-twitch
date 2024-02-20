package api

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type AnalyticsResource struct {
	client *Client

	Extensions *AnalyticsExtensionResource
	Games      *AnalyticsGameResource
}

type AnalyticsDateRange struct {
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
}

func NewAnalyticsResource(client *Client) *AnalyticsResource {
	r := &AnalyticsResource{client: client}
	r.Extensions = NewAnalyticsExtensionResource(client)
	r.Games = NewAnalyticsGameResource(client)
	return r
}

type AnalyticsExtensionResource struct {
	client *Client
}

func NewAnalyticsExtensionResource(client *Client) *AnalyticsExtensionResource {
	return &AnalyticsExtensionResource{client}
}

type AnalyticsExtensionListCall struct {
	client *Client
	opts   []RequestOption
}

type AnalyticsExtensionListResponse struct {
	Header     http.Header
	Data       []ExtensionAnalytics
	Pagination Pagination
}

type ExtensionAnalytics struct {
	ExtensionID string             `json:"extension_id"`
	URL         string             `json:"URL"`
	Type        string             `json:"type"`
	DateRate    AnalyticsDateRange `json:"date_range"`
}

func (r *AnalyticsExtensionResource) List() *AnalyticsExtensionListCall {
	return &AnalyticsExtensionListCall{client: r.client}
}

// ExtensionID If specified, the response contains a report for the specified extension.
// If not specified, the response includes a report for each extension that the authenticated user owns.
func (r *AnalyticsExtensionListCall) ExtensionID(id string) *AnalyticsExtensionListCall {
	r.opts = append(r.opts, SetQueryParameter("extension_id", id))
	return r
}

// Type The type of analytics report to get. Possible values are:
//   - overview_v2
func (r *AnalyticsExtensionListCall) Type(t string) *AnalyticsExtensionListCall {
	r.opts = append(r.opts, SetQueryParameter("type", t))
	return r
}

// StartedAt The start of the date range for the report.
func (r *AnalyticsExtensionListCall) StartedAt(t time.Time) *AnalyticsExtensionListCall {
	r.opts = append(r.opts, SetQueryParameter("started_at", t.Format(time.RFC3339)))
	return r
}

// EndedAt The end of the date range for the report.
func (r *AnalyticsExtensionListCall) EndedAt(t time.Time) *AnalyticsExtensionListCall {
	r.opts = append(r.opts, SetQueryParameter("ended_at", t.Format(time.RFC3339)))
	return r
}

// First The number of records to return. Maximum: 100. Default: 20.
func (r *AnalyticsExtensionListCall) First(n int) *AnalyticsExtensionListCall {
	r.opts = append(r.opts, SetQueryParameter("first", fmt.Sprint(n)))
	return r
}

// After A cursor for forward pagination: the first set of results to return. Provide this value in the after query parameter.
func (r *AnalyticsExtensionListCall) After(cursor string) *AnalyticsExtensionListCall {
	r.opts = append(r.opts, SetQueryParameter("after", cursor))
	return r
}

// Do executes the request.
func (r *AnalyticsExtensionListCall) Do(ctx context.Context, opts ...RequestOption) (*AnalyticsExtensionListResponse, error) {
	res, err := r.client.doRequest(ctx, http.MethodGet, "/analytics/extensions", nil, append(r.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := decodeResponse[ExtensionAnalytics](res)
	if err != nil {
		return nil, err
	}

	return &AnalyticsExtensionListResponse{
		Header:     res.Header,
		Data:       data.Data,
		Pagination: data.Pagination,
	}, nil
}

type AnalyticsGameResource struct {
	client *Client
}

func NewAnalyticsGameResource(client *Client) *AnalyticsGameResource {
	return &AnalyticsGameResource{client}
}

type AnalyticsGameListCall struct {
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

func (r *AnalyticsGameResource) List() *AnalyticsGameListCall {
	return &AnalyticsGameListCall{client: r.client}
}

// GameID If specified, the response contains a report for the specified game.
// If not specified, the response includes a report for each game that the authenticated user has played.
func (r *AnalyticsGameListCall) GameID(id string) *AnalyticsGameListCall {
	r.opts = append(r.opts, SetQueryParameter("game_id", id))
	return r
}

// StartedAt The start of the date range for the report.
func (r *AnalyticsGameListCall) StartedAt(t time.Time) *AnalyticsGameListCall {
	r.opts = append(r.opts, SetQueryParameter("started_at", t.Format(time.RFC3339)))
	return r
}

// EndedAt The end of the date range for the report.
func (r *AnalyticsGameListCall) EndedAt(t time.Time) *AnalyticsGameListCall {
	r.opts = append(r.opts, SetQueryParameter("ended_at", t.Format(time.RFC3339)))
	return r
}

// First The number of records to return. Maximum: 100. Default: 20.
func (r *AnalyticsGameListCall) First(f int) *AnalyticsGameListCall {
	r.opts = append(r.opts, SetQueryParameter("first", fmt.Sprint(f)))
	return r
}

// After A cursor for forward pagination: the first set of results to return. Provide this value in the after query parameter.
func (r *AnalyticsGameListCall) After(a string) *AnalyticsGameListCall {
	r.opts = append(r.opts, SetQueryParameter("after", a))
	return r
}

// Do executes the request.
func (r *AnalyticsGameListCall) Do(ctx context.Context, opts ...RequestOption) (*AnalyticsGameListResponse, error) {
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
