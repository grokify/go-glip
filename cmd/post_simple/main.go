package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/grokify/go-ringcentral-client/office/v1/util/glipgroups"
	"github.com/grokify/goauth"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/net/http/httpsimple"
	flags "github.com/jessevdk/go-flags"
	"github.com/rs/zerolog/log"

	"github.com/grokify/go-glip/examples"
)

type Options struct {
	goauth.Options
	Group  string   `short:"g" long:"groupname" description:"Group Name" required:"true"`
	URL    string   `short:"U" long:"url" description:"URL"`
	Method string   `short:"X" long:"request" description:"Method"`
	Header []string `short:"H" long:"header" description:"HTTP Headers"`
	Body   string   `short:"d" long:"data" description:"HTTP Body"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal().Err(err).Msg("required properties not present")
		panic("Z")
	}
	fmtutil.MustPrintJSON(opts)

	creds, err := goauth.ReadCredentialsFromSetFile(
		opts.Options.CredsPath, opts.Options.Account, true)
	if err != nil {
		log.Fatal().Err(err).
			Str("credsPath", opts.Options.CredsPath).
			Str("accountKey", opts.Options.Account).
			Msg("cannot read credentials")
	}

	var httpClient *http.Client
	if opts.Options.UseCLI() {
		httpClient, err = creds.NewClientCLI(context.Background(), "mystate")
	} else {
		httpClient, err = creds.NewClient(context.Background())
	}
	if err != nil {
		log.Fatal().Err(err).
			Bool("useCLI", opts.Options.UseCLI()).
			Msg("creds.NewClient() or creds.NewClientCLI()")
	}

	set, err := glipgroups.NewGroupsSetApiRequest(
		httpClient, creds.OAuth2.ServerURL, "Team")
	if err != nil {
		log.Fatal().Err(err)
	}

	log.Printf("Searching %v Teams\n", len(set.GroupsMap))

	group, err := set.FindGroupByName(opts.Group)
	if err != nil {
		log.Fatal().Err(err)
	} else {
		fmt.Printf("Found Team [%v]\n", opts.Group)
	}

	sclient, err := creds.NewSimpleClientHTTP(httpClient)
	if err != nil {
		log.Fatal().Err(err).
			Msg("cannot create simpleclient")
	}

	postURL := fmt.Sprintf("/restapi/v1.0/glip/chats/%s/adaptive-cards", group.ID)
	sreq := httpsimple.Request{
		Method:   http.MethodPost,
		URL:      postURL,
		Body:     examples.ExamplePostBodyCardBytes(),
		BodyType: httpsimple.BodyTypeJSON}
	resp, err := sclient.Do(sreq)
	if err != nil {
		log.Fatal().Err(err).Msg("post request")
	} else if resp.StatusCode > 299 {
		body, err := io.ReadAll(resp.Body)
		log.Fatal().
			Err(err).
			Str("url", postURL).
			Str("body", string(body)).
			Int("statusCode", resp.StatusCode).
			Msg("post request with bad statusCode")
	}
	fmt.Println("DONE")
}
