package examples

import (
	glipwebhook "github.com/grokify/go-glip"
)

func ExampleHookBodyAlert() glipwebhook.GlipWebhookMessage {
	text := ":warning: 4 devices in **San Diego** have gone **Offline** :warning:"
	msg := glipwebhook.GlipWebhookMessage{
		Icon: "https://i.imgur.com/9yILi61.png",
		Body: text,
		Attachments: []glipwebhook.Attachment{
			{
				//Text: ":warning: 4 devices in **San Diego have gone **Offline** :warning:",
				Fields: []glipwebhook.Field{
					{
						Title: "Alert Name",
						Value: "San Diego Office Devices",
						Short: false},
					{
						Title: "Target",
						Value: "San Diego",
						Short: true},
					{
						Title: "Alert Trigger",
						Value: "# of Devices went offline",
						Short: true},
					{
						Title: "Condition",
						Value: "More than 3",
						Short: true},
					{
						Title: "Triggered Value",
						Value: "4 devices",
						Short: true},
					{
						Title: "Report Link",
						Value: "https://www.analytics.ringcentral.com/devices-offline",
						Short: false},
				},
			},
		},
	}
	return msg
}
