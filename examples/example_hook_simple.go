package examples

import (
	glipwebhook "github.com/grokify/go-glip"
)

func ExampleHookBodySimple() glipwebhook.GlipWebhookMessage {
	return glipwebhook.GlipWebhookMessage{
		Icon:  "https://d30y9cdsu7xlg0.cloudfront.net/png/6597-200.png",
		Title: "Jeff is having a Maple Bacon Coffee Porter",
		Body:  "* Location: [The Funky Buddha Lounge](http://www.thefunkybuddha.com)",
		Attachments: []glipwebhook.Attachment{
			{
				Title:        "Maple Bacon Coffee Porter",
				TitleLink:    "https://funkybuddhabrewery.com/our-beers/little-buddha-small-batch/maple-bacon-coffee-porter",
				Color:        "#ff0000",
				AuthorName:   "Funky Buddha Lounge",
				AuthorLink:   "https://funkybuddhabrewery.com",
				AuthorIcon:   "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSymLVWmoonPAlBJWD68MczjvfLUavibXovYisb7vADlM8V7Z47DA",
				Text:         "*The beer that started it all.* Evoking a complete diner-style breakfast in a glass, Maple Bacon Coffee Porter is a complex beer with a multitude of flavors at play. It pours an opaque ebony brew with a frothy tan head.",
				ImageURL:     "https://funkybuddhabrewery.com/sites/default/files/our_beer/MBCP_2017_bottle-mock-22oz.png",
				ThumbnailURL: "https://funkybuddhabrewery.com/sites/default/files/WorldBeerCupGold.png",
				Fields: []glipwebhook.Field{
					{
						Title: "Style",
						Value: "Porter",
						Short: true},
					{
						Title: "Beer Advocate Rating:",
						Value: "[99](http://tinyurl.com/psf4uzq)",
						Short: true},
				},
			},
			{
				Color:    "#00ff2a",
				Text:     "Come down and grab a beer!",
				ImageURL: "http://a.memegen.com/zkqt2e.gif",
			},
		},
	}
}
