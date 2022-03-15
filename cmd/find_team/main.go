package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	ru "github.com/grokify/go-ringcentral-client/office/v1/util"
	"github.com/grokify/goauth/credentials"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
	"github.com/jessevdk/go-flags"

	"github.com/grokify/go-ringcentral-client/office/v1/util/glipgroups"
	"github.com/grokify/go-ringcentral-client/office/v1/util/mergedusers"
)

type Options struct {
	credentials.Options
	Groups    []string `short:"g" long:"group" description:"Group Name" required:"true"`
	LoadUsers []bool   `short:"u" long:"users" description:"List Users"`
}

// main finds Glip groups matching the following command:
// find_team -group "My Group Name"
func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	logutil.FatalErr(err)

	fmtutil.MustPrintJSON(opts)

	creds, err := credentials.ReadCredentialsFromFile(opts.CredsPath, opts.Account, true)
	logutil.FatalErr(err)

	var httpClient *http.Client

	if len(opts.CLI) > 0 {
		httpClient, err = creds.NewClientCli("mystate")
	} else {
		httpClient, err = creds.NewClient(context.Background())
	}
	logutil.FatalErr(err)

	serverURL := "https://platform.ringcentral.com"

	apiClient, err := ru.NewApiClientHttpClientBaseURL(
		httpClient, serverURL)
	if err != nil {
		log.Fatal(err)
	}

	set, err := glipgroups.NewGroupsSetApiRequest(httpClient, serverURL, "Team")
	logutil.FatalErr(err)

	log.Printf("Searching %v Teams\n", len(set.GroupsMap))

	for _, g := range opts.Groups {
		groups := set.FindGroupsByName(g)

		fmtutil.MustPrintJSON(groups)

		if len(opts.LoadUsers) > 0 {
			for _, group := range groups {
				memberCount := len(group.Members)
				for i, memberID := range group.Members {
					n := i + 1
					fmt.Printf("[%v/%v] %v", n, memberCount, memberID)
					info, resp, err := apiClient.GlipApi.LoadPerson(context.Background(), memberID)
					if err != nil {
						log.Fatal(fmt.Sprintf("API ERR: %v\n", err))
					} else if resp.StatusCode >= 300 {
						log.Fatal(fmt.Sprintf("API RESP %v\n", resp.StatusCode))
					}
					fmtutil.MustPrintJSON(info)
				}
			}
		}

		if 1 == 0 {
			for _, group := range groups {
				set, err := mergedusers.NewMergedUsersApiIds(httpClient,
					serverURL,
					group.Members)
				logutil.FatalErr(err)

				//fmtutil.PrintJSON(set)
				for id, user := range set.MergedUserMap {
					thin := user.ToMergedUserThin()
					fmt.Printf("ID [%v] NAME [%v][%v][%v]\n", id,
						user.GlipPersonInfo.FirstName,
						user.GlipPersonInfo.LastName,
						thin.DisplayNumber)
					fmtutil.MustPrintJSON(thin)
				}
				if 1 == 0 {
					fmtutil.MustPrintJSON(set.MergedUserMap["557601020"])
					user := set.MergedUserMap["557601020"]
					thin := user.ToMergedUserThin()
					fmtutil.MustPrintJSON(thin)
				}
				fmt.Printf("NUM_USERS [%v]\n", len(set.MergedUserMap))
				break
			}
		}
	}

	log.Println("DONE")
}
