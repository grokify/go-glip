package glip

import (
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/grokify/mogo/net/httputilmore"
	"github.com/grokify/mogo/net/urlutil"
)

func PostFile(client *http.Client, serverURL, groupID string, filepath string) (*http.Response, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	filepathParts := strings.Split(filepath, "/")
	filename := filepathParts[len(filepathParts)-1]

	query := url.Values{}
	query.Add("groupId", groupID)
	query.Add("name", filename)

	uploadURL, err := urlutil.URLAddQueryValuesString(
		urlutil.JoinAbsolute(
			serverURL, APIPathGlipFiles),
		query)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, uploadURL.String(), file)
	if err != nil {
		return nil, err
	}

	rs := regexp.MustCompile(`(.[^.]+)$`).FindStringSubmatch(filepath)
	if len(rs) < 2 {
		return nil, err
	}
	req.Header.Add(httputilmore.HeaderContentType, mime.TypeByExtension(rs[1]))
	req.Header.Add(httputilmore.HeaderContentDisposition, fmt.Sprintf(`attachment; filename="%s"`, filename))

	return client.Do(req)
}
