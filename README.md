Glip Webhook Client in Go
=========================

## Installation

```bash
$ go get github.com/grokify/glip-go-webhook
```

## Usage

1. Create a webhook URL for a conversation in Glip
2. Use the code below to send a message to the webhook URL

```go
import (
    "fmt"
    "github.com/grokify/glip-go-webhook"
)

func sendMessage() {
    // Can instantiate webhook client with full URL or GUID only
    url := "https://hooks.glip.com/webhook/00001111-2222-3333-4444-555566667777"
    client, err := glipwebhook.NewGlipWebhookClient(url)
    if err != nil {
        panic("BAD URL")
    }

    msg := glipwebhook.GlipWebhookMessage{
        Icon:     "https://raw.githubusercontent.com/grokify/glip-go-webhook/master/glip_gopher_600x600xfff.png",
        Activity: "Gopher [Bot]",
        Title:    "Test Message Title",
        Body:     "Test Message Body"}

    resp, err := client.PostMessage(msg)

    respBodyBytes, err := client.SendMessage(msg)
    if err == nil {
        fmt.Printf("%v\n", string(respBodyBytes))
    }
}
```

### Using `fasthttp` client

Posts can be made using [`fasthttp`](https://github.com/valyala/fasthttp).

```go
import (
    "fmt"
    "github.com/grokify/glip-go-webhook"
)

func sendMessage() {
    // Can instantiate webhook client with full URL or GUID only
    url := "https://hooks.glip.com/webhook/00001111-2222-3333-4444-555566667777"
    client, err := glipwebhook.NewGlipWebhookClient(url)
    if err != nil {
        panic("BAD URL")
    }

    msg := glipwebhook.GlipWebhookMessage{
        Body: "Test Message Body"}

    res, resp, err := client.PostMessageFast(msg)
    if err == nil {
        fmt.Println(string(resp.Body()))
    }
    fasthttp.ReleaseRequest(req)
    fasthttp.ReleaseResponse(resp)
}
```
