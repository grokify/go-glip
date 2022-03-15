package glip

import (
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"github.com/grokify/mogo/net/httputilmore"
	"github.com/grokify/mogo/net/urlutil"
)

func PostFile(client *http.Client, serverURL, groupID string, path string) (*http.Response, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	_, filename := filepath.Split(path)
	query := url.Values{}
	query.Add("groupId", groupID)
	query.Add("name", filename)

	uploadURL, err := urlutil.URLAddQueryValuesString(
		urlutil.JoinAbsolute(serverURL, APIPathGlipFiles),
		query)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, uploadURL.String(), file)
	if err != nil {
		return nil, err
	}

	rs := regexp.MustCompile(`(.[^.]+)$`).FindStringSubmatch(path)
	if len(rs) < 2 {
		return nil, err
	}
	req.Header.Add(httputilmore.HeaderContentType, mime.TypeByExtension(rs[1]))
	req.Header.Add(httputilmore.HeaderContentDisposition, fmt.Sprintf(`attachment; filename="%s"`, filename))

	return client.Do(req)
}
