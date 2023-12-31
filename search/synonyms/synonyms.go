package synonyms

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	searchconfig "ssmt-ssu/search/searchConfig"
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

func GetSynonyms(word string, conf *searchconfig.Config) ([]Syn, error) {

	url := fmt.Sprintf(apiPath+"?key=%s&lang=ru-ru&text=%s",
		conf.YandexApiKey, word)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ssmt-ssu learning project")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	println(resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	response := &Response{}
	err = json.Unmarshal(body, response)
	if err != nil || len(response.Def) == 0 {
		return nil, err
	}
	synonims := response.Def[0].Tr[0].Syn
	return synonims[0:min(3, max(len(synonims)-1, 0))], nil
}
