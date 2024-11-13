package urlshortener

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
)

type HTTPClient struct {
	client *resty.Client
}

type shortenResponse struct {
	Shortened string `json:"shortened"`
}

func (c HTTPClient) Shorten(rawURL string, expiration *time.Time) (string, error) {
	request := c.client.R().
		SetQueryParam("url", url.QueryEscape(rawURL))
	if expiration != nil {
		request.SetQueryParam("expiration", url.QueryEscape(expiration.Format("2006-01-02_15:04:05")))
	}
	httpResponse, err := request.Post("/shorten")
	if err != nil {
		return "", err
	}
	switch httpResponse.StatusCode() {
	case http.StatusOK:
		return shortendUrlFromBody(httpResponse)
	case http.StatusBadRequest:
		return "", errorFromBody(httpResponse.Body())
	default:
		return "", errors.New("unexpected error")
	}
}

func shortendUrlFromBody(httpResponse *resty.Response) (string, error) {
	var response shortenResponse
	err := json.Unmarshal(httpResponse.Body(), &response)
	if err != nil {
		return "", err
	}
	return response.Shortened, nil
}

type unshortenResponse struct {
	Unshortened string `json:"unshortened"`
}

func (c HTTPClient) Unshorten(rawURL string) (string, error) {
	httpResponse, err := c.client.R().
		SetQueryParam("url", url.QueryEscape(rawURL)).
		Get("/unshorten")
	if err != nil {
		return "", err
	}
	switch httpResponse.StatusCode() {
	case http.StatusOK:
		return unshortendUrlFromBody(httpResponse)
	case http.StatusBadRequest:
		return "", errorFromBody(httpResponse.Body())
	default:
		return "", errors.New("unexpected error")
	}
}

func unshortendUrlFromBody(httpResponse *resty.Response) (string, error) {
	var response unshortenResponse
	err := json.Unmarshal(httpResponse.Body(), &response)
	if err != nil {
		return "", err
	}
	return response.Unshortened, nil
}

type errorBody struct {
	Error string `json:"error"`
}

func errorFromBody(body []byte) error {
	var e errorBody
	_ = json.Unmarshal(body, &e)
	switch e.Error {
	case ErrNotFound.Error():
		return ErrNotFound
	case ErrMissingScheme.Error():
		return ErrMissingScheme
	case ErrMissingHostname.Error():
		return ErrMissingHostname
	case ErrInvalidURL.Error():
		return ErrInvalidURL
	case ErrExpired.Error():
		return ErrExpired
	default:
		panic(fmt.Sprintf("unexpected error: %s", e.Error))
	}
	return nil
}

func NewHTTPClientFromResty(client *resty.Client) *HTTPClient {
	return &HTTPClient{client: client}
}

func NewHTTPClient() *HTTPClient {
	client := resty.New()
	return NewHTTPClientFromResty(client)
}
