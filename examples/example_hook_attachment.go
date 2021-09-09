package examples

import (
	"time"

	glipwebhook "github.com/grokify/go-glip"
)

func ExampleHookBodyAttachment() glipwebhook.GlipWebhookMessage {
	msg := glipwebhook.GlipWebhookMessage{
		Icon:  "https://i.imgur.com/9yILi61.png",
		Title: "**Title of the post ♠♥♣♦**",
		Body:  "Body of the post ♠♥♣♦",
		Attachments: []glipwebhook.Attachment{
			{
				Title:        "Attachment Title ♠♥♣♦",
				TitleLink:    "https://example.com/title_link",
				Color:        "#00ff2a",
				AuthorName:   "Author Name ♠♥♣♦",
				AuthorLink:   "https://example.com/author_link",
				AuthorIcon:   "https://upload.wikimedia.org/wikipedia/commons/thumb/f/fd/000080_Navy_Blue_Square.svg/1200px-000080_Navy_Blue_Square.svg.png",
				Text:         "Attachment text ♠♥♣♦",
				Pretext:      "Attachment pretext appears before the attachment block ♠♥♣♦",
				ImageURL:     "https://media3.giphy.com/media/l4FssTixISsPStXRC/giphy.gif",
				ThumbnailURL: "https://funkybuddhabrewery.com/sites/default/files/WorldBeerCupGold.png",
				Fields: []glipwebhook.Field{
					{
						Title: "Field 1 ♠♥♣♦",
						Value: "A short field ♠♥♣♦",
						Short: true},
					{
						Title: "Field 2",
						Value: "This is [a linked short field](https://example.com)",
						Short: true},
					{
						Title: "Field 3 ♠♥♣♦",
						Value: "A long, full-width field with *formatting* and [a link](https://example.com) \n\n ♠♥♣♦",
						Short: false},
				},
				Footer:     "Attachment footer and timestamp ♠♥♣♦",
				FooterIcon: "http://www.iconsdb.com/icons/preview/red/square-ios-app-xxl.png",
				TS:         time.Now().Unix(),
			},
		},
	}
	if 1 == 0 {
		msg.Icon = "https://example.com/post_icon.png"
		msg.Attachments[0].ImageURL = "https://example.com/congrats.gif"
		msg.Attachments[0].FooterIcon = "https://example.com/footer_icon.png"
		msg.Attachments[0].AuthorIcon = "https://example.com/author_icon.png"
	}
	return msg
}
