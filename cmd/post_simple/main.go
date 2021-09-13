package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/grokify/go-glip/examples"
	"github.com/grokify/go-ringcentral-client/office/v1/util/glipgroups"
	"github.com/grokify/oauth2more/credentials"
	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/simplego/net/http/httpsimple"
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog/log"
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
	fmtutil.PrintJSON(opts)

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
		httpClient, err = creds.NewClientCli("mystate")
	} else {
		httpClient, err = creds.NewClient()
	}
	if err != nil {
		log.Fatal().Err(err).
			Bool("useCLI", opts.UseCLI()).
			Msg("creds.NewClient() or creds.NewClientCLI()")
	}

	set, err := glipgroups.NewGroupsSetApiRequest(
		httpClient, creds.Application.ServerURL, "Team")
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

	sclient, err := creds.NewSimpleClient(httpClient)
	if err != nil {
		fmt.Println(string(err.Error()))
		log.Fatal().Err(err).
			Msg("cannot create simpleclient")
	}

	postUrl := fmt.Sprintf("/restapi/v1.0/glip/chats/%s/adaptive-cards", group.ID)
	sreq := httpsimple.SimpleRequest{
		Method: http.MethodPost,
		URL:    postUrl,
		Body:   examples.ExamplePostBodyCardBytes(),
		IsJSON: true}
	resp, err := sclient.Do(sreq)
	if err != nil {
		log.Fatal().Err(err).Msg("post request")
	} else if resp.StatusCode > 299 {
		body, err := io.ReadAll(resp.Body)
		log.Fatal().
			Err(err).
			Str("url", postUrl).
			Str("body", string(body)).
			Int("statusCode", resp.StatusCode).
			Msg("post request with bad statusCode")
	}
	fmt.Println("DONE")
}
