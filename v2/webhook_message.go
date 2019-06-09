package v2

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"
)

const (
	GlipWebhookBaseURLProduction string = "https://hooks.glip.com/webhook/v2/"
	GlipWebhookBaseURLSandbox    string = "https://hooks-glip.devtest.ringcentral.com/webhook/v2/"
	HTTPMethodPost               string = "POST"
	webhookV2Path                string = "/webhook/v2/"
	rxGlipWebhookV2Pattern       string = `^https?://[^/]+/webhook/v2/[^/]+/?$`
	rxGlipWebhookV1Pattern       string = `^(?i)(https?://[^/]+)/webhook/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})/?$`
)

var rxGlipWebhookV1 = regexp.MustCompile(rxGlipWebhookV1Pattern)
var rxGlipWebhookV2 = regexp.MustCompile(rxGlipWebhookV2Pattern)

func ToWebhookV2Uri(input string) string {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return input
	}
	if strings.Index(input, "/") == -1 {
		return GlipWebhookBaseURLProduction + input
	}
	if rxGlipWebhookV2.MatchString(input) {
		return input
	}
	m := rxGlipWebhookV1.FindStringSubmatch(input)
	if len(m) == 3 {
		return m[1] + webhookV2Path + m[2]
	}
	return input
}

type GlipWebhookMessage struct {
	Activity    string       `json:"activity,omitempty"`
	IconEmoji   string       `json:"iconEmoji,omitempty"`
	IconUri     string       `json:"iconUri,omitempty"`
	Text        string       `json:"text,omitempty"`
	Title       string       `json:"title,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Author       Author  `json:"author,omitempty"`
	Color        string  `json:"color,omitempty"`
	Fallback     string  `json:"fallback,omitempty"`
	Fields       []Field `json:"fields,omitempty"`
	Footnote     string  `json:"footnote,omitempty"`
	ImageUri     string  `json:"imageUri,omitempty"`
	Intro        string  `json:"intro,omitempty"`
	Text         string  `json:"text,omitempty"`
	ThumbnailUri string  `json:"thumbnailUri,omitempty"`
	Title        string  `json:"title,omitempty"`
	Type         string  `json:"card,omitempty"`
}

type Author struct {
	Name    string `json:"name,omitempty"`
	IconUri string `json:"iconUri,omitempty"`
	Uri     string `json:"uri,omitempty"`
}

type Field struct {
	Style string `json:"style,omitempty"`
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
}

type Footnote struct {
	IconUri string    `json:"iconUri,omitempty"`
	Text    string    `json:"text,omitempty"`
	Time    time.Time `json:"time,omitempty"`
}

type GlipWebhookResponse struct {
	Status  string           `json:"status,omitempty"`
	Message string           `json:"message,omitempty"`
	Error   GlipWebhookError `json:"error,omitempty"`
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
