package main

import (
	"fmt"
	"sync"
)

// SafeURLMap is a list of strings
type SafeURLMap struct {
	urls map[string]bool
	mux  sync.RWMutex
}

// Add a url to the map
func (m *SafeURLMap) Add(url string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	if m.urls == nil {
		m.urls = make(map[string]bool)
	}
	m.urls[url] = true
}

// Exists checks if the given url is in the map
func (m *SafeURLMap) Exists(url string) (exists bool) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	_, exists = m.urls[url]
	return
}

func crawl(url string, depth int, fetcher Fetcher, urlmap *SafeURLMap, wg *sync.WaitGroup) {
	defer wg.Done()
	if depth <= 0 || urlmap.Exists(url) {
		return
	}
	urlmap.Add(url)
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		wg.Add(1)
		go crawl(u, depth-1, fetcher, urlmap, wg)
	}
	return
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	var wg sync.WaitGroup
	var urlmap SafeURLMap

	wg.Add(1)
	go crawl(url, depth, fetcher, &urlmap, &wg)

	wg.Wait()
}

func main() {
	Crawl("http://golang.org/", 4, fetcher)
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
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
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
