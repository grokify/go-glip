package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/io/ioutil"
	"github.com/grokify/mogo/log/logutil"
	"github.com/grokify/mogo/net/http/httpsimple"
	flags "github.com/jessevdk/go-flags"
	"github.com/valyala/fasthttp"

	glip "github.com/grokify/go-glip"
	"github.com/grokify/go-glip/examples"
)

const (
	// DEFAULT_URL           = "https://hooks.glip.com/webhook/1111-deadbeef-8888"
	ExampleTypeCard       = "card"
	ExampleTypeSimple     = "simple"
	ExampleTypeAlert      = "alert"
	ExampleTypeAttachment = "attachment"
	ExampleTypeSalesforce = "salesforce"
)

type cliOptions struct {
	WebhookURLOrGUID string `short:"u" long:"url" description:"URL or GUID for Webhook" required:"true"`
	Type             string `short:"t" long:"type" description:"Type [simple,attachment,salesforce,alert]"`
	File             string `short:"f" long:"file" description:"File containing JSON to use for body"`
	Data             string `short:"d" long:"data" description:"JSON to use for body"`
}

func getBodyBytes(ctx context.Context, webhookURLOrGUID string, body []byte) error {
	resp, err := httpsimple.Request{
		Method:   http.MethodPost,
		URL:      webhookURLOrGUID,
		Body:     body,
		BodyType: httpsimple.BodyTypeJSON}.Do(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("STATUS [%d]\n", resp.StatusCode)
	fmt.Println(string(ioutil.ReadAllOrError(resp.Body)))
	return nil
}

func main() {
	opts := cliOptions{}
	_, err := flags.Parse(&opts)
	logutil.FatalErr(err)

	client := glip.NewGlipWebhookClientFast(opts.WebhookURLOrGUID, 1)

	if len(strings.TrimSpace(opts.File)) > 0 {
		bytes, err := os.ReadFile(opts.File)
		logutil.FatalErr(err)

		logutil.FatalErr(getBodyBytes(context.Background(), opts.WebhookURLOrGUID, bytes))
	} else if len(strings.TrimSpace(opts.Data)) > 0 {
		logutil.FatalErr(getBodyBytes(context.Background(), opts.WebhookURLOrGUID, []byte(opts.Data)))
	} else if opts.Type == ExampleTypeCard {
		logutil.FatalErr(getBodyBytes(context.Background(), opts.WebhookURLOrGUID,
			examples.ExampleHookBodyCardBytes()))
	} else {
		msgs := []glip.GlipWebhookMessage{}

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
			log.Fatalf("body type not found [%s]", opts.Type)
		}

		logutil.FatalErr(fmtutil.PrintJSON(msgs))
		_, err := fmt.Printf("NUM_MSGS [%v]\n", len(msgs))
		logutil.FatalErr(err)

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
