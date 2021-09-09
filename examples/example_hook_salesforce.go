package examples

import (
	glipwebhook "github.com/grokify/go-glip"
)

func ExampleHookBodySalesforce() glipwebhook.GlipWebhookMessage {
	msg := glipwebhook.GlipWebhookMessage{
		Icon:  "http://www.bridgethegap.com/wp-content/uploads/2017/02/salesforce-best-practices.png",
		Title: "**Top Opportunities**",
		Body:  "Full report: https://my.salesforce.com/00O80000007MOgS",
		Attachments: []glipwebhook.Attachment{
			{
				Color:        "#00ff2a",
				ThumbnailURL: "https://funkybuddhabrewery.com/sites/default/files/WorldBeerCupGold.png",
				Fields: []glipwebhook.Field{
					{
						Title: "Opportunity", Short: true,
						Value: "[Eureka Oil and Gas (0038000001MgG2z)](https://example.com)"},
					{
						Title: "Owner", Short: true,
						Value: "Alice Collins"},
					{
						Title: "Stage", Short: true,
						Value: "Proposal/Quote"},
					{
						Title: "Close Date", Short: true,
						Value: "2017-09-20"},
					{
						Title: "Amount", Short: true,
						Value: "$750,000"},
					{
						Title: "Probability", Short: true,
						Value: "85%"},
				},
			},
			{
				Color:        "#dfdd13",
				ThumbnailURL: "https://funkybuddhabrewery.com/sites/default/files/WorldBeerCupGold.png",
				Fields: []glipwebhook.Field{
					{
						Title: "Opportunity", Short: true,
						Value: "[Pacfic Restaurants (0038000004WhM4a)](https://example.com)"},
					{
						Title: "Owner", Short: true,
						Value: "Justin Lyons"},
					{
						Title: "Stage", Short: true,
						Value: "Discovery"},
					{
						Title: "Close Date", Short: true,
						Value: "2017-09-20"},
					{
						Title: "Amount", Short: true,
						Value: "$500,000"},
					{
						Title: "Probability", Short: true,
						Value: "70%"},
				},
			},
			{
				Color:        "#ff0000",
				ThumbnailURL: "https://funkybuddhabrewery.com/sites/default/files/WorldBeerCupGold.png",
				Fields: []glipwebhook.Field{
					{
						Title: "Opportunity", Short: true,
						Value: "[Electric Company of America (0038000004OrS7y)](https://example.com)"},
					{
						Title: "Owner", Short: true,
						Value: "Matthew West"},
					{
						Title: "Stage", Short: true,
						Value: "Proposal/Quote"},
					{
						Title: "Close Date", Short: true,
						Value: "2017-12-20"},
					{
						Title: "Amount", Short: true,
						Value: "$150,000"},
					{
						Title: "Probability", Short: true,
						Value: "40%"},
				},
			},
		},
	}
	return msg
}
