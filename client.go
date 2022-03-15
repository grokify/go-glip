package glip

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	v2 "github.com/grokify/go-glip/v2"
	"github.com/grokify/mogo/net/httputilmore"
	"github.com/valyala/fasthttp"
)

type GlipWebhookClient struct {
	HTTPClient     *http.Client
	FastClient     *fasthttp.Client
	WebhookURL     string
	webhookVersion int
}

func newGlipWebhookClientCore(urlOrGUID string, webhookVersion int) GlipWebhookClient {
	if webhookVersion != 2 {
		webhookVersion = 1
	}
	client := GlipWebhookClient{webhookVersion: webhookVersion}
	if len(urlOrGUID) > 0 {
		client.WebhookURL = MustNewWebhookURLString(urlOrGUID, webhookVersion)
	}
	return client
}

func NewGlipWebhookClient(urlOrGUID string, webhookVersion int) GlipWebhookClient {
	client := newGlipWebhookClientCore(urlOrGUID, webhookVersion)
	client.HTTPClient = httputilmore.NewHTTPClient()
	return client
}

func NewGlipWebhookClientFast(urlOrGUID string, webhookVersion int) GlipWebhookClient {
	client := newGlipWebhookClientCore(urlOrGUID, webhookVersion)
	client.FastClient = &fasthttp.Client{}
	return client
}

func (client *GlipWebhookClient) PostMessage(message GlipWebhookMessage) (*http.Response, error) {
	return client.PostWebhook(client.WebhookURL, message)
}

func (client *GlipWebhookClient) PostWebhook(url string, message GlipWebhookMessage) (*http.Response, error) {
	if client.webhookVersion == 2 {
		v2url, err := V1ToV2WewbhookURI(url)
		if err != nil {
			return nil, err
		}
		return client.PostWebhookV2(v2url, webhookBodyV1ToV2(message))
	}
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	return httputilmore.DoJSONSimple(client.HTTPClient, http.MethodPost, url, map[string][]string{}, msgBytes)
}

func (client *GlipWebhookClient) PostWebhookV1Bytes(url string, message []byte) (*http.Response, error) {
	return httputilmore.DoJSONSimple(client.HTTPClient, http.MethodPost, url, map[string][]string{}, message)
}

func (client *GlipWebhookClient) PostWebhookV2(url string, message v2.GlipWebhookMessage) (*http.Response, error) {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	return httputilmore.DoJSONSimple(client.HTTPClient, http.MethodPost, url, map[string][]string{}, msgBytes)
}

func (client *GlipWebhookClient) PostWebhookGUID(guid string, message GlipWebhookMessage) (*http.Response, error) {
	return client.PostWebhook(strings.Join([]string{WebhookBaseURL, guid}, ""), message)
}

// Request using fasthttp
// Recycle request and response using fasthttp.ReleaseRequest(req) and
// fasthttp.ReleaseResponse(resp)
func (client *GlipWebhookClient) PostMessageFast(message GlipWebhookMessage) (*fasthttp.Request, *fasthttp.Response, error) {
	return client.PostWebhookFast(client.WebhookURL, message)
}

func (client *GlipWebhookClient) PostWebhookFast(url string, message GlipWebhookMessage) (*fasthttp.Request, *fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	if client.webhookVersion == 2 {
		var err error
		url, err = V1ToV2WewbhookURI(url)
		if err != nil {
			return req, resp, err
		}
		bytes, err := json.Marshal(webhookBodyV1ToV2(message))
		if err != nil {
			return req, resp, err
		}
		req.SetBody(bytes)
	} else {
		bytes, err := json.Marshal(message)
		if err != nil {
			return req, resp, err
		}
		req.SetBody(bytes)
	}

	req.Header.SetRequestURI(url)
	req.Header.SetMethod(http.MethodPost)
	req.Header.Set(httputilmore.HeaderContentType, httputilmore.ContentTypeAppJSONUtf8)

	err := client.FastClient.Do(req, resp)
	return req, resp, err
}

func (client *GlipWebhookClient) PostWebhookGUIDFast(guidOrURL string, message GlipWebhookMessage) (*fasthttp.Request, *fasthttp.Response, error) {
	return client.PostWebhookFast(MustNewWebhookURLString(guidOrURL, client.webhookVersion), message)
}

func webhookBodyV1ToV2(v1msg GlipWebhookMessage) v2.GlipWebhookMessage {
	v2msg := v2.GlipWebhookMessage{
		Activity:    v1msg.Activity,
		IconURI:     v1msg.Icon,
		Text:        v1msg.Body,
		Title:       v1msg.Title,
		Attachments: []v2.Attachment{}}
	for _, v1att := range v1msg.Attachments {
		v2msg.Attachments = append(v2msg.Attachments, attachmentV1ToV2(v1att))
	}

	return v2msg
}

func attachmentV1ToV2(v1att Attachment) v2.Attachment {
	v2att := v2.Attachment{
		Color:        v1att.Color,
		Fields:       []v2.Field{},
		ImageURI:     v1att.ImageURL,
		Intro:        v1att.Pretext,
		Text:         v1att.Text,
		ThumbnailURI: v1att.ThumbnailURL,
		Title:        v1att.Title,
		Type:         v1att.Type,
	}
	if len(strings.TrimSpace(v2att.Type)) == 0 {
		v2att.Type = "Card"
	}
	if len(v1att.AuthorName) > 0 {
		v2att.Author = &v2.Author{
			Name:    v1att.AuthorName,
			IconURI: v1att.AuthorIcon,
			URI:     v1att.AuthorLink}
	}
	if len(strings.TrimSpace(v1att.FooterIcon)) > 0 || len(strings.TrimSpace(v1att.Footer)) > 0 {
		v2att.Footnote = &v2.Footnote{
			IconURI: v1att.FooterIcon,
			Text:    v1att.Footer,
		}
		if v1att.TS > 0 {
			v2att.Footnote.Time = time.Unix(v1att.TS, 0)
		}
	}

	for _, v1field := range v1att.Fields {
		v2field := v2.Field{
			Title: v1field.Title,
			Value: v1field.Value}
		if v1field.Short {
			v2field.Style = v2.FieldStyleShort
		} else {
			v2field.Style = v2.FieldStyleLong
		}
		v2att.Fields = append(v2att.Fields, v2field)
	}

	return v2att
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
