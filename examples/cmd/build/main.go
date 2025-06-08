package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/grokify/go-glip/examples"
)

func main() {
	exps := examples.ExampleWebhooks()

	for _, exp := range exps {
		min, err := json.Marshal(exp.Message)
		if err != nil {
			log.Fatal(err)
		}
		max, err := json.MarshalIndent(exp.Message, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile(
			fmt.Sprintf("example_hook_%s_min.json", exp.Stub), min, 0600)
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile(
			fmt.Sprintf("example_hook_%s_sp2.json", exp.Stub), max, 0600)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("DONE")
}
