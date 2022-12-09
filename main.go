package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Stats struct {
	Players []struct {
		Name     string `json:"nomaffiche"`
		Position string `json:"position"`
		Club     string `json:"club"`
		Data     []struct {
			Message string      `json:"message"`
			Value   interface{} `json:"value"`
		} `json:"criteres"`
	} `json:"joueurs"`
}

func main() {
	err := getStats("https://gamezone.autumnnationsseries.com/v1/private/stats?lg=en", "POST")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	statsContent, err := os.Open("stats.json")

	if err != nil {
		fmt.Printf("Got error %v", err)
	}

	defer statsContent.Close()

	byteResult, err := ioutil.ReadAll(statsContent)

	if err != nil {
		fmt.Printf("Got error: %v", err)
	}

	var stats Stats

	err = json.Unmarshal(byteResult, &stats)
	if err != nil {
		fmt.Printf("Got error: %v", err)
	}

	// Do something with stats
	fmt.Println(stats)
}

func getStats(url, method string) error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	var dataStr = "{\"credentials\":{\"critereRecherche\":{\"nom\":\"\",\"club\":\"\",\"position\":\"\"},\"critereTri\":\"moyenne_points\",\"loadSelect\":0,\"pageIndex\":0,\"pageSize\":15}}"
	var jsonStr = []byte(dataStr)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}

	req.Header.Set("accept", "application.json")
	req.Header.Add("accept-language", "en-US,en;q=0.9")
	req.Header.Add("authorization", "Token eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpYXQiOjE2Njg5NTY2MzgsImV4cCI6MTY3MTU0ODYzOCwianRpIjoidE9PRm0xa2ZVdlBEVm5EYzN4eTZ6UT09IiwiaXNzIjoiaHR0cHM6XC9cL2dhbWV6b25lLmF1dHVtbm5hdGlvbnNzZXJpZXMuY29tXC9mYW50YXN5Iiwic3ViIjp7ImlkIjoiMjY4ODUxIiwibWFpbCI6Imdhdmluc3BhdHRvbkBnbWFpbC5jb20iLCJtYW5hZ2VyIjoiR2F2aW4iLCJpZGwiOiIxIiwiaWRnIjoiMSIsImZ1c2VhdSI6IkV1cm9wZVwvTG9uZG9uIiwibWVyY2F0byI6MCwiaWRqZyI6IjE2MDc5IiwiaXNhZG1pbmNsaWVudCI6ZmFsc2UsImlzYWRtaW4iOmZhbHNlLCJpc3N1cGVyYWRtaW4iOmZhbHNlLCJ2aXAiOmZhbHNlLCJpZGVudGl0eSI6Ijg2MSIsImlnbm9yZWNvZGUiOmZhbHNlLCJjb2RlIjoiODYxLjEiLCJjb2RlRjUiOiI4NjEuMSJ9fQ.fZo4MD9caOE__c_-mmWFiJ9OoSRJWtVXQ1G9kgeeryo")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-access-key", "861@3.23@")

	fmt.Println("Getting stats...")
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}

	defer response.Body.Close()

	out, err := os.OpenFile("stats.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}

	defer out.Close()
	io.Copy(out, response.Body)

	return err
}
