package api

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type AnalyticsExtensionResource struct {
	client *Client
}

func NewAnalyticsExtensionResource(client *Client) *AnalyticsExtensionResource {
	return &AnalyticsExtensionResource{client}
}

type AnalyticsExtensionListRequest struct {
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

func (r *AnalyticsExtensionResource) List() *AnalyticsExtensionListRequest {
	return &AnalyticsExtensionListRequest{client: r.client}
}

// ExtensionID If specified, the response contains a report for the specified extension.
// If not specified, the response includes a report for each extension that the authenticated user owns.
func (r *AnalyticsExtensionListRequest) ExtensionID(id string) *AnalyticsExtensionListRequest {
	r.opts = append(r.opts, SetQueryParameter("extension_id", id))
	return r
}

// Type The type of analytics report to get. Possible values are:
//   - overview_v2
func (r *AnalyticsExtensionListRequest) Type(t string) *AnalyticsExtensionListRequest {
	r.opts = append(r.opts, SetQueryParameter("type", t))
	return r
}

// StartedAt The start of the date range for the report.
func (r *AnalyticsExtensionListRequest) StartedAt(t time.Time) *AnalyticsExtensionListRequest {
	r.opts = append(r.opts, SetQueryParameter("started_at", t.Format(time.RFC3339)))
	return r
}

// EndedAt The end of the date range for the report.
func (r *AnalyticsExtensionListRequest) EndedAt(t time.Time) *AnalyticsExtensionListRequest {
	r.opts = append(r.opts, SetQueryParameter("ended_at", t.Format(time.RFC3339)))
	return r
}

// First The number of records to return. Maximum: 100. Default: 20.
func (r *AnalyticsExtensionListRequest) First(n int) *AnalyticsExtensionListRequest {
	r.opts = append(r.opts, SetQueryParameter("first", fmt.Sprint(n)))
	return r
}

// After A cursor for forward pagination: the first set of results to return. Provide this value in the after query parameter.
func (r *AnalyticsExtensionListRequest) After(cursor string) *AnalyticsExtensionListRequest {
	r.opts = append(r.opts, SetQueryParameter("after", cursor))
	return r
}

// Do executes the request.
func (r *AnalyticsExtensionListRequest) Do(ctx context.Context, opts ...RequestOption) (*AnalyticsExtensionListResponse, error) {
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
