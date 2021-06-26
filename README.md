# IMDB Web Crawler

[![Build Status](https://img.shields.io/github/workflow/status/dominikrys/web-crawler/Continuous%20Integration?style=flat-square)](https://github.com/dominikrys/web-crawler/actions)

Web crawler for getting information on people born on a specified day from [IMDB](https://www.imdb.com/). The crawler returns information from the profiles of the people born on a specified date, where the people are sorted by popularity. Based off [Michael Okoko's blog post](https://blog.logrocket.com/web-scraping-with-go-and-colly/).

Note that there is rate limiting in place as the client may be blocked if too many requests are sent.

The aim of this project was to learn about Go and web scraping/crawling.

## Build and Run Instructions

Make sure [Go](https://golang.org/) is installed.

To compile, run:

```bash
go build ./crawler.go
```

Then, to run the program:

```bash
./crawler.go --day <day of birthday> --month <month of birthday>
```

Alternatively, for development, `go run` can be used:

```bash
go run . --day <day of birthday> --month <month of birthday>
```

To get more help on how to run the program, run:

```bash
./crawler --help
```

## TODO

- Store result in a database, or otherwise do something with the result. Also mention what I do in the README
  - Mention how to run MongoDB instance
  - Tweak CI job to make this work. May possibly need some tweaking?
  - Optional printing JSON? Optional MongoDB writing? Specify MongoDB database and names?
- Reflect repo headline with mongodb
- refactor stuff so no comments are needed
- add a test
