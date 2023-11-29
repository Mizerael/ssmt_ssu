package main

import (
	"fmt"
	"ssmt-ssu/parser/parser"
	"sync"
)

var startLink = "/scp-series"

func main() {
	wg := new(sync.WaitGroup)
	for i := 0; i < 9; i++ {
		wg.Add(1)
		link := startLink
		if i != 0 {
			link += fmt.Sprintf("-%d", i)
		}
		go func() {
			defer wg.Done()
			parser.ParsePage(link)
		}()

	}
	wg.Wait()
	// parser.ParsePage(startLink)
}
