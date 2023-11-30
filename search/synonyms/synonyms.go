package synonyms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Head struct {
}

type Def struct {
	Text string `json:"text"`
	Pos  string `json:"pos"`
	Tr   []Tr   `json:"tr"`
}

type Tr struct {
	Text string `json:"text"`
	Pos  string `json:"pos"`
	Fr   int    `json:"fr"`
	Syn  []Syn  `json:"syn"`
}

type Syn struct {
	Text string `json:"text"`
	// Pos  string `json:"pos"`
	// Fr   int    `json:"fr"`
}

type Response struct {
	Head Head  `json:"head"`
	Def  []Def `json:"def"`
}

var apiPath = "https://dictionary.yandex.net/api/v1/dicservice.json/lookup"

func GetSynonyms(word string) ([]Syn, error) {

	apiKey := ""

	url := fmt.Sprintf(apiPath+"?key=%s&lang=ru-ru&text=%s", apiKey, word)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Go Synonyms Client")
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
	response := &Response{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}
	synonims := response.Def[0].Tr[0].Syn
	return synonims[0:min(3, len(synonims)-1)], nil
}
