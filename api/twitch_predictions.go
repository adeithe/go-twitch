package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// PredictionsResource represents the Twitch Predictions API.
type PredictionsResource struct {
	client *Client
}

// NewPredictionsResource creates a new PredictionsResource.
func NewPredictionsResource(client *Client) *PredictionsResource {
	return &PredictionsResource{client}
}

// PredictionsListCall represents a GET call to a Twitch Predictions API endpoint.
type PredictionsListCall struct {
	resource *PredictionsResource
	opts     []RequestOption
}

// PredictionsListResponse represents the response from a GET request to /helix/predictions.
type PredictionsListResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Prediction data returned by the Twitch API.
	Data []Prediction
	// Pagination is the Pagination data returned by the Twitch API.
	Pagination Pagination
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// List creates a new GET request to /helix/predictions.
//
// Gets a list of Channel Points Predictions that the broadcaster created.
//
// # Authorization
//
// Requires a user access token that includes the channel:read:predictions or channel:manage:predictions scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#get-predictions
func (r *PredictionsResource) List(broadcasterID string) *PredictionsListCall {
	c := &PredictionsListCall{resource: r}
	return c.
		BroadcasterID(broadcasterID)
}

// BroadcasterID sets the BroadcasterID query parameter.
func (api *PredictionsListCall) BroadcasterID(broadcasterID string) *PredictionsListCall {
	api.opts = append(api.opts, SetQueryParameter("broadcaster_id", broadcasterID))
	return api
}

// ID adds to the ID query parameter.
func (api *PredictionsListCall) ID(ids ...string) *PredictionsListCall {
	for _, id := range ids {
		api.opts = append(api.opts, AddQueryParameter("id", id))
	}
	return api
}

// After sets the After query parameter.
func (api *PredictionsListCall) After(after string) *PredictionsListCall {
	api.opts = append(api.opts, SetQueryParameter("after", after))
	return api
}

// First sets the First query parameter.
func (api *PredictionsListCall) First(first int) *PredictionsListCall {
	api.opts = append(api.opts, SetQueryParameter("first", first))
	return api
}

// Do executes the request.
func (api *PredictionsListCall) Do(ctx context.Context, opts ...RequestOption) (*PredictionsListResponse, error) {
	res, err := api.resource.client.DoRequest(ctx, "GET", "/helix/predictions", nil, append(api.opts, opts...)...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Prediction](res)
	if err != nil {
		return nil, err
	}

	return &PredictionsListResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Pagination: data.Pagination,
		Request:    res.Request,
	}, nil
}

// PredictionsInsertCall represents a POST call to a Twitch Predictions API endpoint.
type PredictionsInsertCall struct {
	resource *PredictionsResource
	body     map[string]any
}

// PredictionsInsertResponse represents the response from a POST request to /helix/predictions.
type PredictionsInsertResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Prediction data returned by the Twitch API.
	Data []Prediction
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Insert creates a new POST request to /helix/predictions.
//
// Creates a Channel Points Prediction.
//
// With a Channel Points Prediction, the broadcaster poses a question and viewers try to predict the outcome. The prediction runs as soon as it's created. The broadcaster may run only one prediction at a time.
//
// # Authorization
//
// Requires a user access token that includes the channel:manage:predictions scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#create-prediction
func (r *PredictionsResource) Insert(broadcasterID string, title string, predictionWindow int) *PredictionsInsertCall {
	c := &PredictionsInsertCall{resource: r, body: make(map[string]any)}
	return c.
		BroadcasterID(broadcasterID).
		Title(title).
		PredictionWindow(predictionWindow)
}

// BroadcasterID sets the BroadcasterID body parameter.
func (api *PredictionsInsertCall) BroadcasterID(broadcasterID string) *PredictionsInsertCall {
	api.body["broadcaster_id"] = broadcasterID
	return api
}

// Title sets the Title body parameter.
func (api *PredictionsInsertCall) Title(title string) *PredictionsInsertCall {
	api.body["title"] = title
	return api
}

// Outcome sets the Outcome body parameter.
func (api *PredictionsInsertCall) Outcome(outcomes ...OutboundChoice) *PredictionsInsertCall {
	api.body["outcomes"] = outcomes
	return api
}

// PredictionWindow sets the PredictionWindow body parameter.
func (api *PredictionsInsertCall) PredictionWindow(predictionWindow int) *PredictionsInsertCall {
	api.body["prediction_window"] = predictionWindow
	return api
}

// Do executes the request.
func (api *PredictionsInsertCall) Do(ctx context.Context, opts ...RequestOption) (*PredictionsInsertResponse, error) {
	bs, err := json.Marshal(api.body)
	if err != nil {
		return nil, err
	}

	res, err := api.resource.client.DoRequest(ctx, "POST", "/helix/predictions", bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Prediction](res)
	if err != nil {
		return nil, err
	}

	return &PredictionsInsertResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}

// PredictionsModifyCall represents a PATCH call to a Twitch Predictions API endpoint.
type PredictionsModifyCall struct {
	resource *PredictionsResource
	body     map[string]any
}

// PredictionsModifyResponse represents the response from a PATCH request to /helix/predictions.
type PredictionsModifyResponse struct {
	// Status is the HTTP status text returned by the Twitch API. For example, "200 OK".
	Status string
	// StatusCode is the HTTP status code returned by the Twitch API. For example, 200.
	StatusCode int
	// Header contains the HTTP headers from the Twitch API response.
	Header http.Header
	// Data is the Prediction data returned by the Twitch API.
	Data []Prediction
	// Request is the HTTP request that was sent to the Twitch API.
	Request *http.Request
}

// Modify creates a new PATCH request to /helix/predictions.
//
// Locks, resolves, or cancels a Channel Points Prediction.
//
// # Authorization
//
// Requires a user access token that includes the channel:manage:predictions scope.
//
// Check the [Official Twitch Documentation] for more information.
//
// [Official Twitch Documentation]: https://dev.twitch.tv/docs/api/reference/#end-prediction
func (r *PredictionsResource) Modify(broadcasterID string, id string, status string) *PredictionsModifyCall {
	c := &PredictionsModifyCall{resource: r, body: make(map[string]any)}
	return c.
		BroadcasterID(broadcasterID).
		ID(id).
		Status(status)
}

// BroadcasterID sets the BroadcasterID body parameter.
func (api *PredictionsModifyCall) BroadcasterID(broadcasterID string) *PredictionsModifyCall {
	api.body["broadcaster_id"] = broadcasterID
	return api
}

// ID sets the ID body parameter.
func (api *PredictionsModifyCall) ID(id string) *PredictionsModifyCall {
	api.body["id"] = id
	return api
}

// Status sets the Status body parameter.
func (api *PredictionsModifyCall) Status(status string) *PredictionsModifyCall {
	api.body["status"] = status
	return api
}

// WinningOutcomeID sets the WinningOutcomeID body parameter.
func (api *PredictionsModifyCall) WinningOutcomeID(winningOutcomeID string) *PredictionsModifyCall {
	api.body["winning_outcome_id"] = winningOutcomeID
	return api
}

// Do executes the request.
func (api *PredictionsModifyCall) Do(ctx context.Context, opts ...RequestOption) (*PredictionsModifyResponse, error) {
	bs, err := json.Marshal(api.body)
	if err != nil {
		return nil, err
	}

	res, err := api.resource.client.DoRequest(ctx, "PATCH", "/helix/predictions", bytes.NewReader(bs), opts...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	data, err := decodeResponse[Prediction](res)
	if err != nil {
		return nil, err
	}

	return &PredictionsModifyResponse{
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Data:       data.Data,
		Request:    res.Request,
	}, nil
}
