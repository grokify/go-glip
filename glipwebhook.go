package glipwebhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/grokify/gotilla/net/httputil"
	"github.com/valyala/fasthttp"
)

const (
	GLIP_WEBHOOK_BASE_URL = "https://hooks.glip.com/webhook/"
	CONTENT_TYPE_JSON     = "application/json"
	CONTENT_TYPE_HEADER   = "Content-Type"
	HTTP_METHOD_POST      = "POST"
)

type GlipWebhookClient struct {
	HttpClient *http.Client
	FastClient fasthttp.Client
	WebhookUrl string
}

func newGlipWebhookClientCore(urlOrGuid string) (GlipWebhookClient, error) {
	client := GlipWebhookClient{}
	if len(urlOrGuid) > 0 {
		url, err := client.BuildWebhookURL(urlOrGuid)
		if err != nil {
			return client, err
		}
		client.WebhookUrl = url
	}
	return client, nil
}

func NewGlipWebhookClient(urlOrGuid string) (GlipWebhookClient, error) {
	client, err := newGlipWebhookClientCore(urlOrGuid)
	if err != nil {
		return client, err
	}
	client.HttpClient = httputil.NewHttpClient()
	return client, nil
}

func NewGlipWebhookClientFast(urlOrGuid string) (GlipWebhookClient, error) {
	client, err := newGlipWebhookClientCore(urlOrGuid)
	if err != nil {
		return client, err
	}
	client.FastClient = fasthttp.Client{}
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
	return client.PostWebhook(client.WebhookUrl, message)
}

func (client *GlipWebhookClient) PostWebhook(url string, message GlipWebhookMessage) (*http.Response, error) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return &http.Response{}, err
	}

	req, err := http.NewRequest(HTTP_METHOD_POST, url, bytes.NewBuffer(messageBytes))
	if err != nil {
		return &http.Response{}, err
	}

	req.Header.Set(CONTENT_TYPE_HEADER, CONTENT_TYPE_JSON)
	return client.HttpClient.Do(req)
}

func (client *GlipWebhookClient) PostWebhookGUID(guid string, message GlipWebhookMessage) (*http.Response, error) {
	return client.PostWebhook(strings.Join([]string{GLIP_WEBHOOK_BASE_URL, guid}, ""), message)
}

// Request using fasthttp
// Recycle request and response using fasthttp.ReleaseRequest(req) and
// fasthttp.ReleaseResponse(resp)
func (client *GlipWebhookClient) PostMessageFast(message GlipWebhookMessage) (*fasthttp.Request, *fasthttp.Response, error) {
	return client.PostWebhookFast(client.WebhookUrl, message)
}

func (client *GlipWebhookClient) PostWebhookFast(url string, message GlipWebhookMessage) (*fasthttp.Request, *fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	bytes, err := json.Marshal(message)
	if err != nil {
		return req, resp, err
	}
	req.SetBody(bytes)

	req.Header.SetMethod(HTTP_METHOD_POST)
	req.Header.SetRequestURI(url)
	req.Header.Set(CONTENT_TYPE_HEADER, CONTENT_TYPE_JSON)

	err = client.FastClient.Do(req, resp)
	return req, resp, err
}

func (client *GlipWebhookClient) PostWebhookGUIDFast(guid string, message GlipWebhookMessage) (*fasthttp.Request, *fasthttp.Response, error) {
	return client.PostWebhookFast(strings.Join([]string{GLIP_WEBHOOK_BASE_URL, guid}, ""), message)
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
