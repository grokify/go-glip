package glip

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/grokify/mogo/mime/mimeutil"
	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/grokify/mogo/net/urlutil"
)

func PostFile(client *http.Client, serverURL, groupID, path string) (*http.Response, error) {
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

	req.Header.Add(httputilmore.HeaderContentType,
		mimeutil.MustTypeByFile(path, true))
	req.Header.Add(httputilmore.HeaderContentDisposition,
		fmt.Sprintf(`attachment; filename="%s"`, filename))

	return client.Do(req)
}
