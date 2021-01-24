package main

import (
	"context"
	"fmt"
	"log"
	"os"

	om "github.com/grokify/oauth2more"
	"github.com/grokify/simplego/config"
	hum "github.com/grokify/simplego/net/httputilmore"
)

func main() {
	err := config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env")
	if err != nil {
		log.Fatal(err)
	}

	client, err := om.NewClientTokenJSON(
		context.Background(),
		[]byte(os.Getenv("RINGCENTRAL_TOKEN_JSON")),
	)
	if err != nil {
		log.Fatal(err)
	}

	url := "https://platform.ringcentral.com/restapi/v1.0/glip/persons/~"

	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	hum.PrintResponse(resp, true)
	fmt.Println("DONE")
}
