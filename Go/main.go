package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ServerEmbedConfigData struct {
	DashboardId    string `json:"DashboardId"`
	ServerUrl      string `json:"ServerUrl"`
	UserEmail      string `json:"UserEmail"`
	EmbedSecret    string `json:"EmbedSecret"`
	EmbedType      string `json:"EmbedType"`
	Environment    string `json:"Environment"`
	ExpirationTime string `json:"ExpirationTime"`
	SiteIdentifier string `json:"SiteIdentifier"`
}

type ClientEmbedConfigData struct {
	DashboardId    string `json:"DashboardId"`
	ServerUrl      string `json:"ServerUrl"`
	EmbedType      string `json:"EmbedType"`
	Environment    string `json:"Environment"`
	SiteIdentifier string `json:"SiteIdentifier"`
}

var embedSecret string
var userEmail string

var serverEmbedConfigData ServerEmbedConfigData // Create an instance of ServerEmbedConfigData struct
var clientEmbedConfigData ClientEmbedConfigData // Create an instance of ClientEmbedConfigData struct

func main() {
	loadConfig() // Load configuration values
	http.HandleFunc("/authorizationServer", authorizationServer)
	http.HandleFunc("/getServerDetails", getServerDetails)
	log.Fatal(http.ListenAndServe(":8086", nil))
}

func loadConfig() {
	data, err := ioutil.ReadFile("embedConfig.json")
	if err != nil {
		log.Fatal("Error: embedConfig.json file not found.")
	}

	var serverEmbedConfigData ServerEmbedConfigData
	if err := json.Unmarshal(data, &serverEmbedConfigData); err != nil {
		log.Fatal("Error unmarshaling embedConfig.json:", err)
	}

	embedSecret = serverEmbedConfigData.EmbedSecret
	userEmail = serverEmbedConfigData.UserEmail
}

func getServerDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	data, err := ioutil.ReadFile("embedConfig.json")
	if err != nil {
		log.Fatal("Error: embedConfig.json file not found.")
	}

	err = json.Unmarshal(data, &clientEmbedConfigData)
	response, err := json.Marshal(clientEmbedConfigData)
	w.Write(response)
}

func authorizationServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	if len(body) > 0 {
		if queryString, err := unmarshal(string(body)); err != nil {
			log.Println("error converting", err)
		} else {
			serverAPIUrl := queryString.(map[string]interface{})["dashboardServerApiUrl"].(string)
			embedQueryString := queryString.(map[string]interface{})["embedQuerString"].(string)
			embedQueryString += "&embed_user_email=" + userEmail
			timeStamp := time.Now().Unix()
			embedQueryString += "&embed_server_timestamp=" + strconv.FormatInt(timeStamp, 10)
			signatureString, err := getSignatureUrl(embedQueryString)
			embedDetails := "/embed/authorize?" + embedQueryString + "&embed_signature=" + signatureString
			query := serverAPIUrl + embedDetails
			log.Println(query)
			result, err := http.Get(query)
			if err != nil {
				log.Println(err)
			}
			log.Println(result)
			response, err := ioutil.ReadAll(result.Body)
			if err != nil {
				log.Fatalln(err)
			}
			w.Write(response)
		}
	}
}

func getSignatureUrl(queryData string) (string, error) {
	encoding := ([]byte(embedSecret))
	messageBytes := ([]byte(queryData))
	hmacsha1 := hmac.New(sha256.New, encoding)
	hmacsha1.Write(messageBytes)
	sha := base64.StdEncoding.EncodeToString(hmacsha1.Sum(nil))
	return sha, nil
}

func unmarshal(data string) (interface{}, error) {
	var iface interface{}
	decoder := json.NewDecoder(strings.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&iface); err != nil {
		return nil, err
	}
	return iface, nil
}
