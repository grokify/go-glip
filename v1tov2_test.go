package glip

import (
	"testing"
)

var v1ToV2WewbhookUriTests = []struct {
	v    string
	want string
}{
	{"11112222-3333-4444-5555-666677778888", "https://hooks.glip.com/webhook/v2/11112222-3333-4444-5555-666677778888"},
	{"https://hooks.glip.com/webhook/v2/11112222-3333-4444-5555-666677778888", "https://hooks.glip.com/webhook/v2/11112222-3333-4444-5555-666677778888"},
	{"https://hooks.glip.com/webhook/11112222-3333-4444-5555-666677778888", "https://hooks.glip.com/webhook/v2/11112222-3333-4444-5555-666677778888"},
	{"https://hooks.glip.com/webhook/11112222-3333-4444-5555-666677778888/", "https://hooks.glip.com/webhook/v2/11112222-3333-4444-5555-666677778888"},
	{"https://hooks.glip.com/WEBHOOK/11112222-3333-4444-5555-666677778888/", "https://hooks.glip.com/webhook/v2/11112222-3333-4444-5555-666677778888"},
}

func TestV1ToV2WewbhookUri(t *testing.T) {
	for _, tt := range v1ToV2WewbhookUriTests {
		got := V1ToV2WewbhookUri(tt.v)
		if got != tt.want {
			t.Errorf("glipwebhook.V1ToV2WewbhookUri(\"%v\") Mismatch: want[%v], got [%v]", tt.v, tt.want, got)
		}
	}
}
