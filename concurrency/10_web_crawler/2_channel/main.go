package main

import (
	"fmt"
	"log"
)

type CommandType int

const (
	GetCommand = iota
	SetCommand
	IncCommand
)

type Command struct {
	ty        CommandType
	url       string
	replyChan chan bool
}

// startURLManager starts a goroutine that serves as a manager for our
// counters datastore. Returns a channel that's used to send commands to the
// manager.
func startURLManager(initvals map[string]bool) chan<- Command {
	urlMap := make(map[string]bool)
	for k, v := range initvals {
		urlMap[k] = v
	}

	cmds := make(chan Command)

	go func() {
		for cmd := range cmds {
			switch cmd.ty {
			case GetCommand:
				if val, ok := urlMap[cmd.url]; ok {
					cmd.replyChan <- val
				} else {
					cmd.replyChan <- false
				}
			case SetCommand:
				urlMap[cmd.url] = true
				cmd.replyChan <- true
			default:
				log.Fatal("unknown command type", cmd.ty)
			}
		}
	}()
	return cmds
}

var urlManager = startURLManager(map[string]bool{})

// crawledStatus checks if a URL has already been crawled or not.
func crawledStatus(url string) bool {

	replyChan := make(chan bool)
	urlManager <- Command{ty: GetCommand, url: url, replyChan: replyChan}
	reply := <-replyChan

	// If url already exists, return true.
	if reply == true {
		return true
	}
	// Else, set the crawl status to true but return false for existing status.
	urlManager <- Command{ty: SetCommand, url: url, replyChan: replyChan}
	_ = <-replyChan

	return false
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, done chan bool) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	if depth <= 0 {
		done <- true
		return
	}
	// Don't fetch the same URL twice.
	if crawledStatus(url) {
		done <- true
		return
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		done <- true
		return
	}
	fmt.Printf("found: %s %q\n", url, body)

	// Wait for the children crawls.
	childrenDone := make(chan bool, 1)
	for _, u := range urls {
		go Crawl(u, depth-1, fetcher, childrenDone)
	}
	for i := 0; i < len(urls); i++ {
		<-childrenDone
	}

	// Mark the the crawl function as finished.
	done <- true
	return
}

func main() {
	done := make(chan bool)
	go Crawl("https://golang.org/", 4, fetcher, done)
	<-done
	return
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
