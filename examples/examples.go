package examples

import (
	"github.com/grokify/go-glip"
)

type ExampleWebhook struct {
	Stub    string
	Message glip.GlipWebhookMessage
}

func ExampleWebhooks() []ExampleWebhook {
	return []ExampleWebhook{
		{
			Stub:    "alert",
			Message: ExampleHookBodyAlert()},
		{
			Stub:    "attachment",
			Message: ExampleHookBodyAttachment()},
		{
			Stub:    "salesforce",
			Message: ExampleHookBodySalesforce()},
		{
			Stub:    "simple",
			Message: ExampleHookBodySimple()},
	}
}
