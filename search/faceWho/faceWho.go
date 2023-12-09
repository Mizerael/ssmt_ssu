package facewho

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	searchconfig "ssmt-ssu/search/searchConfig"
	"strconv"
	"strings"
	"sync"

	"github.com/blevesearch/bleve/v2"
)

const RuToEnApi = "https://api-inference.huggingface.co/models/Helsinki-NLP/opus-mt-ru-en"
const EnToRuApi = "https://api-inference.huggingface.co/models/Helsinki-NLP/opus-mt-en-ru"
const SummarizeApi = "https://api-inference.huggingface.co/models/facebook/bart-large-cnn"
const SimilarirtyApi = "https://api-inference.huggingface.co/models/sentence-transformers/all-MiniLM-L6-v2"

type SimmilarityRequest struct {
	SourceSentence string   `json:"source_sentence"`
	Sentences      []string `json:"sentences"`
}

type Request struct {
	Inputs SimmilarityRequest `json:"inputs"`
}

type Kv struct {
	Key   string
	Value float64
}

func Translate(api string, text string, conf *searchconfig.Config) (string, error) {
	var translateText string
	var currentTextBlock string
	client := &http.Client{}
	sentences := strings.Split(text, ". ")
	sentencesLen := len(sentences)
	var wg sync.WaitGroup
	for i := 0; i < sentencesLen; i++ {
		currentTextBlock += sentences[i] + "."
		if i > 0 && i%3 == 0 || i == sentencesLen-1 {
			wg.Add(1)
			values := map[string]string{"inputs": currentTextBlock}
			data, _ := json.Marshal(values)
			req, err := http.NewRequest("POST", api, bytes.NewBuffer(data))
			if err != nil {
				log.Printf("err: %v\n", err)
				return "", err
			}
			req.Header.Set("User-Agent", "ssmt-ssu learning project")
			req.Header.Set("Authorization", "Bearer "+conf.HuggingfaceApiKey)
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				return "", err
			}
			defer resp.Body.Close()

			println(resp.Status)
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", err
			}
			var result []map[string]interface{}
			err = json.Unmarshal(body, &result)
			if err != nil {
				fmt.Printf("err: %v\n", err)
				return "", err
			}
			currentTextBlock = ""
			translateText += fmt.Sprintf("%v", result[0]["translation_text"])
			wg.Done()
		}
		wg.Wait()
	}
	return translateText, nil
}

func Summarize(api string, text string, conf *searchconfig.Config) (string, error) {
	client := &http.Client{}

	values := map[string]string{"inputs": text}
	data, _ := json.Marshal(values)
	req, err := http.NewRequest("POST", api, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("err: %v\n", err)
		return "", err
	}
	req.Header.Set("User-Agent", "ssmt-ssu learning project")
	req.Header.Set("Authorization", "Bearer "+conf.HuggingfaceApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	println(resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result []map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return "", err
	}

	return fmt.Sprintf("%v", result[0]["summary_text"]), nil
}

func Simmilarity(api string, source string,
	searchResult *bleve.SearchResult, conf *searchconfig.Config) []Kv {
	client := &http.Client{}
	var sentences []string
	result := make(map[string]float64)
	for _, v := range searchResult.Hits {
		sentences = append(sentences, fmt.Sprintf("%v", v.Fields["contains"]))
	}

	values := Request{
		Inputs: SimmilarityRequest{
			SourceSentence: source,
			Sentences:      sentences,
		},
	}
	data, _ := json.Marshal(values)
	req, err := http.NewRequest("POST", api, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("err: %v\n", err)
	}
	req.Header.Set("User-Agent", "ssmt-ssu learning project")
	req.Header.Set("Authorization", "Bearer "+conf.HuggingfaceApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("err: %v\n", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("err: %v\n", err)
	}

	str := strings.Split(strings.Replace(strings.Replace(string(body), "]", "", -1),
		"[", "", -1), ",")
	for i, v := range searchResult.Hits {
		if n, err := strconv.ParseFloat(str[i], 64); err == nil {
			result[fmt.Sprintf("%v", v.Fields["name"])] = n
		}
	}
	var ss []Kv
	for k, v := range result {
		ss = append(ss, Kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	return ss[:conf.Count]
}
