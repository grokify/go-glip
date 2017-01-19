package glipwebhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/grokify/gotilla/net/httputil"
)

const (
	GLIP_WEBHOOK_BASE_URL = "https://hooks.glip.com/webhook/"
)

type GlipWebhookClient struct {
	HttpClient *http.Client
	WebhookUrl string
}

func NewGlipWebhookClient(urlOrGuid string) (GlipWebhookClient, error) {
	client := GlipWebhookClient{}
	url, err := client.BuildWebhookURL(urlOrGuid)
	if err != nil {
		return client, err
	}
	client.WebhookUrl = url
	client.HttpClient = httputil.NewHttpClient()
	return client, nil
}

func (client *GlipWebhookClient) BuildWebhookURL(urlOrGuid string) (string, error) {
	if len(urlOrGuid) < 36 {
		return "", errors.New("Webhook GUID or URL is required.")
	}
	rx := regexp.MustCompile(`^[0-9A-Za-z-]+$`)
	match := rx.FindString(urlOrGuid)
	if len(match) > 0 {
		return fmt.Sprintf("%v%v", GLIP_WEBHOOK_BASE_URL, urlOrGuid), nil
	}
	return urlOrGuid, nil
}

func (client *GlipWebhookClient) SendMessage(message GlipWebhookMessage) ([]byte, error) {
	resp, err := client.PostMessage(message)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (client *GlipWebhookClient) PostMessage(message GlipWebhookMessage) (*http.Response, error) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return &http.Response{}, err
	}
	return client.PostJSON(client.WebhookUrl, messageBytes)
}

func (client *GlipWebhookClient) PostJSON(url string, bodyBytes []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return &http.Response{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	httpClient := httputil.NewHttpClient()
	return httpClient.Do(req)
}

type GlipWebhookMessage struct {
	Icon     string `json:"icon,omitempty"`
	Activity string `json:"activity,omitempty"`
	Title    string `json:"title",omitempty`
	Body     string `json:"body,omitempty"`
}

type GlipWebhookResponse struct {
	Status  string           `json:"status,omitempty"`
	Message string           `json:"message,omitempty"`
	Error   GlipWebhookError `json:error,omitempty`
}

type GlipWebhookError struct {
	Code           string                   `json:"code,omitempty"`
	Message        string                   `json:"message,omitempty"`
	HttpStatusCode int                      `json:"http_status_code,omitempty"`
	ResponseData   string                   `json:"response_data,omitempty"`
	Response       GlipWebhookErrorResponse `json:"response,omitempty"`
}

func (gwerr *GlipWebhookError) Inflate() {
	if len(gwerr.ResponseData) > 2 {
		res := GlipWebhookResponse{}
		err := json.Unmarshal([]byte(gwerr.ResponseData), &res)
		if err == nil {
			gwerr.Response = res
		}
	}
}

type GlipWebhookErrorResponse struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Validation bool   `json:validation"`
}
