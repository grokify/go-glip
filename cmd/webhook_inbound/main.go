package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/simplego/io/ioutilmore"
	"github.com/grokify/simplego/net/http/httpsimple"
	"github.com/jessevdk/go-flags"
	"github.com/valyala/fasthttp"

	glipwebhook "github.com/grokify/go-glip"
	"github.com/grokify/go-glip/examples"
)

const (
	DEFAULT_URL           = "https://hooks.glip.com/webhook/1111-deadbeef-8888"
	ExampleTypeCard       = "card"
	ExampleTypeSimple     = "simple"
	ExampleTypeAlert      = "alert"
	ExampleTypeAttachment = "attachment"
	ExampleTypeSalesforce = "salesforce"
)

type cliOptions struct {
	WebhookUrlOrGuid string `short:"u" long:"url" description:"URL or GUID for Webhook" required:"true"`
	Type             string `short:"t" long:"type" description:"Type [simple,attachment,salesforce,alert]"`
	File             string `short:"f" long:"file" description:"File containing JSON to use for body"`
	Data             string `short:"d" long:"data" description:"JSON to use for body"`
}

func getBodyBytes(webhookUrlOrGuid string, body []byte) {
	resp, err := httpsimple.Do(httpsimple.SimpleRequest{
		Method: http.MethodPost,
		URL:    webhookUrlOrGuid,
		Body:   body,
		IsJSON: true})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("STATUS [%d]\n", resp.StatusCode)
	fmt.Println(string(ioutilmore.ReadAllOrError(resp.Body)))
}

func main() {
	opts := cliOptions{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	client, err := glipwebhook.NewGlipWebhookClientFast(opts.WebhookUrlOrGuid, 1)
	if err != nil {
		log.Fatal(err)
	}

	if len(strings.TrimSpace(opts.File)) > 0 {
		bytes, err := os.ReadFile(opts.File)
		if err != nil {
			log.Fatal(err)
		}
		getBodyBytes(opts.WebhookUrlOrGuid, bytes)
	} else if len(strings.TrimSpace(opts.Data)) > 0 {
		getBodyBytes(opts.WebhookUrlOrGuid, []byte(opts.Data))
	} else if opts.Type == ExampleTypeCard {
		getBodyBytes(opts.WebhookUrlOrGuid,
			examples.ExampleHookBodyCardBytes())
	} else {
		msgs := []glipwebhook.GlipWebhookMessage{}

		opts.Type = strings.ToLower(strings.TrimSpace(opts.Type))
		switch opts.Type {
		case ExampleTypeSimple:
			msgs = append(msgs, examples.ExampleHookBodySimple())
		case ExampleTypeSalesforce:
			msgs = append(msgs, examples.ExampleHookBodySalesforce())
		case ExampleTypeAlert:
			msgs = append(msgs, examples.ExampleHookBodyAlert())
		case ExampleTypeAttachment:
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
}
