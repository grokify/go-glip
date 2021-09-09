package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"

	glipwebhook "github.com/grokify/go-glip"
	"github.com/grokify/go-glip/examples"
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

type cliOptions struct {
	WebhookUrlOrGuid string `short:"u" long:"url" description:"URL or GUID for Webhook" required:"true"`
	Type             string `short:"t" long:"type" description:"Type [simple,card,salesforce,alert]" required:"true"`
}

func main() {
	opts := cliOptions{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	client, err := glipwebhook.NewGlipWebhookClientFast(opts.WebhookUrlOrGuid)
	if err != nil {
		log.Fatal(err)
	}

	msgs := []glipwebhook.GlipWebhookMessage{}

	opts.Type = strings.ToLower(strings.TrimSpace(opts.Type))
	switch opts.Type {
	case "simple":
		msgs = append(msgs, examples.ExampleHookBodySimple())
	case "salesforce":
		msgs = append(msgs, examples.ExampleHookBodySalesforce())
	case "alert":
		msgs = append(msgs, examples.ExampleHookBodyAlert())
	case "card":
		msgs = append(msgs, examples.ExampleHookBodyAttachment())
	default:
		log.Fatal(fmt.Sprintf("body type not found [%s]", opts.Type))
	}

	fmtutil.PrintJSON(msgs)
	fmt.Printf("NUM_MSGS [%v]\n", len(msgs))

	for _, msg := range msgs {
		req, resp, err := client.PostMessageFast(msg)
		if err == nil {
			fmt.Println(string(resp.Body()))
		}
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}
}
