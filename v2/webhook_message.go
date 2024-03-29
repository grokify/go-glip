package v2

import (
	"encoding/json"
	"time"
)

const (
	GlipWebhookBaseURLProduction string = "https://hooks.glip.com/webhook/v2/"                     // #nosec G101
	GlipWebhookBaseURLSandbox    string = "https://hooks-glip.devtest.ringcentral.com/webhook/v2/" // #nosec G101
	HTTPMethodPost               string = "POST"
	FieldStyleLong               string = "Long"
	FieldStyleShort              string = "Short"
)

type GlipWebhookMessage struct {
	Activity    string       `json:"activity,omitempty"`
	IconEmoji   string       `json:"iconEmoji,omitempty"`
	IconURI     string       `json:"iconUri,omitempty"`
	Text        string       `json:"text,omitempty"`
	Title       string       `json:"title,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

func NewGlipWebhookMessage() GlipWebhookMessage {
	return GlipWebhookMessage{Attachments: []Attachment{}}
}

type Attachment struct {
	Author       *Author   `json:"author,omitempty"`
	Color        string    `json:"color,omitempty"`
	Fallback     string    `json:"fallback,omitempty"`
	Fields       []Field   `json:"fields,omitempty"`
	Footnote     *Footnote `json:"footnote,omitempty"`
	ImageURI     string    `json:"imageUri,omitempty"`
	Intro        string    `json:"intro,omitempty"`
	Text         string    `json:"text,omitempty"`
	ThumbnailURI string    `json:"thumbnailUri,omitempty"`
	Title        string    `json:"title,omitempty"`
	Type         string    `json:"type,omitempty"`
}

type Author struct {
	Name    string `json:"name,omitempty"`
	IconURI string `json:"iconUri,omitempty"`
	URI     string `json:"uri,omitempty"`
}

type Field struct {
	Style string `json:"style,omitempty"` // ['Short','Long']
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
}

type Footnote struct {
	IconURI string    `json:"iconUri,omitempty"`
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
