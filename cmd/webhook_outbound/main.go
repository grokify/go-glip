package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/grokify/go-glip"
	rc "github.com/grokify/go-ringcentral-client/office/v1/client"
	ru "github.com/grokify/go-ringcentral-client/office/v1/util"
	"github.com/grokify/goauth"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
	"github.com/rs/zerolog/log"
)

type RingCentralConfig struct {
	WebhookURL string `env:"RINGCENTRAL_WEBHOOK_URL"`
	AppPort    int64  `env:"PORT"`
}

type GenericEvent struct {
	UUID           string    `json:"uuid,omitempty"`
	Event          string    `json:"event,omitempty"`
	Timestamp      time.Time `json:"timestamp,omitempty"`
	SubscriptionId string    `json:"subscriptionId,omitempty"`
	OwnerId        string    `json:"ownerId,omitempty"`
}

type GlipEvent struct {
	UUID           string          `json:"uuid,omitempty"`
	Event          string          `json:"event,omitempty"`
	Timestamp      time.Time       `json:"timestamp,omitempty"`
	SubscriptionId string          `json:"subscriptionId,omitempty"`
	OwnerId        string          `json:"ownerId,omitempty"`
	Body           TextMessageBody `json:"body,omitempty"`
}

type TextMessageBody struct {
	Id               string             `json:"id,omitempty"`
	GroupId          string             `json:"groupId,omitempty"`
	Type             string             `json:"type,omitempty"`
	Text             string             `json:"text,omitempty"`
	CreatorId        string             `json:"creatorId,omitempty"`
	CreationTime     time.Time          `json:"creationTime,omitempty"`
	LastModifiedTime time.Time          `json:"lastModifiedTime,omitempty"`
	Mentions         []GlipEventMention `json:"mentions,omitempty"`
	EventType        string             `json:"eventType,omitempty"`
}

type GlipEventMention struct {
	Id   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

func WebhookHandler(res http.ResponseWriter, req *http.Request) {
	vtHeader := glip.HeaderValidationToken
	vt := req.Header.Get(vtHeader)
	if len(strings.TrimSpace(vt)) > 0 {
		res.Header().Set(vtHeader, vt)
		fmt.Printf("VALIDATION-TOKEN: %v", vt)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	logutil.FatalErr(err)

	evt := &GenericEvent{}
	err = json.Unmarshal(body, evt)
	logutil.FatalErr(err)

	subFmt := "/restapi/v1.0/subscription/.*threshold"
	rx := regexp.MustCompile(subFmt)
	m := rx.FindString(subFmt)
	if len(m) > 0 {
		//renew()
		return
	}

	//log.Info(jsonutil.MustMarshalString(req, true))
	//fmt.Println(string(body))

	log.Info().
		Str("body", string(body)).
		Msg("hook-body")
}

func getRingCentralAPIClient() (*rc.APIClient, error) {
	fmt.Println(os.Getenv("RINGCENTRAL_CORP_METABOT_TOKEN"))
	rcHttpClient, err := goauth.NewClientTokenJSON(
		context.Background(),
		[]byte(os.Getenv("RINGCENTRAL_CORP_METABOT_TOKEN")))
	if err != nil {
		return nil, err
	}

	return ru.NewApiClientHttpClientBaseURL(
		rcHttpClient,
		os.Getenv("RINGCENTRAL_SERVER_URL"),
	)
}

func createWebhook(webhookURL string) error {
	log.Info().Msg("Creating Hook...")
	apiClient, err := getRingCentralAPIClient()
	if err != nil {
		return err
	}

	req := rc.CreateSubscriptionRequest{
		EventFilters: []string{
			"/restapi/v1.0/glip/groups",
			"/restapi/v1.0/glip/posts",
			"/restapi/v1.0/subscription/~?threshold=60&interval=15",
		},
		DeliveryMode: rc.NotificationDeliveryModeRequest{
			TransportType: "WebHook",
			Address:       webhookURL,
		},
		//ExpiresIn: int32(ExpiresIn),
	}
	log.Info().Msg(jsonutil.MustMarshalString(req, true))

	info, resp, err := apiClient.PushNotificationsApi.CreateSubscription(
		context.Background(),
		req,
	)
	if err != nil {
		log.Fatal().Err(err)
	} else if resp.StatusCode >= 300 {
		log.Fatal().
			Int("status", resp.StatusCode).
			Msg("bad_status")
	}
	return fmtutil.PrintJSON(info)
}

func main() {
	var webhookURL string
	flag.StringVar(&webhookURL, "webhookurl", "", "Config file path")
	flag.Parse()

	err := config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env")
	logutil.FatalErr(err)

	appCfg := RingCentralConfig{}
	err = env.Parse(&appCfg)
	logutil.FatalErr(err)

	if len(webhookURL) == 0 {
		webhookURL = appCfg.WebhookURL
	}

	http.HandleFunc("/webhook", WebhookHandler)

	done := make(chan bool)
	go http.ListenAndServe(fmt.Sprintf(":%v", appCfg.AppPort), nil)
	log.Printf("Server started at port %v", appCfg.AppPort)
	//time.Sleep(3 * time.Second)
	//createWebhook(webhookURL)
	<-done
}
