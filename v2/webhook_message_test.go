package v2

import (
	"testing"
)

var toWebhookV2UriTests = []struct {
	v    string
	want string
}{
	{"11112222-3333-4444-5555-666677778888", "https://hooks.glip.com/webhook/v2/11112222-3333-4444-5555-666677778888"},
	{"https://hooks.glip.com/webhook/v2/11112222-3333-4444-5555-666677778888", "https://hooks.glip.com/webhook/v2/11112222-3333-4444-5555-666677778888"},
	{"https://hooks.glip.com/webhook/11112222-3333-4444-5555-666677778888", "https://hooks.glip.com/webhook/v2/11112222-3333-4444-5555-666677778888"},
	{"https://hooks.glip.com/webhook/11112222-3333-4444-5555-666677778888/", "https://hooks.glip.com/webhook/v2/11112222-3333-4444-5555-666677778888"},
	{"https://hooks.glip.com/WEBHOOK/11112222-3333-4444-5555-666677778888/", "https://hooks.glip.com/webhook/v2/11112222-3333-4444-5555-666677778888"},
}

func TestToWebhookV2Uri(t *testing.T) {
	for _, tt := range toWebhookV2UriTests {
		got := ToWebhookV2Uri(tt.v)
		if got != tt.want {
			t.Errorf("v2.ToV2WebhookUri(\"%v\") Mismatch: want[%v], got [%v]", tt.v, tt.want, got)
		}
	}
}
