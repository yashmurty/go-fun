package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

var urlStore SafeURLStore = SafeURLStore{urlMap: make(map[string]bool)}

// SafeURLStore declares maps for safe concurrent use.
type SafeURLStore struct {
	mutex  sync.Mutex
	urlMap map[string]bool
	wg     sync.WaitGroup
}

// crawledStatus checks if a URL has already been crawled or not.
func (str *SafeURLStore) crawledStatus(url string) bool {
	str.mutex.Lock()
	defer str.mutex.Unlock()

	// If url already exists, return true.
	if _, ok := str.urlMap[url]; ok {
		return true
	}
	// Else, set the crawl status to true but return false for existing status.
	str.urlMap[url] = true
	return false
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	defer urlStore.wg.Done()
	if depth <= 0 {
		return
	}
	// Don't fetch the same URL twice.
	if urlStore.crawledStatus(url) {
		return
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		urlStore.wg.Add(1)
		go Crawl(u, depth-1, fetcher)
	}
	return
}

func main() {
	urlStore.wg.Add(1)
	go Crawl("https://golang.org/", 6, fetcher)
	urlStore.wg.Wait()
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
