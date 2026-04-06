# Glip Webhook Client in Go

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/grokify/go-glip/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/grokify/go-glip/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/grokify/go-glip/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/grokify/go-glip/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/grokify/go-glip/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/grokify/go-glip/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/go-glip
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/go-glip
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/go-glip
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/go-glip
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fgo-glip
 [loc-svg]: https://tokei.rs/b1/github/grokify/go-glip
 [repo-url]: https://github.com/grokify/go-glip
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/go-glip/blob/master/LICENSE

## Installation

```bash
$ go get github.com/grokify/go-glip
```

## Usage

1. Create a webhook URL for a conversation in Glip
2. Use the code below to send a message to the webhook URL

```go
import (
    "fmt"
    "github.com/grokify/go-glip"
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
    "github.com/grokify/go-glip"
)

func sendMessage() {
    // Can instantiate webhook client with full URL or GUID only
    url := "https://hooks.glip.com/webhook/00001111-2222-3333-4444-555566667777"
    client, err := glipwebhook.NewGlipWebhookClientFast(url)
    if err != nil {
        panic("BAD URL")
    }

    msg := glipwebhook.GlipWebhookMessage{
        Body: "Test Message Body"}

    req, resp, err := client.PostMessageFast(msg)
    if err == nil {
        fmt.Println(string(resp.Body()))
    }
    fasthttp.ReleaseRequest(req)
    fasthttp.ReleaseResponse(resp)
}
```

You can reuse the client for different Webhook URLs or GUIDs as follows:

```go
// Webhook URL
res, resp, err := client.PostWebhookFast(url, msg)

// Webhook GUID
res, resp, err := client.PostWebhookGUIDFast(guid, msg)
```
