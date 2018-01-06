package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/grokify/go-glip"
	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
)

const (
	DEFAULT_URL = "https://hooks.glip.com/webhook/1111-deadbeef-8888"
)

func loadEnv() error {
	envPaths := []string{}
	if len(os.Getenv("ENV_PATH")) > 0 {
		envPaths = append(envPaths, os.Getenv("ENV_PATH"))
	}
	return godotenv.Load(envPaths...)
}

func getPostSimple() glipwebhook.GlipWebhookMessage {
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

func getPostAttachment() glipwebhook.GlipWebhookMessage {
	msg := glipwebhook.GlipWebhookMessage{
		Icon:  "https://i.imgur.com/9yILi61.png",
		Title: "**Title of the post**",
		Body:  "Body of the post",
		Attachments: []glipwebhook.Attachment{
			{
				Title:        "Attachment Title",
				TitleLink:    "https://example.com/title_link",
				Color:        "#00ff2a",
				AuthorName:   "Author Name",
				AuthorLink:   "https://example.com/author_link",
				AuthorIcon:   "https://upload.wikimedia.org/wikipedia/commons/thumb/f/fd/000080_Navy_Blue_Square.svg/1200px-000080_Navy_Blue_Square.svg.png",
				Text:         "Attachment text",
				Pretext:      "Attachment pretext appears before the attachment block",
				ImageURL:     "https://media3.giphy.com/media/l4FssTixISsPStXRC/giphy.gif",
				ThumbnailURL: "https://funkybuddhabrewery.com/sites/default/files/WorldBeerCupGold.png",
				Fields: []glipwebhook.Field{
					{
						Title: "Field 1",
						Value: "A short field",
						Short: true},
					{
						Title: "Field 2",
						Value: "This is [a linked short field](https://example.com)",
						Short: true},
					{
						Title: "Field 3",
						Value: "A long, full-width field with *formatting* and [a link](https://example.com) \n\n http://example.com/hello/world",
						Short: true},
				},
				Footer:     "Attachment footer and timestamp",
				FooterIcon: "http://www.iconsdb.com/icons/preview/red/square-ios-app-xxl.png", //},
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

func getPostSalesforce() glipwebhook.GlipWebhookMessage {
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

func main() {
	var hookUrl string
	flag.StringVar(&hookUrl, "hookurl", "", "Config file path")
	flag.Parse()

	if len(hookUrl) == 0 {
		if err := loadEnv(); err != nil {
			panic(err)
		}
		hookUrl = os.Getenv("GLIP_WEBHOOK_URL")
	}

	client, err := glipwebhook.NewGlipWebhookClientFast(hookUrl)
	if err != nil {
		log.Fatal(err)
	}

	msgs := []glipwebhook.GlipWebhookMessage{}

	addPostSimple := true
	addPostAttachment := true
	addPostSalesforce := true

	if addPostSimple {
		msgs = append(msgs, getPostSimple())
	}
	if addPostAttachment {
		msgs = append(msgs, getPostAttachment())
	}
	if addPostSalesforce {
		msgs = append(msgs, getPostSalesforce())
	}

	fmtutil.PrintJSON(msgs)

	for _, msg := range msgs {
		req, resp, err := client.PostMessageFast(msg)
		if err == nil {
			fmt.Println(string(resp.Body()))
		}
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}
}
