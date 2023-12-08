package facewho

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	searchconfig "ssmt-ssu/search/searchConfig"
	"strings"
	"sync"
)

const RuToEnApi = "https://api-inference.huggingface.co/models/Helsinki-NLP/opus-mt-ru-en"
const EnToRuApi = "https://api-inference.huggingface.co/models/Helsinki-NLP/opus-mt-en-ru"

func Translate(api string, text string, conf *searchconfig.Config) (string, error) {
	var translateText string
	var currentTextBlock string
	client := &http.Client{}
	sentences := strings.Split(text, ". ")
	sentencesLen := len(sentences)
	var wg sync.WaitGroup
	for i := 0; i < sentencesLen; i++ {
		currentTextBlock += sentences[i] + "."
		if i > 0 && i%6 == 0 || i == sentencesLen-1 {
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
