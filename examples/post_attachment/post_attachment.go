package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/grokify/gotilla/config"
	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/grokify/gotilla/net/httputil"
	"github.com/grokify/oauth2-util-go/services/ringcentral"
	"golang.org/x/oauth2"
)

// main finds Glip groups matching the following command:
// find_team -group "My Group Name"
func main() {
	var wantGroupName, filepath string
	flag.StringVar(&wantGroupName, "group", "All Employees", "Glip Group Name")
	flag.StringVar(&filepath, "file", "/path/to/myfile.png", "Filepath")
	flag.Parse()

	config.LoadDotEnv()

	err := ringcentral.SetHostnameForURL(os.Getenv("RC_SERVER_URL"))
	if err != nil {
		log.Fatal(err)
	}

	client := getRingCentralClient()

	set := getGroupsSet(client, "Team")

	log.Printf("Searching %v Teams\n", len(set.GroupsMap))

	groups := set.FindGroupsByName(wantGroupName)

	fmtutil.PrintJSON(groups)

	for i, group := range groups {
		log.Printf("%d) %v %v\n", i, group.ID, group.Name)
		if 1 == 1 {
			resp, err := postTestMessageWithAttachment(client, group.ID, filepath)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Status %v\n", resp.StatusCode)
			bytes, err := httputil.ResponseBody(resp)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("%v\n", string(bytes))
		}
	}

	log.Println("DONE")
}

func postFile(client *http.Client, groupId string, filepath string) (*http.Response, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return &http.Response{}, err
	}

	filepathParts := strings.Split(filepath, "/")
	filename := filepathParts[len(filepathParts)-1]

	uploadURL := ringcentral.BuildURL("/glip/posts", true, query)

	req, err := http.NewRequest("POST", uploadURL, file)
	if err != nil {
		return &http.Response{}, nil
	}

	rs := regexp.MustCompile(`(.[^.]+)$`).FindStringSubmatch(filepath)
	if len(rs) < 2 {
		return &http.Response{}, err
	}
	req.Header.Add("Content-Type", mime.TypeByExtension(rs[1]))
	req.Header.Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	return client.Do(req)
}

func getGroupsSet(client *http.Client, groupType string) GroupsSet {
	set := GroupsSet{GroupsMap: map[string]Group{}}

	query := url.Values{}
	query.Add("recordCount", "250")

	if len(groupType) > 0 {
		query.Add("type", groupType)
	}

	for {
		groupsURL := ringcentral.BuildURL("/glip/groups", true, query)
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

func getRingCentralClient() *http.Client {
	o2Config := &oauth2.Config{
		ClientID:     os.Getenv("RC_APP_KEY"),
		ClientSecret: os.Getenv("RC_APP_SECRET"),
		Endpoint:     ringcentral.NewEndpoint(ringcentral.Hostname)}

	tok, err := o2Config.PasswordCredentialsToken(
		oauth2.NoContext,
		os.Getenv("RC_USER_USERNAME"),
		os.Getenv("RC_USER_PASSWORD"))

	log.Println(tok.AccessToken)
	if err != nil {
		log.Fatal(err)
	}
	return o2Config.Client(oauth2.NoContext, tok)
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
