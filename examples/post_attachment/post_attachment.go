package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/grokify/simplego/config"
	"github.com/grokify/simplego/fmt/fmtutil"
	httputil "github.com/grokify/simplego/net/httputilmore"
	"github.com/pkg/errors"

	"github.com/grokify/go-glip/examples"
	ru "github.com/grokify/go-ringcentral/clientutil"
	"github.com/grokify/go-ringcentral/clientutil/glipgroups"
	ro "github.com/grokify/oauth2more/ringcentral"
)

var RingCentralServerURL = "https://platform.ringcentral.com"

// main finds Glip groups matching the following command:
// find_team -group "My Group Name"
func main() {
	err := config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env")
	if err != nil {
		panic(err)
	}

	var wantGroupName, filepath string
	flag.StringVar(&wantGroupName, "group", "All Employees", "Glip Group Name")
	flag.StringVar(&filepath, "file", "/path/to/myfile.png", "Filepath")
	flag.Parse()

	serverURL := os.Getenv("RINGCENTRAL_SERVER_URL")

	err = ro.SetHostnameForURL(serverURL)
	if err != nil {
		log.Fatal(errors.Wrap(err, "SetHostnameForURL"))
	}

	httpClient, err := ro.NewHttpClientEnvFlexStatic("")
	if err != nil {
		log.Fatal(errors.Wrap(err, "getHttpClientEnv"))
	}

	apiClient, err := ru.NewApiClientHttpClientBaseURL(httpClient, serverURL)
	if err != nil {
		log.Fatal(err)
	}

	set, err := glipgroups.NewGroupsSetApiRequest(httpClient, serverURL, "Team")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Searching %v Teams\n", len(set.GroupsMap))

	group, err := set.FindGroupByName(wantGroupName)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Found Team [%v]\n", wantGroupName)
	}

	if 1 == 1 {
		info, resp, err := apiClient.GlipApi.CreatePost(
			context.Background(), group.ID, examples.GetExamplePostAlertWarning())
		if err != nil {
			log.Fatal(err)
		} else if resp.StatusCode >= 300 {
			log.Fatal(fmt.Sprintf("Status [%v]", resp.StatusCode))
		}
		fmtutil.PrintJSON(info)
		info, resp, err = apiClient.GlipApi.CreatePost(
			context.Background(), group.ID, examples.GetExamplePostAlertSOS())
		if err != nil {
			log.Fatal(err)
		} else if resp.StatusCode >= 300 {
			log.Fatal(fmt.Sprintf("Status [%v]", resp.StatusCode))
		}
		fmtutil.PrintJSON(info)
	}

	if 1 == 0 {
		resp, err := postFile(httpClient, group.ID, filepath)
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

	log.Println("DONE")
}

func postFile(client *http.Client, groupId string, filepath string) (*http.Response, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return &http.Response{}, err
	}

	filepathParts := strings.Split(filepath, "/")
	filename := filepathParts[len(filepathParts)-1]

	uploadURL := ro.BuildURL(RingCentralServerURL, "/glip/posts", true, url.Values{})

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
