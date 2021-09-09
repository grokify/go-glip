package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"regexp"

	ru "github.com/grokify/go-ringcentral-client/office/v1/util"
	"github.com/grokify/oauth2more/credentials"
	ro "github.com/grokify/oauth2more/ringcentral"
	"github.com/grokify/simplego/config"
	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/simplego/os/osutil"
	"github.com/jessevdk/go-flags"

	"github.com/grokify/go-ringcentral-client/office/v1/util/glipgroups"
	"github.com/grokify/go-ringcentral-client/office/v1/util/mergedusers"
)

type Options struct {
	Group     string `short:"g" long:"group" description:"Group Name" required:"true"`
	LoadUsers []bool `short:"u" long:"users" description:"List Users"`
}

// main finds Glip groups matching the following command:
// find_team -group "My Group Name"
func main() {
	if err := config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env"); err != nil {
		log.Fatal(err)
	}

	/*
		var wantGroupName string
		flag.StringVar(&wantGroupName, "group", "All Employees", "Glip Group Name")
		flag.Parse()
	*/
	opts := Options{}

	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("USER [%v]\n", os.Getenv("RINGCENTRAL_USERNAME"))

	if 1 == 0 {
		fmtutil.PrintJSON(osutil.EnvFiltered(regexp.MustCompile(`RINGCENTRAL`)))
	}

	httpClient, err := ro.NewClientPassword(
		credentials.ApplicationCredentials{
			ServerURL:    os.Getenv("RINGCENTRAL_SERVER_URL"),
			ClientID:     os.Getenv("RINGCENTRAL_CLIENT_ID"),
			ClientSecret: os.Getenv("RINGCENTRAL_CLIENT_SECRET")},
		credentials.PasswordCredentials{
			Username: os.Getenv("RINGCENTRAL_USERNAME"),
			Password: os.Getenv("RINGCENTRAL_PASSWORD")})
	if err != nil {
		log.Fatal(fmt.Printf("AUTH: %v\n", err))
	}

	apiClient, err := ru.NewApiClientHttpClientBaseURL(
		httpClient, os.Getenv("RINGCENTRAL_SERVER_URL"))
	if err != nil {
		log.Fatal(err)
	}

	set, err := glipgroups.NewGroupsSetApiRequest(httpClient, os.Getenv("RINGCENTRAL_SERVER_URL"), "Team")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Searching %v Teams\n", len(set.GroupsMap))

	groups := set.FindGroupsByName(opts.Group)

	fmtutil.PrintJSON(groups)

	if len(opts.LoadUsers) > 0 {
		for _, group := range groups {
			memberCount := len(group.Members)
			for i, memberId := range group.Members {
				n := i + 1
				fmt.Printf("[%v/%v] %v", n, memberCount, memberId)
				info, resp, err := apiClient.GlipApi.LoadPerson(context.Background(), memberId)
				if err != nil {
					log.Fatal(fmt.Sprintf("API ERR: %v\n", err))
				} else if resp.StatusCode >= 300 {
					log.Fatal(fmt.Sprintf("API RESP %v\n", resp.StatusCode))
				}
				fmtutil.PrintJSON(info)
			}
		}
	}

	if 1 == 1 {
		for _, group := range groups {
			set, err := mergedusers.NewMergedUsersApiIds(httpClient,
				os.Getenv("RINGCENTRAL_SERVER_URL"),
				group.Members)
			if err != nil {
				log.Fatal(err)
			}
			//fmtutil.PrintJSON(set)
			for id, user := range set.MergedUserMap {
				thin := user.ToMergedUserThin()
				fmt.Printf("ID [%v] NAME [%v][%v][%v]\n", id,
					user.GlipPersonInfo.FirstName,
					user.GlipPersonInfo.LastName,
					thin.DisplayNumber)
				fmtutil.PrintJSON(thin)
			}
			if 1 == 0 {
				fmtutil.PrintJSON(set.MergedUserMap["557601020"])
				user := set.MergedUserMap["557601020"]
				thin := user.ToMergedUserThin()
				fmtutil.PrintJSON(thin)
			}
			fmt.Printf("NUM_USERS [%v]\n", len(set.MergedUserMap))
			break
		}

	}

	log.Println("DONE")
}
