# edgarwebcrawler
utillity for crawling edgar and reading RSS feeds

To get started with sampling the RSS feed:

first initialize a discarder. This is the built in default discarder that holds all seen links in memory and discards them by the RSS Guid. 
this one will hold 100 urls in memory
```go
SAMPLE_SIZE = 100
discarder := edgarwebcrawler.NewInMemorySampleDiscarderById(SAMPLE_SIZE)
```

next initialize the url provider.
you want the rss sample size to match to discarder's cache size. the second parameter is the sampling interval (in seconds), and lastly the discarder
```go
urlProvider := edgarwebcrawler.NewUrlFromRssProvider(
  fmt.Sprintf("https://www.sec.gov/cgi-bin/browse-edgar?action=getcurrent&type=4&start=-1&count=%d&output=rss", SAMPLE_SIZE),
  2,
  discarder,
  )
```

now all you need to do is create a string channel for the urls and start the provider
```go
urlChannel := make(chan string, SAMPLE_SIZE)
err := urlProvider.Start(urlChannel)
if err != nil {
  panic(err)
}
for {
  select {
    case url := <-urlChannel:
      // handle the url
  }
}
```
