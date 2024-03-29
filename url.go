package glip

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grokify/mogo/encoding/guid"
	"github.com/grokify/mogo/net/urlutil"
)

const (
	rxGlipWebhookV1CaptureFmt = `^(?i)https?://[^/]+/webhook/([^/#?]+)`    // #nosec G101
	rxGlipWebhookV2CaptureFmt = `^(?i)https?://[^/]+/webhook/v2/([^/#?]+)` // #nosec G101
)

type WebhookURL struct {
	webhookID       string
	originalVersion int
	originalInput   string
}

func MustNewWebhookURLString(input string, webhookVersion int) string {
	hookURL, err := NewWebhookURL(input)
	if err != nil {
		return input
	}
	if webhookVersion == 2 {
		return hookURL.V2URL()
	}
	return hookURL.V1URL()
}

func NewWebhookURL(input string) (WebhookURL, error) {
	input = strings.TrimSpace(input)
	wu := WebhookURL{
		originalInput: input}
	r1 := regexp.MustCompile(rxGlipWebhookV2CaptureFmt)
	m1 := r1.FindStringSubmatch(input)
	if len(m1) > 0 {
		wu.webhookID = m1[1]
		wu.originalVersion = 2
		return wu, nil
	}
	r2 := regexp.MustCompile(rxGlipWebhookV1CaptureFmt)
	m2 := r2.FindStringSubmatch(input)
	if len(m2) > 0 {
		wu.webhookID = m2[1]
		wu.originalVersion = 1
		return wu, nil
	}
	if urlutil.IsHTTP(input, true, true) {
		return wu, fmt.Errorf("is not Glip webhook URL [%s]", input)
	}

	if strings.Contains(input, "/") ||
		strings.Contains(input, "#") ||
		strings.Contains(input, "?") ||
		strings.Contains(input, " ") {
		return wu, fmt.Errorf("id has chars [/#? ] [%s]", input)
	}

	return WebhookURL{
		webhookID:     input,
		originalInput: input}, nil
}

func (w *WebhookURL) IsGUID() bool {
	return guid.ValidGUIDHex(w.webhookID)
}

func (w *WebhookURL) V1URL() string {
	return GlipWebhookV1BaseURLProduction + w.webhookID
}

func (w *WebhookURL) V2URL() string {
	return GlipWebhookV2BaseURLProduction + w.webhookID
}

func (w *WebhookURL) OriginalInput() string {
	return w.originalInput
}

func (w *WebhookURL) OriginalVersion() int {
	return w.originalVersion
}

func (w *WebhookURL) ID() string {
	return w.webhookID
}
