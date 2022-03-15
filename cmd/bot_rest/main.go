package main

import (
	"context"
	"fmt"
	"os"

	"github.com/grokify/goauth"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/log/logutil"
	"github.com/grokify/mogo/net/httputilmore"
)

func main() {
	err := config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env")
	logutil.FatalErr(err)

	client, err := goauth.NewClientTokenJSON(
		context.Background(),
		[]byte(os.Getenv("RINGCENTRAL_TOKEN_JSON")),
	)
	logutil.FatalErr(err)

	url := "https://platform.ringcentral.com/restapi/v1.0/glip/persons/~"

	resp, err := client.Get(url)
	logutil.FatalErr(err)

	logutil.FatalErr(httputilmore.PrintResponse(resp, true))
	fmt.Println("DONE")
}
