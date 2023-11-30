package spellshecker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Item struct {
	Code int      `json:"code"`
	Pos  int      `json:"pos"`
	Row  int      `json:"row"`
	Col  int      `json:"col"`
	Len  int      `json:"len"`
	Word string   `json:"word"`
	S    []string `json:"s"`
}

type JSON struct {
	Items []Item `json:""`
}

var apiPath = "https://speller.yandex.net/services/spellservice.json/checkTexts?text="

func Spellcheck(word string) ([]string, error) {

	url := fmt.Sprintf(apiPath + strings.Replace(word, " ", "+", 0))
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ssmt-ssu Spellcheck Client")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	println(resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	b := string(body)
	response := &JSON{}
	err = json.Unmarshal([]byte(b), &response.Items)
	if err != nil {
		return []string{}, err
	}
	checkedValue := response.Items[0].S
	println(checkedValue)
	return checkedValue[0:min(3, len(checkedValue)-1)], nil
}
