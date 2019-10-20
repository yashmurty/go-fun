# Concurrency

https://tour.golang.org/concurrency/10

### Exercise: Web Crawler
In this exercise you'll use Go's concurrency features to parallelize a web crawler.

Modify the Crawl function to fetch URLs in parallel without fetching the same URL twice.

Hint: you can keep a cache of the URLs that have been fetched on a map, but maps alone are not safe for concurrent use!

### Part 1/2 : 1_mutex

It achieves safe concurrent use of a map via `sync.Mutex`.

### Part 2/2 : 2_channel

It achieves safe concurrent use of a map via `channel`.
