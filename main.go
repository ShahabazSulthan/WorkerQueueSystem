package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type Value struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

var (
	WorkerQueue = make(chan InputData, 10)
	Wg          sync.WaitGroup
)

type InputData struct {
	Event       string           `json:"ev"`
	EvenType    string           `json:"et"`
	AppID       string           `json:"id"`
	UserID      string           `json:"uid"`
	MessageID   string           `json:"mid"`
	PageTitle   string           `json:"t"`
	PageURL     string           `json:"p"`
	BrowserLang string           `json:"l"`
	ScreenSize  string           `json:"sc"`
	Attributes  map[string]Value `json:"-"`
	Traits      map[string]Value `json:"-"`
}

type TransformedData struct {
	Event       string           `json:"event"`
	EventType   string           `json:"event_type"`
	AppID       string           `json:"app_id"`
	UserID      string           `json:"user_id"`
	MessageID   string           `json:"message_id"`
	PageTitle   string           `json:"page_title"`
	PageURL     string           `json:"page_url"`
	BrowserLang string           `json:"browser_language"`
	ScreenSize  string           `json:"screen_size"`
	Attributes  map[string]Value `json:"attributes"`
	Traits      map[string]Value `json:"traits"`
}

func main() {
	for i := 0; i < 5; i++ {
		Wg.Add(1)
		go Worker()
	}

	http.HandleFunc("/event", HandleEvent)

	fmt.Println("Server Running 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))

	close(WorkerQueue)
	Wg.Wait()
}

func HandleEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only PostRequest are allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var input InputData
	if err := json.Unmarshal(body, &input); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	input.Attributes = ParseDynamicKeys(body, "atrk", "atrv", "atrt")
	input.Traits = ParseDynamicKeys(body, "uatrk", "uatrv", "uatrt")

	WorkerQueue <- input

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Request is being peocessed"))
}

func PostToWebhook(data TransformedData) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal data %v", err)
		fmt.Println("Failed to marshal :", err)
		return
	}

	webhookUrl := "https://webhook.site/68c410fb-7327-42cc-8467-a90a7f59c506"
	resp, err := http.Post(webhookUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to sent data webhook: %v ", err)
		fmt.Println("Failed to sent webhook : ", err)
		return
	}

	defer resp.Body.Close()

	log.Printf("Data sent to webhook ststud : %s ", resp.Status)
}

func TransFormData(input InputData) TransformedData {
	return TransformedData{
		Event:       input.Event,
		EventType:   input.EvenType,
		AppID:       input.AppID,
		UserID:      input.UserID,
		MessageID:   input.MessageID,
		PageTitle:   input.PageTitle,
		PageURL:     input.PageURL,
		BrowserLang: input.BrowserLang,
		ScreenSize:  input.ScreenSize,
		Attributes:  input.Attributes,
		Traits:      input.Traits,
	}
}

func Worker() {
	defer Wg.Done()

	for data := range WorkerQueue {
		transform := TransFormData(data)
		PostToWebhook(transform)
	}
}

func ParseDynamicKeys(body []byte, keyPrefix, valuePrefix, typePrefix string) map[string]Value {
	parsed := make(map[string]Value)
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return parsed
	}

	for k := range raw {
		if len(k) > len(keyPrefix) && k[:len(keyPrefix)] == keyPrefix {
			index := k[len(keyPrefix):]
			parsed[raw[keyPrefix+index].(string)] = Value{
				Value: raw[valuePrefix+index].(string),
				Type:  raw[typePrefix+index].(string),
			}
		}
	}
	return parsed
}
