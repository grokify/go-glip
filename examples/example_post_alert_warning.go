package examples

import (
	rc "github.com/grokify/go-ringcentral-client/office/v1/client"
)

func ExamplePostBodyAlertWarning() rc.GlipCreatePost {
	reqBody := rc.GlipCreatePost{
		Text: ":warning: **4 devices** in **San Diego** have gone **Offline** :warning:",
		Attachments: []rc.GlipMessageAttachmentInfoRequest{
			{
				Type:  "Card",
				Color: "#ffa500",
				Fields: []rc.GlipMessageAttachmentFieldsInfo{
					{
						Title: "Alert Name",
						Value: "San Diego Office Devices",
						Style: "Long"},
					{
						Title: "Target",
						Value: "San Diego",
						Style: "Short"},
					{
						Title: "Alert Trigger",
						Value: "# of Devices went offline",
						Style: "Short"},
					{
						Title: "Condition",
						Value: "More than 3",
						Style: "Short"},
					{
						Title: "Triggered Value",
						Value: "4 devices",
						Style: "Short"},
					{
						Title: "Report Link",
						Value: "https://www.analytics.ringcentral.com/devices-offline",
						Style: "Long"},
				},
			},
		},
	}
	return reqBody
}
