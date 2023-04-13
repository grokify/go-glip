package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/grokify/goauth"
	ro "github.com/grokify/goauth/ringcentral"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"

	glip "github.com/grokify/go-glip"
)

// main finds Glip groups matching the following command:
// find_team -group "My Group Name"
func main() {
	_, err := config.LoadDotEnv([]string{os.Getenv("ENV_PATH"), "./.env"}, 1)
	if err != nil {
		panic(err)
	}

	var wantGroupName, filepath string
	flag.StringVar(&wantGroupName, "group", "All Employees", "Glip Group Name")
	flag.StringVar(&filepath, "file", "/path/to/myfile.png", "Filepath")
	flag.Parse()

	fmt.Printf("GROUP %v\n", wantGroupName)
	fmt.Printf("FILE %v\n", filepath)

	httpClient, err := ro.NewClientPassword(
		goauth.NewCredentialsOAuth2Env("RINGCENTRAL_"))
	if err != nil {
		panic(err)
	}

	set, err := getGroupsSet(httpClient, "Team")
	logutil.FatalErr(err)

	log.Printf("Searching %v Teams for %v\n", len(set.GroupsMap), wantGroupName)

	groups := set.FindGroupsByName(wantGroupName)
	fmtutil.MustPrintJSON(groups)

	for i, group := range groups {
		log.Printf("%d) %v %v\n", i, group.ID, group.Name)

		resp, err := glip.PostFile(httpClient,
			os.Getenv("RINGCENTRAL_SERVER_URL"),
			group.ID, filepath)
		logutil.FatalErr(err)

		log.Printf("Status %v\n", resp.StatusCode)
		bytes, err := io.ReadAll(resp.Body)
		logutil.FatalErr(err)

		log.Printf("%v\n", string(bytes))
	}

	log.Println("DONE")
}

func getGroupsSet(client *http.Client, groupType string) (GroupsSet, error) {
	set := GroupsSet{GroupsMap: map[string]Group{}}

	query := url.Values{}
	query.Add("recordCount", "250")

	if len(groupType) > 0 {
		query.Add("type", groupType)
	}

	for {
		groupsURL, err := ro.BuildURL(os.Getenv("RINGCENTRAL_SERVER_URL"), "/glip/groups", true, query)
		if err != nil {
			return set, err
		}
		resp, err := client.Get(groupsURL)
		if err != nil {
			return set, err
		}

		groupsResp, err := GetGroupsResponseFromHTTPResponse(resp)
		if err != nil {
			return set, err
		}
		set.AddGroups(groupsResp.Records)

		if len(groupsResp.Navigation.PrevPageToken) > 0 {
			query.Add("pageToken", groupsResp.Navigation.PrevPageToken)
		} else {
			break
		}
	}
	return set, nil
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
	bytes, err := io.ReadAll(resp.Body)
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
