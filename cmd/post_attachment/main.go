package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"regexp"
	"strings"

	ru "github.com/grokify/go-ringcentral-client/office/v1/util"
	"github.com/grokify/go-ringcentral-client/office/v1/util/glipgroups"
	"github.com/grokify/goauth/credentials"
	"github.com/grokify/simplego/config"
	"github.com/grokify/simplego/fmt/fmtutil"
	"github.com/grokify/simplego/net/urlutil"
	"github.com/rs/zerolog/log"

	"github.com/grokify/go-glip"
	"github.com/grokify/go-glip/examples"
)

// main finds Glip groups matching the following command:
// find_team -group "My Group Name"
func main() {
	err := config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env")
	if err != nil {
		panic(err)
	}

	var credsFile, credsKey, wantGroupName, filepath string
	flag.StringVar(&credsFile, "credsfile", "/path/to/rccreds.json", "RC Creds File")
	flag.StringVar(&credsKey, "credskey", "PROD", "RC Creds Key")
	flag.StringVar(&wantGroupName, "group", "All Employees", "Glip Group Name")
	flag.StringVar(&filepath, "file", "/path/to/myfile.png", "Filepath")
	flag.Parse()

	credsSet, err := credentials.ReadFileCredentialsSet(credsFile, true)
	if err != nil {
		log.Fatal().Err(err)
	}
	creds, err := credsSet.Get(credsKey)
	if err != nil {
		log.Fatal().Err(err)
	}
	httpClient, err := creds.NewClient(context.Background())
	if err != nil {
		log.Fatal().Err(err)
	}

	apiClient, err := ru.NewApiClientHttpClientBaseURL(
		httpClient, creds.OAuth2.ServerURL)
	if err != nil {
		log.Fatal().Err(err)
	}

	set, err := glipgroups.NewGroupsSetApiRequest(
		httpClient, creds.OAuth2.ServerURL, "Team")
	if err != nil {
		log.Fatal().Err(err)
	}

	log.Printf("Searching %v Teams\n", len(set.GroupsMap))

	group, err := set.FindGroupByName(wantGroupName)
	if err != nil {
		log.Fatal().Err(err)
	} else {
		fmt.Printf("Found Team [%v]\n", wantGroupName)
	}

	if 1 == 1 {
		info, resp, err := apiClient.GlipApi.CreatePost(
			context.Background(), group.ID, examples.ExamplePostBodyAlertWarning())
		if err != nil {
			log.Fatal().Err(err)
		} else if resp.StatusCode >= 300 {
			log.Fatal().Msg(fmt.Sprintf("Status [%v]", resp.StatusCode))
		}
		fmtutil.PrintJSON(info)
		info, resp, err = apiClient.GlipApi.CreatePost(
			context.Background(), group.ID, examples.ExamplePostBodyAlertSOS())
		if err != nil {
			log.Fatal().Err(err)
		} else if resp.StatusCode >= 300 {
			log.Fatal().
				Int("status", resp.StatusCode).
				Msg("response")
		}
		fmtutil.PrintJSON(info)
	}

	if 1 == 0 {
		resp, err := postFile(httpClient, creds.OAuth2.ServerURL, group.ID, filepath)
		if err != nil {
			log.Fatal().Err(err)
		}
		log.Printf("Status %v\n", resp.StatusCode)
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal().Err(err)
		}
		log.Info().
			Msg(string(bytes))
	}

	log.Info().Msg("DONE")
}

func postFile(client *http.Client, serverURL, groupId string, filepath string) (*http.Response, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return &http.Response{}, err
	}

	filepathParts := strings.Split(filepath, "/")
	filename := filepathParts[len(filepathParts)-1]

	uploadURL := urlutil.JoinAbsolute(serverURL, glip.ApiPathGlipPosts)
	//uploadURL := ro.BuildURL(serverURL, "/glip/posts", true, url.Values{})

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
