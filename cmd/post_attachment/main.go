package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/grokify/go-glip"
	ru "github.com/grokify/go-ringcentral-client/office/v1/util"
	"github.com/grokify/go-ringcentral-client/office/v1/util/glipgroups"
	"github.com/grokify/goauth/credentials"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
	"github.com/rs/zerolog/log"

	"github.com/grokify/go-glip/examples"
)

// main finds Glip groups matching the following command:
// find_team -group "My Group Name"
func main() {
	err := config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env")
	logutil.FatalErr(err)

	var credsFile, credsKey, wantGroupName, filepath string
	flag.StringVar(&credsFile, "credsfile", "/path/to/rccreds.json", "RC Creds File")
	flag.StringVar(&credsKey, "credskey", "PROD", "RC Creds Key")
	flag.StringVar(&wantGroupName, "group", "All Employees", "Glip Group Name")
	flag.StringVar(&filepath, "file", "/path/to/myfile.png", "Filepath")
	flag.Parse()

	credsSet, err := credentials.ReadFileCredentialsSet(credsFile, true)
	logutil.FatalErr(err)

	creds, err := credsSet.Get(credsKey)
	logutil.FatalErr(err)

	httpClient, err := creds.NewClient(context.Background())
	logutil.FatalErr(err)

	apiClient, err := ru.NewApiClientHttpClientBaseURL(
		httpClient, creds.OAuth2.ServerURL)
	logutil.FatalErr(err)

	set, err := glipgroups.NewGroupsSetApiRequest(
		httpClient, creds.OAuth2.ServerURL, "Team")
	logutil.FatalErr(err)

	log.Printf("Searching %v Teams\n", len(set.GroupsMap))

	group, err := set.FindGroupByName(wantGroupName)
	logutil.FatalErr(err)

	fmt.Printf("Found Team [%v]\n", wantGroupName)

	if 1 == 1 {
		info, resp, err := apiClient.GlipApi.CreatePost(
			context.Background(), group.ID, examples.ExamplePostBodyAlertWarning())
		if err != nil {
			log.Fatal().Err(err)
		} else if resp.StatusCode >= 300 {
			log.Fatal().Msg(fmt.Sprintf("Status [%v]", resp.StatusCode))
		}
		fmtutil.MustPrintJSON(info)
		info, resp, err = apiClient.GlipApi.CreatePost(
			context.Background(), group.ID, examples.ExamplePostBodyAlertSOS())
		if err != nil {
			log.Fatal().Err(err)
		} else if resp.StatusCode >= 300 {
			log.Fatal().
				Int("status", resp.StatusCode).
				Msg("response")
		}
		fmtutil.MustPrintJSON(info)
	}

	if 1 == 0 {
		resp, err := glip.PostFile(httpClient, creds.OAuth2.ServerURL, group.ID, filepath)
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

/*
func postFile(client *http.Client, serverURL, groupId string, filepath string) (*http.Response, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	filepathParts := strings.Split(filepath, "/")
	filename := filepathParts[len(filepathParts)-1]

	uploadURL := urlutil.JoinAbsolute(serverURL, glip.APIPathGlipPosts)
	//uploadURL := ro.BuildURL(serverURL, "/glip/posts", true, url.Values{})

	req, err := http.NewRequest(http.MethodPost, uploadURL, file)
	if err != nil {
		return nil, nil
	}

	rs := regexp.MustCompile(`(.[^.]+)$`).FindStringSubmatch(filepath)
	if len(rs) < 2 {
		return nil, err
	}
	req.Header.Add(httputilmore.HeaderContentType, mime.TypeByExtension(rs[1]))
	req.Header.Add(httputilmore.HeaderContentDisposition, fmt.Sprintf(`attachment; filename="%s"`, filename))

	return client.Do(req)
}
*/
