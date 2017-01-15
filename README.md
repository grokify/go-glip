Glip Webhook Client in Go
=========================

## Installation

```bash
$ go get github.com/grokify/glip-go-webhook
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/grokify/glip-go-webhook"
)

func main() {
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

    bytes, err := client.SendMessage(msg)
    if err != nil {
        panic("BAD RESPONSE")
    }
    fmt.Printf("%v\n", string(bytes))

    fmt.Println("DONE")
}
```
