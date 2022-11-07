package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/grokify/mogo/config"
	"github.com/valyala/fasthttp"

	glip "github.com/grokify/go-glip"
)

func main() {
	err := config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env")
	if err != nil {
		panic(err)
	}

	var hookURL string
	flag.StringVar(&hookURL, "hookurl", "https://hooks.glip.com/webhook/1111-deadbeef-8888", "Config file path")
	flag.Parse()

	client := glip.NewGlipWebhookClientFast(hookURL, 1)

	msg := glip.GlipWebhookMessage{
		Icon:  "https://d30y9cdsu7xlg0.cloudfront.net/png/6597-200.png",
		Title: "Jeff is having a Maple Bacon Coffee Porter",
		Body:  "* Location: [The Funky Buddha Lounge](http://www.thefunkybuddha.com)",
		Attachments: []glip.Attachment{
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
				Fields: []glip.Field{
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
		},
	}

	req, resp, err := client.PostMessageFast(msg)
	if err == nil {
		fmt.Println(string(resp.Body()))
	}
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
}
