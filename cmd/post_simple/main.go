package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/grokify/go-ringcentral-client/office/v1/util/glipgroups"
	"github.com/grokify/goauth/credentials"
	"github.com/grokify/gohttp/httpsimple"
	"github.com/grokify/mogo/fmt/fmtutil"
	flags "github.com/jessevdk/go-flags"
	"github.com/rs/zerolog/log"

	"github.com/grokify/go-glip/examples"
)

type Options struct {
	credentials.Options
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

	creds, err := credentials.ReadCredentialsFromFile(
		opts.CredsPath, opts.Account, true)
	if err != nil {
		log.Fatal().Err(err).
			Str("credsPath", opts.CredsPath).
			Str("accountKey", opts.Account).
			Msg("cannot read credentials")
	}

	var httpClient *http.Client
	if opts.UseCLI() {
		httpClient, err = creds.NewClientCLI("mystate")
	} else {
		httpClient, err = creds.NewClient(context.Background())
	}
	if err != nil {
		log.Fatal().Err(err).
			Bool("useCLI", opts.UseCLI()).
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
	sreq := httpsimple.SimpleRequest{
		Method: http.MethodPost,
		URL:    postURL,
		Body:   examples.ExamplePostBodyCardBytes(),
		IsJSON: true}
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
