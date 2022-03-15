package glip

import (
	"testing"
)

var v1ToV2WewbhookURITests = []struct {
	input string
	v2url string
}{
	{"11112222-3333-4444-5555-666677778888", "https://hooks.ringcentral.com/webhook/v2/11112222-3333-4444-5555-666677778888"},
	{"https://hooks.glip.com/webhook/v2/11112222-3333-4444-5555-666677778888", "https://hooks.ringcentral.com/webhook/v2/11112222-3333-4444-5555-666677778888"},
	{"https://hooks.glip.com/webhook/11112222-3333-4444-5555-666677778888", "https://hooks.ringcentral.com/webhook/v2/11112222-3333-4444-5555-666677778888"},
	{"https://hooks.glip.com/webhook/11112222-3333-4444-5555-666677778888/", "https://hooks.ringcentral.com/webhook/v2/11112222-3333-4444-5555-666677778888"},
	{"https://hooks.glip.com/WEBHOOK/11112222-3333-4444-5555-666677778888/", "https://hooks.ringcentral.com/webhook/v2/11112222-3333-4444-5555-666677778888"},
}

func TestV1ToV2WewbhookURI(t *testing.T) {
	for _, tt := range v1ToV2WewbhookURITests {
		got, err := V1ToV2WewbhookURI(tt.input)
		if err != nil {
			t.Errorf("glipwebhook.V1ToV2WewbhookUri(\"%s\") Error [%s]", tt.input, err.Error())
		}
		if got != tt.v2url {
			t.Errorf("glipwebhook.V1ToV2WewbhookUri(\"%s\") Mismatch: want[%v], got [%v]", tt.input, tt.v2url, got)
		}
	}
}
