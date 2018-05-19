package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	ru "github.com/grokify/go-ringcentral/clientutil"
	"github.com/grokify/gotilla/fmt/fmtutil"
	httputil "github.com/grokify/gotilla/net/httputilmore"
	ro "github.com/grokify/oauth2more/ringcentral"
	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
)

type Options struct {
	Group     string `short:"g" long:"group" description:"Group Name" required:"true"`
	LoadUsers []bool `short:"u" long:"users" description:"List Users"`
}

func loadEnv() error {
	envPaths := []string{}
	if len(os.Getenv("ENV_PATH")) > 0 {
		envPaths = append(envPaths, os.Getenv("ENV_PATH"))
	}
	return godotenv.Load(envPaths...)
}

// main finds Glip groups matching the following command:
// find_team -group "My Group Name"
func main() {
	if err := loadEnv(); err != nil {
		panic(err)
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

	httpClient, err := ro.NewClientPassword(
		ro.ApplicationCredentials{
			ServerURL:    os.Getenv("RINGCENTRAL_SERVER_URL"),
			ClientID:     os.Getenv("RINGCENTRAL_CLIENT_ID"),
			ClientSecret: os.Getenv("RINGCENTRAL_CLIENT_SECRET")},
		ro.PasswordCredentials{
			Username:  os.Getenv("RINGCENTRAL_USERNAME"),
			Extension: os.Getenv("RINGCENTRAL_EXTENSION"),
			Password:  os.Getenv("RINGCENTRAL_PASSWORD")})
	if err != nil {
		panic(err)
	}

	apiClient, err := ru.NewApiClientHttpClientBaseURL(
		httpClient, os.Getenv("RINGCENTRAL_SERVER_URL"))
	if err != nil {
		log.Fatal(err)
	}

	set := getGroupsSet(httpClient, "Team")

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
					log.Fatal(err)
				} else if resp.StatusCode >= 300 {
					log.Fatal(fmt.Sprintf("API RESP %v", resp.StatusCode))
				}
				fmtutil.PrintJSON(info)
			}
		}
	}

	log.Println("DONE")
}

func getGroupsSet(client *http.Client, groupType string) GroupsSet {
	set := GroupsSet{GroupsMap: map[string]Group{}}

	query := url.Values{}
	query.Add("recordCount", "250")

	if len(groupType) > 0 {
		query.Add("type", groupType)
	}

	for {
		groupsURL := ro.BuildURL(os.Getenv("RINGCENTRAL_SERVER_URL"), "/glip/groups", true, query)
		resp, err := client.Get(groupsURL)
		if err != nil {
			log.Fatal(err)
		}

		groupsResp, err := GetGroupsResponseFromHTTPResponse(resp)
		if err != nil {
			log.Fatal(err)
		}
		set.AddGroups(groupsResp.Records)

		if len(groupsResp.Navigation.PrevPageToken) > 0 {
			query.Add("pageToken", groupsResp.Navigation.PrevPageToken)
		} else {
			break
		}
	}
	return set
}

type GroupsSet struct {
	GroupsMap map[string]Group
}

func (set *GroupsSet) AddGroups(groups []Group) {
	for _, grp := range groups {
		set.GroupsMap[grp.ID] = grp
	}
}

func (set *GroupsSet) FindGroupsByName(groupName string) []Group {
	groups := []Group{}
	for _, group := range set.GroupsMap {
		if groupName == group.Name {
			groups = append(groups, group)
		}
	}
	return groups
}

type Group struct {
	ID               string    `json:"id,omitempty"`
	Name             string    `json:"name,omitempty"`
	Description      string    `json:"description,omitempty"`
	CreationTime     time.Time `json:"creationTime,omitempty"`
	LastModifiedTime time.Time `json:"lastModifiedTime,omitempty"`
	Members          []string  `json:"members,omitempty"`
}

type GetGroupsResponse struct {
	Records    []Group    `json:"records,omitempty"`
	Navigation Navigation `json:"navigation,omitempty"`
}

func GetGroupsResponseFromHTTPResponse(resp *http.Response) (GetGroupsResponse, error) {
	bytes, err := httputil.ResponseBody(resp)
	if err != nil {
		return GetGroupsResponse{}, err
	}
	var apiResp GetGroupsResponse
	err = json.Unmarshal(bytes, &apiResp)
	return apiResp, err
}

type Navigation struct {
	PrevPageToken string `json:"prevPageToken,omitempty"`
	NextPageToken string `json:"nextPageToken,omitempty"`
}
