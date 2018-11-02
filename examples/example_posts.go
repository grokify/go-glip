package examples

import (
	rc "github.com/grokify/go-ringcentral/client"
)

func GetExamplePostAlertWarning() rc.GlipCreatePost {
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

func GetExamplePostAlertSOS() rc.GlipCreatePost {
	reqBody := rc.GlipCreatePost{
		Text: ":sos: **6 devices** in **San Diego** have gone **Offline** :sos:",
		Attachments: []rc.GlipMessageAttachmentInfoRequest{
			{
				Type:  "Card",
				Color: "#ff0000",
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
						Value: "More than 5",
						Style: "Short"},
					{
						Title: "Triggered Value",
						Value: "6 devices",
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
