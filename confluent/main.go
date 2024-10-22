package main

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"time"

	"sheeps.demo.go/confluent"
)

type ErroDetails struct {
	ID     string `json:"id"`
	Code   string `json:"code"`
	Status string `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

type ErrorResponse struct {
	Errors []ErroDetails `json:"errors"`
}

type Message struct {
	Data []struct {
		Metadata struct {
			Self string `json:"self"`
		} `json:"metadata"`
	} `json:"data"`
}

type SingleMessage struct {
	Data []struct {
		Spec struct {
			DisplayName           string `json:"display_name"`
			KafkaBoostrapEndpoint string `json:"kafka_boostrap_endpoint"`
			Region                string `json:"region"`
		} `json:"spec"`
	} `json:"data"`
}

var (
	ApiKey      = "your"
	AccessToken = "your"
	ApiUrl      = "http://localhost:40097/envrionments"
	pci         = false
	env         = "dev"
)

func main() {

	pciStr := "pci"
	if !pci {
		pciStr = "no" + pciStr
	}

	names := make(chan string)
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	msg, err := confluent.HttpGet[Message](ctx, ApiUrl, confluent.WithBasicAuthorization(ApiKey, AccessToken))

	if err := confluent.HasError(err); err != nil {
		var errRes ErrorResponse
		err.Unmarshal(&errRes)
		fmt.Printf("Error: %v\n%v\n", err.Status, err.Err)
		fmt.Printf("Error: %v\n", errRes)
		return
	}

	go func() {
		wg.Wait()
		close(names)
	}()

	for _, data := range msg.Data {
		wg.Add(1)
		go resolveHostname(ctx, data.Metadata.Self, names, &wg)
	}

	for name := range names {
		if findCurrentEnv(pciStr, env, name) {
			fmt.Printf("ACHEI: %v\n", name)
			cancel() // cancel all requests
			break
		}
	}

}

func resolveHostname(ctx context.Context, url string, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	msg, err := confluent.HttpGet[SingleMessage](ctx, url, confluent.WithBasicAuthorization(ApiKey, AccessToken))
	if err := confluent.HasError(err); err != nil {
		var errRes ErrorResponse
		fmt.Printf("Error: %v\n%v\n", err.Status, err.Err)
		err.Unmarshal(&errRes)
		fmt.Printf("Error: %v\n", errRes)
		return
	}

	for _, data := range msg.Data {
		ch <- data.Spec.DisplayName
	}
}

func findCurrentEnv(pciStr, envStr, name string) bool {

	matches := regexp.MustCompile(`(?i)((no)?pci)`).FindStringSubmatch(name)
	result := len(matches) >= 1 && matches[1] == pciStr

	matches = regexp.MustCompile(fmt.Sprintf(`(?i)(%v)`, envStr)).FindStringSubmatch(name)
	result = result && len(matches) > 0 && matches[0] == envStr

	return result
}
