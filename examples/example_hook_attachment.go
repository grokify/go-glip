package examples

import (
	"time"

	glipwebhook "github.com/grokify/go-glip"
)

func ExampleHookBodyAttachment() glipwebhook.GlipWebhookMessage {
	return glipwebhook.GlipWebhookMessage{
		Icon:     "https://i.imgur.com/9yILi61.png",
		Activity: "Activity of the post ♠♥♣♦",
		Title:    "**Title of the post ♠♥♣♦**",
		Body:     "Body of the post ♠♥♣♦",
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
				ThumbnailURL: "https://raw.githubusercontent.com/grokify/go-glip/master/docs/example_thumbnail-url.png",
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
				FooterIcon: "https://raw.githubusercontent.com/grokify/go-glip/master/docs/example_footer-icon.png",
				TS:         time.Now().Unix(),
			},
		},
	}
}
