package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/grokify/go-glip"
	rc "github.com/grokify/go-ringcentral-client/office/v1/client"
	ru "github.com/grokify/go-ringcentral-client/office/v1/util"
	"github.com/grokify/goauth"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
	"github.com/jessevdk/go-flags"
	//	"github.com/rs/zerolog/log"
)

type RingCentralConfig struct {
	WebhookURL string `env:"RINGCENTRAL_WEBHOOK_URL"`
	AppPort    int64  `env:"PORT"`
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	vt := r.Header.Get(glip.HeaderValidationToken)
	if len(strings.TrimSpace(vt)) > 0 {
		w.Header().Set(glip.HeaderValidationToken, vt)
		fmt.Printf("%s: %v", glip.HeaderValidationToken, vt)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	evt := &glip.GenericEvent{}
	err = json.Unmarshal(body, evt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	subFmt := "/restapi/v1.0/subscription/.*threshold"
	rx := regexp.MustCompile(subFmt)
	m := rx.FindString(subFmt)
	if len(m) > 0 {
		// renew()
		return
	}

	//log.Info(jsonutil.MustMarshalString(req, true))
	//fmt.Println(string(body))

	log.Printf("hook body [%s]", string(body))
}

func getRingCentralAPIClient() (*rc.APIClient, error) {
	fmt.Println(os.Getenv("RINGCENTRAL_CORP_METABOT_TOKEN"))
	rcHTTPClient, err := goauth.NewClientTokenJSON(
		context.Background(),
		[]byte(os.Getenv("RINGCENTRAL_CORP_METABOT_TOKEN")))
	if err != nil {
		return nil, err
	}

	return ru.NewApiClientHttpClientBaseURL(
		rcHTTPClient,
		os.Getenv("RINGCENTRAL_SERVER_URL"),
	)
}

func createWebhook(rcAPIClient *rc.APIClient, webhookURL string) error {
	log.Print("Creating Hook...")

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
	log.Print(jsonutil.MustMarshalString(req, true))

	info, resp, err := rcAPIClient.PushNotificationsApi.CreateSubscription(
		context.Background(),
		req,
	)
	if err != nil {
		log.Fatal(err)
	} else if resp.StatusCode >= 300 {
		log.Fatalf("bad_status [%d]", resp.StatusCode)
	}
	return fmtutil.PrintJSON(info)
}

type Options struct {
	WebhookURL    string `short:"w" long:"webhookurl" description:"Webhook URL"`
	CreateWebhook []bool `short:"c" long:"create" description:"Create webhook"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	logutil.FatalErr(err)

	err = config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env")
	logutil.FatalErr(err)

	appCfg := RingCentralConfig{}
	err = env.Parse(&appCfg)
	logutil.FatalErr(err)

	if len(opts.WebhookURL) == 0 {
		opts.WebhookURL = appCfg.WebhookURL
	}

	http.HandleFunc("/webhook", WebhookHandler)

	done := make(chan bool)
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%v", appCfg.AppPort), nil); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Printf("Server started at port %v", appCfg.AppPort)
	//time.Sleep(3 * time.Second)
	if len(opts.CreateWebhook) > 0 {
		rcAPIClient, err := getRingCentralAPIClient()
		logutil.FatalErr(err)
		logutil.FatalErr(
			createWebhook(rcAPIClient, opts.WebhookURL))
	}
	<-done
}
