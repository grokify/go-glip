package glipwebhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/grokify/mogo/net/httputilmore"
)

const (
	GlipWebhookBaseURL = "https://hooks.glip.com/webhook/" // #nosec G101
)

type GlipWebhookClient struct {
	HTTPClient *http.Client
	WebhookURL string
}

func NewGlipWebhookClient(urlOrGUID string) (GlipWebhookClient, error) {
	client := GlipWebhookClient{}
	url, err := client.BuildWebhookURL(urlOrGUID)
	if err != nil {
		return client, err
	}
	client.WebhookURL = url
	client.HTTPClient = httputilmore.NewHTTPClient()
	return client, nil
}

func (client *GlipWebhookClient) BuildWebhookURL(urlOrGUID string) (string, error) {
	if len(urlOrGUID) < 36 {
		return "", errors.New("webhook GUID or URL is required")
	}
	rx := regexp.MustCompile(`^[0-9A-Za-z-]+$`)
	match := rx.FindString(urlOrGUID)
	if len(match) > 0 {
		return fmt.Sprintf("%v%v", GlipWebhookBaseURL, urlOrGUID), nil
	}
	return urlOrGUID, nil
}

func (client *GlipWebhookClient) SendMessage(message GlipWebhookMessage) ([]byte, error) {
	resp, err := client.PostMessage(message)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (client *GlipWebhookClient) PostMessage(message GlipWebhookMessage) (*http.Response, error) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return &http.Response{}, err
	}
	return client.PostJSON(client.WebhookURL, messageBytes)
}

func (client *GlipWebhookClient) PostJSON(url string, bodyBytes []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return &http.Response{}, err
	}
	req.Header.Set(httputilmore.HeaderContentType, httputilmore.ContentTypeAppJSONUtf8)
	httpClient := httputilmore.NewHTTPClient()
	return httpClient.Do(req)
}

type GlipWebhookMessage struct {
	Icon     string `json:"icon,omitempty"`
	Activity string `json:"activity,omitempty"`
	Title    string `json:"title,omitempty"`
	Body     string `json:"body,omitempty"`
}

type GlipWebhookResponse struct {
	Status  string           `json:"status,omitempty"`
	Message string           `json:"message,omitempty"`
	Error   GlipWebhookError `json:"error,omitempty"`
}

type GlipWebhookError struct {
	Code           string                   `json:"code,omitempty"`
	Message        string                   `json:"message,omitempty"`
	HTTPStatusCode int                      `json:"http_status_code,omitempty"`
	ResponseData   string                   `json:"response_data,omitempty"`
	Response       GlipWebhookErrorResponse `json:"response,omitempty"`
}

func (gwerr *GlipWebhookError) Inflate() {
	if len(gwerr.ResponseData) > 2 {
		res := GlipWebhookErrorResponse{}
		err := json.Unmarshal([]byte(gwerr.ResponseData), &res)
		if err == nil {
			gwerr.Response = res
		}
	}
}

type GlipWebhookErrorResponse struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Validation bool   `json:"validation"`
}
